package setup_cli

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/cloudfoundry-incubator/lattice/ltc/cli_app_factory"
	"github.com/cloudfoundry-incubator/lattice/ltc/config"
	"github.com/cloudfoundry-incubator/lattice/ltc/config/config_helpers"
	"github.com/cloudfoundry-incubator/lattice/ltc/config/persister"
	"github.com/cloudfoundry-incubator/lattice/ltc/config/target_verifier"
	"github.com/cloudfoundry-incubator/lattice/ltc/config/target_verifier/receptor_client_factory"
	"github.com/cloudfoundry-incubator/lattice/ltc/exit_handler"
	"github.com/codegangsta/cli"
	"github.com/pivotal-golang/lager"
)

const (
	latticeCliHomeVar = "LATTICE_CLI_HOME"
)

var (
	latticeVersion string // provided by linker argument at compile-time
)

func NewCliApp() *cli.App {
	config := config.New(persister.NewFilePersister(config_helpers.ConfigFileLocation(ltcConfigRoot())))

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	exitHandler := exit_handler.New(signalChan, os.Exit)
	go exitHandler.Run()

	targetVerifier := target_verifier.New(receptor_client_factory.MakeReceptorClient)
	app := cli_app_factory.MakeCliApp(latticeVersion, ltcConfigRoot(), exitHandler, config, logger(), targetVerifier, os.Stdout)
	return app
}

func logger() lager.Logger {
	logger := lager.NewLogger("ltc")
	var logLevel lager.LogLevel

	if os.Getenv("LTC_LOG_LEVEL") == "DEBUG" {
		logLevel = lager.DEBUG
	} else {
		logLevel = lager.INFO
	}

	logger.RegisterSink(lager.NewWriterSink(os.Stderr, logLevel))
	return logger
}

func ltcConfigRoot() string {
	if os.Getenv(latticeCliHomeVar) != "" {
		return os.Getenv(latticeCliHomeVar)
	}

	return os.Getenv("HOME")
}

func InjectHelpTemplate(badFlags string) {
	cli.CommandHelpTemplate = fmt.Sprintf(`%sNAME:
   {{join .Names ", "}} - {{.Usage}}
{{with .ShortName}}
ALIAS:
   {{.Aliases}}
{{end}}
USAGE:
   {{.Description}}{{with .Flags}}
OPTIONS:
{{range .}}   {{.}}
{{end}}{{else}}
{{end}}`, badFlags)
}

func MatchArgAndFlags(flags []string, args []string) string {
	var badFlag string
	var lastPassed bool
	multipleFlagErr := false
Loop:
	for _, arg := range args {
		prefix := ""
		//only take flag name, ignore value after '='
		arg = strings.Split(arg, "=")[0]
		if arg == "--h" || arg == "-h" || arg == "--help" || arg == "-help" {
			continue Loop
		}
		if strings.HasPrefix(arg, "--") {
			prefix = "--"
		} else if strings.HasPrefix(arg, "-") {
			prefix = "-"
		}
		arg = strings.TrimLeft(arg, prefix)
		//skip verification for negative integers, e.g. -i -10
		if lastPassed {
			lastPassed = false
			if _, err := strconv.ParseInt(arg, 10, 32); err == nil {
				continue Loop
			}
		}
		if prefix != "" {
			for _, flag := range flags {
				for _, f := range strings.Split(flag, ", ") {
					flag = strings.TrimSpace(f)
					if flag == arg {
						lastPassed = true
						continue Loop
					}
				}
			}
			if badFlag == "" {
				badFlag = fmt.Sprintf("\"%s%s\"", prefix, arg)
			} else {
				multipleFlagErr = true
				badFlag = badFlag + fmt.Sprintf(", \"%s%s\"", prefix, arg)
			}
		}
	}
	if multipleFlagErr && badFlag != "" {
		badFlag = fmt.Sprintf("%s %s", "Unknown flags:", badFlag)
	} else if badFlag != "" {
		badFlag = fmt.Sprintf("%s %s", "Unknown flag", badFlag)
	}
	return badFlag
}

func RequestHelp(args []string) bool {
	for _, v := range args {
		if v == "-h" || v == "--help" {
			return true
		}
	}
	return false
}

func CallCoreCommand(args []string, cliApp *cli.App) {
	err := cliApp.Run(args)
	if err != nil {
		os.Exit(1)
	}
}

func GetCommandFlags(app *cli.App, command string) []string {
	cmd, err := GetByCmdName(app, command)
	if err != nil {
		return []string{}
	}
	var flags []string
	for _, flag := range cmd.Flags {
		switch t := flag.(type) {
		default:
		case cli.StringSliceFlag:
			flags = append(flags, t.Name)
		case cli.IntFlag:
			flags = append(flags, t.Name)
		case cli.StringFlag:
			flags = append(flags, t.Name)
		case cli.BoolFlag:
			flags = append(flags, t.Name)
		case cli.DurationFlag:
			flags = append(flags, t.Name)
		}
	}
	return flags
}

func GetByCmdName(app *cli.App, cmdName string) (cmd *cli.Command, err error) {
	cmd = app.Command(cmdName)
	if cmd == nil {
		for _, c := range app.Commands {
			if c.ShortName == cmdName {
				return &c, nil
			}
		}
		err = errors.New("Command not found")
	}
	return
}
