//helpful link for match: https://golang.org/src/path/filepath/match_test.go
//GetApps()'s apps-summary is WAY too limited!!! They REALLY need to improve it >:(
//

package main

import (
	"errors"
	"fmt" //standard
	"os"
	//"reflect" //used to see type of object
	"strconv"
	"github.com/guidowb/cf-go-client/panic" //panics 
	"strings"
	"github.com/cloudfoundry/cli/plugin/models"
	//"github.com/cloudfoundry/cli/cf/api"
	//"github.com/cloudfoundry/cli/cf/formatters"
	"github.com/cloudfoundry/cli/plugin" //standard//https://github.com/cloudfoundry/cli/blob/8c310da376377c53f001d916708c056ce1558959/plugin/plugin.go

	"path/filepath" //for matches//https://golang.org/pkg/path/filepath/
	"github.com/cloudfoundry/cli/cf/terminal" //for table || https://github.com/cloudfoundry/cli/blob/4a108fd21d6633b250f6d9f46e870967cae96ac0/cf/terminal/table.go


	//for implementing T
	"github.com/cloudfoundry/cli/cf/trace"
	. "github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/i18n/detection"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	//for adding onto table
	"github.com/cloudfoundry/cli/cf/formatters"
	//"github.com/cloudfoundry/cli/cf/ui_helpers"
	//Prompts
	//"github.com/codegangsta/cli"


)

//Wildcard is this plugin
type Wildcard struct {
	ui 				terminal.UI
	matchedApps 	[]plugin_models.GetAppsModel
	pattern 		string
	err 			error
}

//GetMetadata returns metatada
func (cmd *Wildcard) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "wildcard",
		Version: plugin.VersionType{ //leavealone
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{  //****** array of command structures
			{
				Name:     "wildcard-apps",
				Alias:	  "wc-a",
				HelpText: "List all apps in the target space matching the wildcard",
				UsageDetails: plugin.Usage{
					Usage: "cf wildcard-apps APP_NAME_WITH_WILDCARD",
				},
			}, 
			{
				Name:     "wildcard-delete",
				Alias:	  "wc-d",
				HelpText: "Delete apps in the target space matching the wildcard",
				UsageDetails: plugin.Usage{
					Usage: "cf wildcard-delete APP_NAME_WITH_WILDCARD",
				},
			},
		},
	}
}

func main() { 
	plugin.Start(newWildcard())
}

func newWildcard() (*Wildcard) {
	return &Wildcard{ui: terminal.NewUI(os.Stdin, terminal.NewTeePrinter())}
}

func (cmd *Wildcard) usage(args []string) error {
	badArgs := 2 != len(args)
	if badArgs {
		if args[0] == "wildcard-apps" {
			return errors.New("Usage: cf wildcard-apps\n\tcf wildcard-apps APP_NAME_WITH_WILDCARD")
		} else if args[0] == "wildcard-delete" {
			return errors.New("Usage: cf wildcard-delete\n\tcf wildcard-delete APP_NAME_WITH_WILDCARD")
		}
	}
	return nil
}

//Run runs the plugin
//called everytime user executes the command
func (cmd *Wildcard) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "wildcard-apps" { //checking is very imp.
		cmd.WildcardCommandApps(cliConnection, args)
	} else if args[0] == "wildcard-delete" {
		cmd.WildcardCommandDelete(cliConnection, args)
	}
}

//WildcardCommand creates a new instance of this plugin
//this is the actual implementation
//one method per command
func InitializeCliDependencies() {
	errorHandler := func(err error) {
		if err != nil {
			fmt.Sprintf("Config error: %s", err)
		}
	}
	cc_config := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), errorHandler)
	T = Init(cc_config, &detection.JibberJabberDetector{})
	if os.Getenv("CF_TRACE") != "" {
		trace.Logger = trace.NewLogger(os.Getenv("CF_TRACE"))
	} else {
		trace.Logger = trace.NewLogger(cc_config.Trace())
	}
}
//Q: How come when I add the following, the command exits?
// func (cmd *Wildcard) introduction(cliConnection plugin.CliConnection, args []string) {
// 	currOrg, _ := cliConnection.GetCurrentOrg()
// 	currSpace, _ := cliConnection.GetCurrentSpace()
// 	currUsername, _ := cliConnection.Username()
// 	cmd.ui.Say(T("Getting apps in org {{.OrgName}} / space {{.SpaceName}} as {{.Username}}...",
// 		map[string]interface{}{
// 			"OrgName":   terminal.EntityNameColor(currOrg.Name),
// 			"SpaceName": terminal.EntityNameColor(currSpace.Name),
// 			"Username":  terminal.EntityNameColor(currUsername)}))
// 	cmd.ui.Ok()
// 	cmd.ui.Say("")
// }

func (cmd *Wildcard) getMatchedApps(cliConnection plugin.CliConnection, args []string) ([]plugin_models.GetAppsModel) {
	if err := cmd.usage(args); err != nil {
		fmt.Println(err) //printing
		os.Exit(1) //failure
	}
	cmd.pattern = args[1]
	//cmd.introduction(cliConnection, args)
	output, _ := cliConnection.GetApps()
	for i := 0; i < (len(output)); i++ {
		ok, _ := filepath.Match(cmd.pattern, output[i].Name)
		if ok {
			cmd.matchedApps = append(cmd.matchedApps, output[i])
		}
	}
	if len(cmd.matchedApps) <= 0 {
		fmt.Printf("No apps matching %q found", cmd.pattern)
		fmt.Println("")
		os.Exit(1)
	}
	return cmd.matchedApps
}
func (cmd *Wildcard) WildcardCommandApps(cliConnection plugin.CliConnection, args []string) {
	cmd.getMatchedApps(cliConnection, args)
	InitializeCliDependencies()
	defer panic.HandlePanics()
	table := terminal.NewTable(cmd.ui, []string{T("name"), T("requested state"), T("instances"), T("memory"), T("disk"), T("urls")})
	for _, app := range cmd.matchedApps {
		var urls []string
		for _, route := range app.Routes {
			if route.Host == "" { 
				urls = append(urls, route.Domain.Name)
			}
			urls = append(urls, fmt.Sprintf("%s.%s", route.Host, route.Domain.Name))
		}
		table.Add(
			app.Name,
			app.State, 
			strconv.Itoa(app.RunningInstances),
			formatters.ByteSize(app.Memory*formatters.MEGABYTE),
			formatters.ByteSize(app.DiskQuota*formatters.MEGABYTE),
			strings.Join(urls, ", "),
		)
	}
	table.Print()
}

func (cmd *Wildcard) WildcardCommandDelete(cliConnection plugin.CliConnection, args []string) {
	cmd.WildcardCommandApps(cliConnection, args)
	response := cmd.ui.Ask("Would you like to delete the apps (i)nteractively, (a)ll, or (c)ancel this command?")
	if !strings.EqualFold(response,"a") && !strings.EqualFold(response,"all") && !strings.EqualFold(response,"i") && !strings.EqualFold(response,"interactively") {
		fmt.Printf("Delete cancelled")
		fmt.Println("")
		os.Exit(1)
	} else {
		for _, app := range cmd.matchedApps {
			if strings.EqualFold(response,"i") || strings.EqualFold(response,"interactively"){
				cliConnection.CliCommandWithoutTerminalOutput("delete", app.Name)
			} else if strings.EqualFold(response,"a") || strings.EqualFold(response,"all") {
				confirmation := cmd.ui.Confirm("Really delete all apps matching %q ?", cmd.pattern)
				if !confirmation {
					cmd.ui.Warn(T("Delete all cancelled"))
					os.Exit(1)
				} else {
					fmt.Println("Deleting all apps matching %q ", cmd.pattern)
					cliConnection.CliCommandWithoutTerminalOutput("delete", app.Name, "-f")
					
				}
				cmd.ui.Ok()
			} else {
				return
			}
		}
	}
}






