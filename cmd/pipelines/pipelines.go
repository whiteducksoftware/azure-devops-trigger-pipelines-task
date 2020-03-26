/*
Copyright © 2020 Stefan Kürzeder <stefan.kuerzeder@whiteduck.de

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package pipelines

import (
	"errors"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/cmd/pipelines/trigger"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/pkg/azure"
)

// Cmd represents the pipelines command
var Cmd = &cobra.Command{
	Use:   "pipelines",
	Short: "Utilities for Azure DevOps pipelines",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		url, err := flags.GetString(azure.OrganizationUrlFlagName.Name)
		if err != nil {
			return err
		}

		pat, err := flags.GetString(azure.PersonalAccessTokenFlagName.Name)
		if err != nil {
			return err
		}

		// Check for a valid devops organization url
		urlMatchA, err := regexp.MatchString(`(?i)^https://dev.azure.com/[a-zA-Z0-9]{1,256}$`, url)
		if err != nil {
			return err
		}

		urlMatchB, err := regexp.MatchString(`(?i)^https://[a-zA-Z0-9]{1,256}.visualstudio.com$`, url)
		if err != nil {
			return err
		}

		urlMatchC, err := regexp.MatchString(`(?i)^https://(www\.)?([^.]+)\.([^.]+.)+/tfs/DefaultCollection$`, url)
		if err != nil {
			return err
		}

		if urlMatchA == false && urlMatchB == false && urlMatchC {
			return errors.New("Invalid Orgranization Url has been passed")
		}

		patMatch, err := regexp.MatchString(`^[a-z0-9]{52}$`, pat)
		if err != nil {
			return err
		}

		if patMatch == false {
			return errors.New("Invalid PAT has been passed")
		}

		return nil
	},
}

func init() {
	// Add sub-Commands
	Cmd.AddCommand(trigger.Cmd)

	// Define flags
	azure.AddFlags(Cmd, []azure.FlagDefinition{azure.OrganizationUrlFlagName, azure.PersonalAccessTokenFlagName})
	Cmd.MarkPersistentFlagRequired(azure.OrganizationUrlFlagName.Name)
	Cmd.MarkPersistentFlagRequired(azure.PersonalAccessTokenFlagName.Name)
}
