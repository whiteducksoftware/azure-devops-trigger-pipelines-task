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
package trigger

import (
	"errors"
	"fmt"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelines"
	"github.com/spf13/cobra"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task-go/pkg/azure"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task-go/pkg/utils"
)

var (
	ProjectFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "project",
		Shorthand:   "p",
		Description: "",
	}
	RepositoryNameFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "repo",
		Description: "",
	}
	TargetRefNameFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "targetRefName",
		Description: "",
	}
	TargetVersionFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "targetVersion",
		Description: "",
	}
)

var Cmd = &cobra.Command{
	Use:   "trigger",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("No pipeline names have been passed")
		}

		flags := cmd.Flags()

		project, err := flags.GetString(ProjectFlag.Name)
		if err != nil {
			return err
		}

		connection, err := azure.GetDevOpsClient(flags)
		if err != nil {
			return err
		}

		pipelinesClient := pipelines.NewClient(cmd.Context(), connection)
		result, err := azure.ListDevOpsPipelines(cmd.Context(), pipelinesClient, &project)
		if err != nil {
			return err
		}

		for _, pipeline := range result {
			if utils.StringInSlice(*pipeline.Name, args) {
				// ToDO: These parameters seems to get ignored by the DevOps Api
				// https://github.com/microsoft/azure-devops-go-api/issues/55
				var repositoryResourceParameters map[string]pipelines.RepositoryResourceParameters
				if utils.IsFlagPassed(RepositoryNameFlag.Name, flags) && utils.IsFlagPassed(TargetRefNameFlag.Name, flags) && utils.IsFlagPassed(TargetVersionFlag.Name, flags) {
					repoName, err := flags.GetString(RepositoryNameFlag.Name)
					if err != nil {
						return err
					}

					targetRefName, err := flags.GetString(TargetRefNameFlag.Name)
					if err != nil {
						return err
					}

					targetVersion, err := flags.GetString(TargetVersionFlag.Name)
					if err != nil {
						return err
					}

					repositoryResourceParameters = map[string]pipelines.RepositoryResourceParameters{
						repoName: pipelines.RepositoryResourceParameters{
							RefName: to.StringPtr(targetRefName),
							Version: to.StringPtr(targetVersion),
						},
					}
				}

				runResult, err := pipelinesClient.RunPipeline(cmd.Context(), pipelines.RunPipelineArgs{
					Project:    &project,
					PipelineId: pipeline.Id,
					RunParameters: &pipelines.RunPipelineParameters{
						Resources: &pipelines.RunResourcesParameters{
							Repositories: &repositoryResourceParameters,
						},
					},
				})
				if err != nil {
					return err
				}

				if *runResult.State != pipelines.RunStateValues.InProgress && *runResult.State != pipelines.RunStateValues.Completed {
					return fmt.Errorf("Unkown error occured, result: %s, state: %s", string(*runResult.Result), string(*runResult.State))
				}
			}
		}

		return nil
	},
}

func init() {
	azure.AddFlags(Cmd, []azure.FlagDefinition{ProjectFlag, RepositoryNameFlag, TargetRefNameFlag, TargetVersionFlag}, false)
	Cmd.MarkFlagRequired(ProjectFlag.Name)
}
