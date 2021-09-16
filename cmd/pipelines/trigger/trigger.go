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
		Description: "Name of the Azure DevOps Project",
	}
	targetRefNameFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "targetRefName",
		Shorthand:   "r",
		Default:     "",
		Description: "(Optional) Specify the GitRef on which the Pipeline should run",
	}
	targetVersionFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "targetVersion",
		Shorthand:   "v",
		Default:     "",
		Description: "(Optional) Specify the Commit Hash on which the Pipeline should run",
	}
	waitForCompletionFlag azure.FlagDefinition = azure.FlagDefinition{
		Name:        "waitForCompletion",
		Shorthand:   "w",
		Default:     false,
		Description: "(Optional) Specify if the task should block until the target pipeline is completed",
	}
)

// Cmd represents the trigger command
var Cmd = &cobra.Command{
	Use:   "trigger",
	Short: "Triggers the specified pipeline(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no pipeline names have been passed")
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
		result, err := azure.GetDevOpsPiplines(cmd.Context(), pipelinesClient, project, args)
		if err != nil {
			return err
		}

		repositoryResourceParameters := pipelines.RepositoryResourceParameters{}
		if utils.IsFlagPassed(targetVersionFlag.Name, flags) {
			targetVersion, err := flags.GetString(targetVersionFlag.Name)
			if err != nil {
				return err
			}

			repositoryResourceParameters.Version = to.StringPtr(targetVersion)
		}
		if utils.IsFlagPassed(targetRefNameFlag.Name, flags) {
			targetRefName, err := flags.GetString(targetRefNameFlag.Name)
			if err != nil {
				return err
			}

			repositoryResourceParameters.RefName = to.StringPtr(targetRefName)
		}

		for _, pipeline := range result {
			runResult, err := pipelinesClient.RunPipeline(cmd.Context(), pipelines.RunPipelineArgs{
				Project:    &project,
				PipelineId: pipeline.Id,
				RunParameters: &pipelines.RunPipelineParameters{
					Resources: &pipelines.RunResourcesParameters{
						Repositories: &map[string]pipelines.RepositoryResourceParameters{"self": repositoryResourceParameters},
					},
				},
			})
			if err != nil {
				return err
			}

			if *runResult.State != pipelines.RunStateValues.InProgress && *runResult.State != pipelines.RunStateValues.Completed {
				log.Errorf("unkown error occured, result: %s, state: %s", string(*runResult.Result), string(*runResult.State))
			}

			log.Info(*runResult.Url)

			if wait {
			__Completion_Check_Loop: // labeled for better code readablilty
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
						switch *run.Result {
						case pipelines.RunResultValues.Succeeded:
							json, err := json.Marshal(*run)
							if err != nil {
								log.Errorf("failed to serialize result into json, %s", err.Error())
							} else {
								log.Info(string(json))
							}
						case pipelines.RunResultValues.Failed:
							log.Errorf("pipeline %s failed", *pipeline.Name)
						case pipelines.RunResultValues.Canceled:
							log.Warnf("pipeline %s was canceled", *pipeline.Name)
						default:
							log.Errorf("pipeline %s, unknown error occured, result: %s", *pipeline.Name, string(*run.Result))
						}

						break __Completion_Check_Loop
					}

					time.Sleep(10 * time.Second)
				}
			} else {
				json, err := json.Marshal(*runResult)
				if err != nil {
					log.Errorf("pipeline %s, failed to serialize result into json, %s", *&pipeline.Name, err.Error())
				} else {
					log.Info(string(json))
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
