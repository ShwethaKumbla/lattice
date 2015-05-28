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

	BeforeEach(func() {
		cliApp = NewCliApp()
	})

	Describe("NewCliApp", func() {
		It("Runs registered command without error", func() {
			commandRan := false
			cliApp.Commands = []cli.Command{
				cli.Command{
					Name:   "print-a-unicorn",
					Action: func(ctx *cli.Context) { commandRan = true },
				},
			}

			cliAppArgs := []string{"ltc", "print-a-unicorn"}
			err := cliApp.Run(cliAppArgs)

			Expect(err).NotTo(HaveOccurred())
		})
	})

})
