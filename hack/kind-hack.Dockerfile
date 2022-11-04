FROM kindest/node:v1.24.7

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates curl gnupg lsb-release && \
  mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg && \
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" > /etc/apt/sources.list.d/docker.list && \
  apt-get update && apt-get install -y docker-ce-cli && rm -rf /var/lib/apt/lists/* && \
  docker context create underneath --description "docker on underneath host" --docker host=unix:///mnt/host/docker.sock && \
  docker context use underneath
