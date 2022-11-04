# cluster-api-provider-kind

Cluster API provider for KinD or CAPK (pronounced cap-cake) is a provider for creating Kubernetes clusters using KinD and Cluster API.

## Description

This is an infrastructure provider for CAPI using `kind` to generate the cluster. It is intended to be used as an example and for testing purposes

In order to test the provider inside another kind cluster you need to hack a little the kind "bootstrap" cluster.


## Getting Started

Create a boostrap cluster with kind hacked with docker(with `examples/bootstrap-cluster.yaml`:

```bash
kind create cluster --config examples/bootstrap-cluster.yaml
```

Configure clusterctl to use the repo of this cluster-api-provider-kind:

```bash
$ cat ~/.cluster-api/clusterctl.yaml

providers:
  - name: "kind"
    url: "https://github.com/mbovo/cluster-api-provider-kind/releases/latest/infrastructure-components.yaml"
    type: "InfrastructureProvider"
```

Boostrap cluster-api:

`KIND_DOCKER_FSGROUP`set the securityContext fsGroup to make docker works inside the kind node. **Must be set to the group used by docker on the underneath host**.
You can get it with `getent group docker` on the underneath host.(e.g. your linux host or the VM where docker runs on MacOS)
The default (101) is the one used by [coLima](https://github.com/abiosoft/colima) in macOs.

```bash

```bash
export KIND_DOCKER_FSGROUP=101
clusterctl init --core cluster-api  --infrastructure kind -v5
```

Generate KindCluster:

```bash
clusterctl generate cluster my-cluster --kubernetes-version v1.24.7 --control-plane-machine-count=1 --worker-machine-count=1 -i kind:v0.1.5 | kubectl apply -f -
```

## Improvements

- [ ] Add errors to the status as per [Infrastucture provider contract](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/cluster.html#infrastructure-provider)
- [ ] Export created cluster Kubeconfig to the end user
- [ ] Automatically retrieve kind images given a Kubernetes version
- [ ] Add a [ControlPlane provider](https://cluster-api.sigs.k8s.io/developer/architecture/controllers/control-plane.html) to configure the controlplane and stopping Kind before starting Kubernetes (using create option `CreateWithStopBeforeSettingUpKubernetes`)

## Pitfails

### Github rate limiting

When downloading manifest with clusterctl you can hit the github rate limit. You can set your `GITHUB_TOKEN` env variable to avoid this.

## License

Copyright 2022 Manuel Bovo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
