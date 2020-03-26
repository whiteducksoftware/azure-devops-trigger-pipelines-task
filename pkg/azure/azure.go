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
package azure

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelines"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FlagDefinition struct {
	Name        string
	Shorthand   string
	Default     string
	Description string
	Persistent  bool
}

var (
	OrganizationUrlFlagName FlagDefinition = FlagDefinition{
		Name:        "url",
		Description: "Azure DevOps Organization Url (must be https://dev.azure.com/ORG or https://ORG.visualstudio.com)",
		Persistent:  true,
	}
	PersonalAccessTokenFlagName FlagDefinition = FlagDefinition{
		Name:        "token",
		Description: "Azure DevOps Personal Access Token (PAT)",
		Persistent:  true,
	}
)

func AddFlags(cmd *cobra.Command, flags []FlagDefinition) {
	for _, flag := range flags {
		var flags *pflag.FlagSet
		if flag.Persistent {
			flags = cmd.PersistentFlags()
		} else {
			flags = cmd.Flags()
		}

		if flag.Shorthand != "" {
			flags.StringP(flag.Name, flag.Shorthand, flag.Default, flag.Description)
		} else {
			flags.String(flag.Name, flag.Default, flag.Description)
		}
	}
}

func GetDevOpsClient(flags *pflag.FlagSet) (*azuredevops.Connection, error) {
	orgURL, err := flags.GetString(OrganizationUrlFlagName.Name)
	if err != nil {
		return &azuredevops.Connection{}, err
	}

	pat, err := flags.GetString(PersonalAccessTokenFlagName.Name)
	if err != nil {
		return &azuredevops.Connection{}, err
	}

	return azuredevops.NewPatConnection(orgURL, pat), nil
}

func ListDevOpsPipelines(ctx context.Context, client pipelines.Client, project *string) ([]pipelines.Pipeline, error) {
	output := make([]pipelines.Pipeline, 0)

	var continuationToken string
	for ok := true; ok; ok = continuationToken != "" {
		result, err := client.ListPipelines(ctx, pipelines.ListPipelinesArgs{
			Project:           project,
			ContinuationToken: &continuationToken,
		})
		if err != nil {
			return nil, err
		}

		output = append(output, result.Value...)
		continuationToken = result.ContinuationToken
	}

	return output, nil
}
