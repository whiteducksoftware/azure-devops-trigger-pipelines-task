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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelines"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/pkg/azure"
	"github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/pkg/utils"
)

var (
	projectFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "project",
		Shorthand:   "p",
		Default:     "",
		Description: "",
	}
	targetRefNameFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "targetRefName",
		Default:     "",
		Description: "",
	}
	targetVersionFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "targetVersion",
		Default:     "",
		Description: "",
	}
	waitForCompletionFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "waitForCompletion",
		Shorthand:   "w",
		Default:     false,
		Description: "",
	}
)

// Cmd represents the trigger command
var Cmd = &cobra.Command{
	Use:   "trigger",
	Short: "Triggers the specified pipeline(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("No pipeline names have been passed")
		}

		flags := cmd.Flags()

		project, err := flags.GetString(projectFlag.Name)
		if err != nil {
			return err
		}

		wait, err := flags.GetBool(waitForCompletionFlag.Name)
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

		var repositoryResourceParameters map[string]pipelines.RepositoryResourceParameters
		if utils.IsFlagPassed(targetRefNameFlag.Name, flags) && utils.IsFlagPassed(targetVersionFlag.Name, flags) {
			targetRefName, err := flags.GetString(targetRefNameFlag.Name)
			if err != nil {
				return err
			}

			targetVersion, err := flags.GetString(targetVersionFlag.Name)
			if err != nil {
				return err
			}

			repositoryResourceParameters = map[string]pipelines.RepositoryResourceParameters{
				"self": {
					RefName: to.StringPtr(targetRefName),
					Version: to.StringPtr(targetVersion),
				},
			}
		}

		for _, pipeline := range result {
			if utils.StringInSlice(*pipeline.Name, args) {
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

				log.Info(*runResult.Url)

				if wait {
					for {
						run, err := pipelinesClient.GetRun(cmd.Context(), pipelines.GetRunArgs{
							Project:    &project,
							PipelineId: pipeline.Id,
							RunId:      runResult.Id,
						})
						if err != nil {
							return err
						}

						if run.Result != nil {
							if *run.Result == pipelines.RunResultValues.Succeeded {
								json, err := json.Marshal(*run)
								if err != nil {
									log.Errorln("Failed to serialize result into json", err)
								} else {
									log.Info(string(json))
								}
								break
							}

							if *run.Result == pipelines.RunResultValues.Failed {
								return fmt.Errorf("Pipeline %s failed", *pipeline.Name)
							}
						}

						time.Sleep(1 * time.Second)
					}
				} else {
					json, err := json.Marshal(*runResult)
					if err != nil {
						log.Errorln("Failed to serialize result into json", err)
					} else {
						log.Info(string(json))
					}
				}
			}
		}

		return nil
	},
}

func init() {
	azure.AddFlags(Cmd, []azure.FlagDefinition{projectFlag, targetRefNameFlag, targetVersionFlag, waitForCompletionFlag})
	Cmd.MarkFlagRequired(projectFlag.Name)
}
