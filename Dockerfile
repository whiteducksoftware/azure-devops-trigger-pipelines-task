FROM ubuntu:16.04

# Install azure-cli
RUN apt update && \
    apt install -y --no-install-recommends \
      ca-certificates \
      curl \
      apt-transport-https \
      lsb-release \
      gnupg \
      jq

RUN curl -sL https://packages.microsoft.com/keys/microsoft.asc | \
    gpg --dearmor | \
    tee /etc/apt/trusted.gpg.d/microsoft.asc.gpg > /dev/null

RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" | \
    tee /etc/apt/sources.list.d/azure-cli.list

RUN apt update && \
    apt install -y --no-install-recommends azure-cli

# Add lables
LABEL version="1.0.0"
LABEL maintainer="whiteduck GmbH" 
LABEL name="Trigger Azure Pipelines" 
LABEL description="Container which can trigger Azure pipeline(s)" 

# Add scripts
RUN mkdir -p /opt/azure/pipelines/
ADD scripts /opt/azure/pipelines/

# Fix permissions of the scripts
RUN chmod +x /opt/azure/pipelines/init.sh
RUN chmod +x /opt/azure/pipelines/queue.sh