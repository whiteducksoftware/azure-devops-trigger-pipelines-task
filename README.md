# Azure DevOps Pipelines trigger task

![Build and push Image](https://github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/workflows/Build%20and%20push%20Image/badge.svg)

## Description

Image which can be used in Azure DevOps Pipelines to trigger other pipelines.  
## Usage

```
Usage:
  devops-worker pipelines trigger [flags]

Flags:
  -h, --help                   help for trigger
  -p, --project string         Name of the Azure DevOps Project
      --targetRefName string   (Optional) Specify the GitRef on which the Pipeline should run
      --targetVersion string   (Optional) Specify the Commit Hash on which the Pipeline should run
  -w, --waitForCompletion      Specify if the task should block until the target pipeline is completed

Global Flags:
      --token string   Azure DevOps Personal Access Token (PAT)
      --url string     Azure DevOps Organization Url (must be https://dev.azure.com/ORG or https://ORG.visualstudio.com)
```

## Authentication

For authentication, you can use a Personal Access Token (PAT) or the Build-In `Build Service` User.  
In order to use the `Build Service` User, you need to manually grant the User permission to queue builds on the pipelines, then you can use as token value the value of the `System.AccessToken`  built-in variable. E.g. `--token $(System.AccessToken)`.

## Example:

ToDo
