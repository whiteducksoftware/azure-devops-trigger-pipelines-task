# Azure DevOps Pipelines trigger task
![Build and push Image](https://github.com/whiteducksoftware/azure-devops-trigger-pipelines-task/workflows/Build%20and%20push%20Image/badge.svg)

### Description
Image which can be used in Azure DevOps Pipelines to trigger other pipelines.  

### Variables
- `URL` - **Mandatory**; the fully-qualified URL to the Azure DevOps organization (eg, `https://dev.azure.com/organization` or `https://server.example.com:8080/tfs/DefaultCollection`)
- `Token` – **Mandatory**  
- `Project` – **Mandatory** 
- `AZURE_PIPELINE_NAME`; Note: Name **is** case-sensitive.


### Example:
Single Pipeline:
```yaml
pool:
  vmImage: 'ubuntu-16.04'

container: 
  image: whiteduck/azure-devops-pipeline-trigger-task:v2.0

variables:
  AZURE_DEVOPS_URL: "https://dev.azure.com/demo"
  AZURE_DEVOPS_PROJECT: "demo"
  AZURE_PIPELINE_NAME: "My Demo Pipeline"
  AZURE_DEVOPS_BRANCH: $(Build.SourceBranchName)
  AZURE_DEVOPS_COMMIT: $(Build.SourceVersion)

steps:
  - script: devops-worker pipelines trigger --url $(AZURE_DEVOPS_URL) --token $(AZURE_DEVOPS_TOKEN) --project $(AZURE_DEVOPS_PROJECT) -- $(AZURE_PIPELINE_NAME)
    displayName: "Trigger 'My Demo Pipeline' pipeline"
    env:
      AZURE_DEVOPS_TOKEN: $(AZURE_DEVOPS_TOKEN)
```

Multiple Pipelines:
```yaml
pool:
  vmImage: 'ubuntu-16.04'

container: 
  image: whiteduck/azure-devops-pipeline-trigger-task:v2.0

variables:
  AZURE_DEVOPS_URL: "https://dev.azure.com/demo"
  AZURE_DEVOPS_PROJECT: "demo"
  AZURE_DEVOPS_BRANCH: $(Build.SourceBranchName)
  AZURE_DEVOPS_COMMIT: $(Build.SourceVersion)

steps:
  - script: devops-worker pipelines trigger --url $(AZURE_DEVOPS_URL) --token $(AZURE_DEVOPS_TOKEN) --project $(AZURE_DEVOPS_PROJECT) -- "PipelineA" "PipelineB" "PipelineC"
    displayName: "Trigger 'PipelineA, PipelineB, PipelineC' pipeline"
    env:
      AZURE_DEVOPS_TOKEN: $(AZURE_DEVOPS_TOKEN)
...
```