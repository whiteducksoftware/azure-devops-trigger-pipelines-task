#!/bin/bash
if [ -z ${AZURE_PIPELINE_NAME+x} ]; then
    AZURE_PIPELINE_NAME="$1";
fi

if [ -z "$AZURE_PIPELINE_NAME" ]; then
    echo "AZURE_PIPELINE_NAME is not set." >&2
    exit 1
fi

PIPELINES=$(az pipelines build definition list --name "${AZURE_PIPELINE_NAME}" --output json)

if ! (echo "${PIPELINES}" | jq -e .); then
    echo "Failed to fetch pipelines. Error: ${PIPELINES}"
    exit 1;
fi 

COUNT=$(echo "${PIPELINES}" | jq length)

if [ "$COUNT" -eq 0 ]; 
then
   echo "No pipeline found with name: '${AZURE_PIPELINE_NAME}'". >&2
   exit 1;
fi

if [ "$COUNT" -gt 1 ]; 
then
    echo "Multple pipelines were found with name: '${AZURE_PIPELINE_NAME}'. Pass unique pipeline name and try again." >&2
    exit 1;
fi

BUILD_DEFINITION_ID=$(echo "${PIPELINES}" | jq -r ".[0]?.id //empty")
BUILD_DEFINITION=$(az pipelines build definition show --id "${BUILD_DEFINITION_ID}" --output json)

if ! (echo "${BUILD_DEFINITION}" | jq -e .); then
    echo "Failed to  get pipeline using Id: ${BUILD_DEFINITION_ID}. Error: ${BUILD_DEFINITION}"
    exit 1;
fi

REPOSITORY_NAME=$(echo "${BUILD_DEFINITION}" | jq -r ".repository?.name? //empty")
REPOSITORY_TYPE=$(echo "${BUILD_DEFINITION}" | jq  -r ".repository?.type?  //empty")

if [ -n "$REPOSITORY_NAME" ] && [ -n "$REPOSITORY_TYPE" ] && [ "$REPOSITORY_TYPE" = "TfsGit" ]; 
then
    BUILD_OUTPUT=$(az pipelines build queue --definition-name "${AZURE_PIPELINE_NAME}" --branch "${AZURE_DEVOPS_BRANCH}" --commit-id "${AZURE_DEVOPS_COMMIT}" --output json)
else
    BUILD_OUTPUT=$(az pipelines build queue --definition-name "${AZURE_PIPELINE_NAME}" --output json)
fi

if [ -z "$BUILD_OUTPUT" ];
then
    echo "Failed to queue build."
	exit 1;
fi

if ! (echo "${BUILD_OUTPUT}" | jq -e .); then
    echo "Failed to queue pipeline. Error: ${BUILD_OUTPUT}"
    exit 1;
fi

echo "${BUILD_OUTPUT}"
ERROR="error"
ERROR_COUNT=$(echo "${BUILD_OUTPUT}" | jq -r ".validationResults[]? | select(.result?==\"$ERROR\") | .result" | wc -l )
if [ "$ERROR_COUNT" -gt 0 ];
then
    echo "Failed to queue pipeline." >&2
    MESSAGE=$(echo "${BUILD_OUTPUT}" | jq -r ".validationResults[]?| select(.result?==\"$ERROR\")| .message //empty" | tr '\n' ','| sed  's/.$/./' )
    if [ -n "$MESSAGE" ];
    then
        echo "Validation result contains error. Error: ${MESSAGE}" >&2
    fi
    exit 1;
fi
