/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pipelines

import (
	"errors"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/cmd/pipelines/trigger"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/pkg/azure"
)

// pipelinesCmd represents the pipelines command
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

		if urlMatchA == false && urlMatchB == false {
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
	azure.AddFlags(Cmd, []azure.FlagDefinition{azure.OrganizationUrlFlagName, azure.PersonalAccessTokenFlagName}, true)
	Cmd.MarkPersistentFlagRequired(azure.OrganizationUrlFlagName.Name)
	Cmd.MarkPersistentFlagRequired(azure.PersonalAccessTokenFlagName.Name)
}
