package main

import (
	"fmt"
	"github.com/cloudfoundry/cli/plugin/fakes"
	"github.com/cloudfoundry/cli/plugin/models"
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func fakeError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var _ = Describe("WildcardPlugin", func() {
	var (
		ui                *testterm.FakeUI
		wildcardPlugin    *Wildcard
		fakeCliConnection *fakes.FakeCliConnection
		appsList          []plugin_models.GetAppsModel
	)
	Context("When running wildcard-apps", func() {
		BeforeEach(func() {
			appsList = make([]plugin_models.GetAppsModel, 0)
			appsList = append(appsList,
				plugin_models.GetAppsModel{"spring-music", "", "", 0, 0, 0, 0, nil},
				plugin_models.GetAppsModel{"app321", "", "", 0, 0, 0, 0, nil},
			)
			fakeCliConnection = &fakes.FakeCliConnection{}
			wildcardPlugin = &Wildcard{}
		})

		Describe("When there are matching apps", func() {
			It("prints a table containing only those apps", func() {
				fakeCliConnection.GetAppsReturns(appsList, nil)
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-apps", "app*"})
				})

				Expect(output).To(ContainSubstrings(
					[]string{"app321"},
				))
				Expect(output).ToNot(ContainSubstrings(
					[]string{"spring-music"},
				))
			})
		})

		Describe("When the user provides incorrect input", func() {
			It("prints correct usage", func() {
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-apps", "app*", "123*"})
				})

				Expect(output).To(ContainSubstrings(
					[]string{"Usage"},
					[]string{"cf wildcard-apps"},
					[]string{"APP_NAME_WITH_WILDCARD"},
				))
			})
		})
	})

	Context("wildcard-delete -f", func() {
		Describe("When there are matching apps", func() {
			BeforeEach(func() {
				appsList = make([]plugin_models.GetAppsModel, 0)
				appsList = append(appsList,
					plugin_models.GetAppsModel{"spring-music", "", "", 0, 0, 0, 0, nil},
					plugin_models.GetAppsModel{"app321", "", "", 0, 0, 0, 0, nil},
					plugin_models.GetAppsModel{"apple_pie", "", "", 0, 0, 0, 0, nil},
				)
				fakeCliConnection = &fakes.FakeCliConnection{}
				fakeCliConnection.GetAppsReturns(appsList, nil)
				wildcardPlugin = &Wildcard{}
				ui = &testterm.FakeUI{}
			})

			It("does not prompt when the user provides the -f flag", func() {
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-delete", "app*", "-f"})
				})
				Expect(output).To(ContainSubstrings(
					[]string{"Deleting app app321"},
					[]string{"Deleting app apple_pie"},
				))
				Expect(output).ToNot(ContainSubstrings(
					[]string{"Deleting app spring-music"},
				))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[0]).To(Equal("delete"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[1]).To(Equal("app321"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[2]).To(Equal("-f"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[0]).To(Equal("delete"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[1]).To(Equal("apple_pie"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[2]).To(Equal("-f"))
			})
			It("does not prompt and deletes all mapped routes when the user provides the -f and -r flag", func() {
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-delete", "app*", "-f", "-r"})
				})
				Expect(output).To(ContainSubstrings(
					[]string{"Deleting app app321 and its mapped routes"},
					[]string{"Deleting app apple_pie and its mapped routes"},
				))
				Expect(output).ToNot(ContainSubstrings(
					[]string{"Deleting app spring-music and its mapped routes"},
				))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[0]).To(Equal("delete"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[1]).To(Equal("app321"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[2]).To(Equal("-f"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[3]).To(Equal("-r"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[0]).To(Equal("delete"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[1]).To(Equal("apple_pie"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[2]).To(Equal("-f"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[3]).To(Equal("-r"))
			})
			It("does not matter what the order of the flags, -f and -r, are", func() {
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-delete", "app*", "-r", "-f"})
				})
				Expect(output).To(ContainSubstrings(
					[]string{"Deleting app app321 and its mapped routes"},
					[]string{"Deleting app apple_pie and its mapped routes"},
				))
				Expect(output).ToNot(ContainSubstrings(
					[]string{"Deleting app spring-music and its mapped routes"},
				))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[0]).To(Equal("delete"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[1]).To(Equal("app321"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[2]).To(Equal("-f"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)[3]).To(Equal("-r"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[0]).To(Equal("delete"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[1]).To(Equal("apple_pie"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[2]).To(Equal("-f"))
				Expect(fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)[3]).To(Equal("-r"))
			})

		})
	})
	Context("When there are no matching apps", func() {
		BeforeEach(func() {
			appsList = make([]plugin_models.GetAppsModel, 0)
			appsList = append(appsList,
				plugin_models.GetAppsModel{"spring-music", "", "", 0, 0, 0, 0, nil},
				plugin_models.GetAppsModel{"qwerty", "", "", 0, 0, 0, 0, nil},
				plugin_models.GetAppsModel{"apple_pie", "", "", 0, 0, 0, 0, nil},
			)
			fakeCliConnection = &fakes.FakeCliConnection{}
			wildcardPlugin = &Wildcard{}
			ui = &testterm.FakeUI{}
		})
		Describe("When there are no matching apps", func() {
			It("prints an empty table and informs the user", func() {
				fakeCliConnection.GetAppsReturns(appsList, nil)
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-apps", "foo*"})
				})

				Expect(output).To(ContainSubstrings(
					[]string{"name"},
					[]string{"requested state"},
					[]string{"instances"},
					[]string{"No apps found matching foo*"},
				))
				Expect(output).ToNot(ContainSubstrings(
					[]string{"spring-music"},
					[]string{"app321"},
				))
			})
		})
		Describe("When there are no matching apps", func() {
			It("prints no apps found", func() {
				output := io_helpers.CaptureOutput(func() {
					wildcardPlugin.Run(fakeCliConnection, []string{"wildcard-delete", "foo*", "-f"})
				})
				Expect(output).To(ContainSubstrings(
					[]string{"No apps found matching foo*"},
				))
				Expect(output).ToNot(ContainSubstrings(
					[]string{"Deleting app"},
				))
			})
		})
	})
})
