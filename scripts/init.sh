#!/bin/bash
set -e 

echo "Preparing azure cli..."

az extension add -n azure-devops

if [ -z "$AZURE_DEVOPS_URL" ]; then
    echo "AZURE_DEVOPS_URL is not set." >&2
    exit 1
fi

if [ -z "$AZURE_DEVOPS_PROJECT" ]; then
    echo "AZURE_DEVOPS_PROJECT is not set." >&2
    exit 1
fi

if [ -z "$AZURE_DEVOPS_TOKEN" ]; then
    echo "AZURE_DEVOPS_TOKEN is not set." >&2
    exit 1
fi

az devops configure --defaults organization="${AZURE_DEVOPS_URL}" project="${AZURE_DEVOPS_PROJECT}"
    
echo "${AZURE_DEVOPS_TOKEN}" | az devops login --organization "${AZURE_DEVOPS_URL}"

echo "Done."