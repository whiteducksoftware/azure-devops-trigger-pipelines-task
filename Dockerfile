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
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="Trigger Azure Pipelines" 
LABEL org.label-schema.description="Container which can trigger Azure pipeline(s)" 
LABEL org.label-schema.vcs-ref="https://github.com/whiteducksoftware/azure-devops-trigger-pipelines-task"
LABEL org.label-schema.maintainer="Stefan Kürzeder <stefan.kuerzeder@whiteduck.de>"

# Add scripts
RUN mkdir -p /opt/azure/pipelines/bin
ADD scripts /opt/azure/pipelines/

# Fix permissions of the scripts
RUN chmod +x /opt/azure/pipelines/init.sh
RUN chmod +x /opt/azure/pipelines/queue.sh

# Create bin folder
RUN ln -s /opt/azure/pipelines/init.sh /opt/azure/pipelines/bin/task_init
RUN ln -s /opt/azure/pipelines/queue.sh /opt/azure/pipelines/bin/task_queue

# Add to path
ENV PATH="/opt/azure/pipelines/bin:${PATH}"
