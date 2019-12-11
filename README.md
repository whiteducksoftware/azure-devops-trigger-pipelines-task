# Azure DevOps Pipelines trigger task

### Description
Image which can be used in Azure DevOps Pipelines to trigger other pipelines.  
This project was inspired by [Azure/github-actions/pipelines](https://github.com/Azure/github-actions/tree/master/pipelines).

### Secrets
- `AZURE_DEVOPS_TOKEN` – **Mandatory** 


### Environment variables
- `AZURE_DEVOPS_URL` – **Mandatory**; the fully-qualified URL to the Azure DevOps organization (eg, `https://dev.azure.com/organization` or `https://server.example.com:8080/tfs/DefaultCollection`)
- `AZURE_DEVOPS_PROJECT` – **Mandatory** 
- `AZURE_PIPELINE_NAME` – Optional; If not set, you need to pass it to task_queue: `task_queue *NAME*`  
    Note: `AZURE_PIPELINE_NAME` or `NAME` **is** case-sensitive.

### Notes:
`task_init` is **mandatory** to be called before using `task_queue`

### Example:
Single Pipeline:
```yaml
pool:
  vmImage: 'ubuntu-16.04'

container: 
  image: whiteduck/azure-devops-pipeline-trigger-task

variables:
  AZURE_DEVOPS_URL: "https://dev.azure.com/demo"
  AZURE_DEVOPS_PROJECT: "demo"
  AZURE_PIPELINE_NAME: "My Demo Pipeline"

steps:
  - script: task_init
    displayName: "Initialize the task"

  - script: task_queue
    displayName: "Trigger 'My Demo Pipeline' pipeline"
```

Multiple Pipelines:
```yaml
pool:
  vmImage: 'ubuntu-16.04'

container: 
  image: whiteduck/azure-devops-pipeline-trigger-task

variables:
  AZURE_DEVOPS_URL: "https://dev.azure.com/demo"
  AZURE_DEVOPS_PROJECT: "demo"

steps:
  - script: task_init
    displayName: "Initialize the task"

  - script: task_queue "PipelineA"
    displayName: "Trigger 'PipelineA' pipeline"

  - script: task_queue "PipelineB"
    displayName: "Trigger 'PipelineB' pipeline"

  - script: task_queue "PipelineC"
    displayName: "Trigger 'PipelineC' pipeline"

...
```