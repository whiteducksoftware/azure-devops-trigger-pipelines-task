# Azure DevOps Pipelines trigger task

![Build and push Image](https://github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/workflows/Build%20and%20push%20Image/badge.svg)

## Description

Image which can be used in Azure DevOps Pipelines to trigger other pipelines.  

## Usage

```plaintext
Usage:
  devops-worker pipelines trigger [flags] -- PIPELINES...

Flags:
  -h, --help                   help for trigger
  -p, --project string         Name of the Azure DevOps Project
  -r, --targetRefName string   (Optional) Specify the GitRef on which the Pipeline should run
  -v, --targetVersion string   (Optional) Specify the Commit Hash on which the Pipeline should run
  -w, --waitForCompletion      (Optional) Specify if the task should block until the target pipeline is completed

Global Flags:
      --token string   Azure DevOps Personal Access Token (PAT) / Value of $(System.AccessToken)
      --url string     Azure DevOps Organization Url (must be https://dev.azure.com/ORG or https://ORG.visualstudio.com)
```

## Authentication

For authentication, you can use a Personal Access Token (PAT) or the Build-In `Build Service` User.  
In order to use the `Build Service` User, you need to manually grant the User permission to queue builds on the pipelines, then you can use as token value the value of the `System.AccessToken`  built-in variable. E.g. `--token $(System.AccessToken)`.

## Example

```yaml
pool:
  vmImage: 'ubuntu-latest'

container:
  image: whiteduck/azure-devops-pipeline-trigger-task:2.1.3

steps:
  - script: |
      devops-worker pipelines trigger \
        --url "$(System.CollectionUri)" \
        --token "$(System.AccessToken)" \
        --project "$(System.TeamProject)" \
        --targetRefName "$(Build.SourceBranchName)" \
        --targetVersion "$(Build.SourceVersion)" \
        --waitForCompletion \
        -- "My Demo Pipeline"
    displayName: "Trigger 'My Demo Pipeline' pipeline"
```
