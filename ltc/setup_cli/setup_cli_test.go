package setup_cli_test

import (
	. "github.com/cloudfoundry-incubator/lattice/ltc/setup_cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"


	"github.com/codegangsta/cli"
)

var _ = Describe("SetupCli", func() {
	var (
		cliApp *cli.App
	)


	JustBeforeEach(func() {
		cliApp = NewCliApp()
	})
	
	Describe("MatchArgAndFlags", func() {
		It("Checks for badflag", func() {
			cliAppArgs := []string{"ltc", "create", "--badflag"}
			flags := GetCommandFlags(cliApp, cliAppArgs[1])
			badFlags := MatchArgAndFlags(flags, cliAppArgs[2:])

			Expect(badFlags).To(Equal("Unknown flag \"--badflag\""))

		})

		It("returns if multiple bad flags are passed", func() {
			cliAppArgs := []string{"ltc", "create", "--badflag1", "--badflag2"}
			flags := GetCommandFlags(cliApp, cliAppArgs[1])
			badFlags := MatchArgAndFlags(flags, cliAppArgs[2:])
			InjectHelpTemplate(badFlags)
			Expect(badFlags).To(Equal("Unknown flags: \"--badflag1\", \"--badflag2\""))

		})
	})

	Describe("GetCommandFlags", func() {
		It("returns list of type Flag", func() {
			flaglist := GetCommandFlags(cliApp, "create")
			cmd := cliApp.Command("create")
			for _, flag := range cmd.Flags {
				switch t := flag.(type) {
				default:
				case cli.StringSliceFlag:
					Expect(flaglist).Should(ContainElement(t.Name))
				case cli.IntFlag:
					Expect(flaglist).Should(ContainElement(t.Name))
				case cli.StringFlag:
					Expect(flaglist).Should(ContainElement(t.Name))
				case cli.BoolFlag:
					Expect(flaglist).Should(ContainElement(t.Name))
				case cli.DurationFlag:
					Expect(flaglist).Should(ContainElement(t.Name))
				}
			}
		})
	})
	
	Describe("GetByCmdName", func() {
		It("returns command not found error", func() {
			_, err := GetByCmdName(cliApp, "zz")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Command not found"))
		})
	})
	
	Describe("RequestHelp", func() {
		It("checks for the flag -h", func() {
			cliAppArgs := []string{"ltc", "-h"}
			boolVal := RequestHelp(cliAppArgs[1:])
			Expect(boolVal).To(BeTrue())
		})

		It("checks for the flag --help", func() {
			cliAppArgs := []string{"ltc", "--help"}
			boolVal := RequestHelp(cliAppArgs[1:])
			Expect(boolVal).To(BeTrue())
		})

		It("checks for the unknown flag", func() {
			cliAppArgs := []string{"ltc", "--unknownFlag"}
			boolVal := RequestHelp(cliAppArgs[1:])
			Expect(boolVal).To(BeFalse())
		})
		
	})
})
