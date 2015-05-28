package cli_app_factory_test

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/lattice/ltc/cli_app_factory"
	"github.com/cloudfoundry-incubator/lattice/ltc/config"
	"github.com/cloudfoundry-incubator/lattice/ltc/config/persister"
	"github.com/cloudfoundry-incubator/lattice/ltc/terminal"
	"github.com/codegangsta/cli"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("AppHelp", func() {
	var (
		cliApp       *cli.App
		outputBuffer *gbytes.Buffer
	)

	dummyTemplate := `
		{{range .Commands}}{{range .CommandSubGroups}}{{range .}}
		{{.Name}}
		{{end}}{{end}}{{end}}`

	subCommandTemplate := `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} command{{if .Flags}} [command options]{{end}} [arguments...]
`
	BeforeEach(func() {
		outputBuffer = gbytes.NewBuffer()

		cliApp = cli_app_factory.MakeCliApp(
			"",
			"~/",
			nil,
			config.New(persister.NewMemPersister()),
			nil,
			nil,
			terminal.NewUI(nil, outputBuffer, nil),
		)
	})

	It("shows help for all commands", func() {

		cliCommands := cliApp.Commands
		Expect(cliCommands).NotTo(BeEmpty())

		cli_app_factory.ShowHelp(outputBuffer, dummyTemplate, cliApp)

		outputBytes, err := ioutil.ReadAll(outputBuffer)
		Expect(err).NotTo(HaveOccurred())
		for _, command := range cliCommands {
			commandName := strings.TrimSpace(strings.Join(command.Names(), ", "))
			Expect(string(outputBytes)).To(ContainSubstring(commandName))
		}
	})

	Describe("commandPrintHelp", func() {
		It("Shows help for particular command", func() {
			commandRan := false
			subCommand := cli.Command{
				Name:        "print-a-command",
				ShortName:   "p",
				Description: "Print command",
				Usage:       "print-a-command [arguments]",
				Action:      func(ctx *cli.Context) { commandRan = true },
			}
			cli_app_factory.ShowHelp(outputBuffer, subCommandTemplate, subCommand)
			outputBytes, err := ioutil.ReadAll(outputBuffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(outputBytes)).To(ContainSubstring(subCommand.Name))
		})
	})
})
