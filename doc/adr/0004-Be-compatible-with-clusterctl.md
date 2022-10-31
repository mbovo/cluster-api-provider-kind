# 4. Be compatible with clusterctl

Date: 2022-10-26

## Status

Accepted

## Context

If we want to use the clusterctl tool to install the provider, we need to be compatible with the [provider contract](https://cluster-api.sigs.k8s.io/clusterctl/provider-contract.html)
This means adding

- metadata.yaml file with the supported version
- components.yaml file with all the required resources to be installed alongside the controller

## Decision

`components.yaml` is the same output of `make install` but redirect to file.
A new make target `make components` is added to generate it and goreleaser will add it to the release assets.

`metadata.yaml` is a new file with the following content:

```yaml
apiVersion: clusterctl.cluster.x-k8s.io/v1alpha3
releaseSeries:
    # these two fields are our infrastructure provider version
  - major: 0
    minor: 1
    contract: v1beta1 # must be the version of clusterctl used to install the provider
```

Additional `templates/cluster-template.yaml` is looked by `clusterctl` in the `generate` phase. The content must reflect what the KindCluster CR would look like if it was created by the user and a Cluster object that reference it.

```yaml
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "${CLUSTER_NAME}"
spec:
  clusterNetwork:
    pods:
      cidrBlocks: ["192.168.0.0/16"]
  infrastructureRef:
    kind: KindCluster
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    name: "${CLUSTER_NAME}"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: KindCluster
metadata:
  name: "${CLUSTER_NAME}"
spec:
  workerCount: ${WORKER_MACHINE_COUNT}
  controlPlaneCount: ${CONTROL_PLANE_MACHINE_COUNT}
  k8sVersion: ${KUBERNETES_VERSION}
```

The templates will be filled by clusterctl.
We took [this as example](https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/main/templates/cluster-template-eks.yaml)

## Consequences

Without clusterctl an end user will need to install the controller and create all resources by hand.
*NOTE*: KinD doesn't support to be executed inside a KinD cluster, so for testing we will need to use a different approach.
