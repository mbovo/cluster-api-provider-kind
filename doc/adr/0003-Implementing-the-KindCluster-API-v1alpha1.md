# 3. Implementing the kindcluster api v1alpha1

Date: 2022-10-25

## Status

Accepted
Amended by [5. Kindcluster api v1beta1](0005-KindCluster-API-v1beta1.md)

## Context

Implement the KindCluster API version `v1alpha1`

## Decision

The CustomeResource will have the following structure:

```yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: KindCluster
metadata: Kubernetes Metadata Object
spec:
  workerCount: 0 #int32, optional  number of workers in the kind cluster target
  controlPlaneCount: 1 #int32, optional number of control plane nodes in the kind cluster target
  k8sVersion: #string, optional version of kubernetes to use in the kind cluster target
  controlPlaneEndpoint: #clusterv1.APIEndpoint, optional
  image: #string, optional  image to use for the kind cluster target (default: kindest/node)
status:
  ready: #bool, optional, if the cluster is ready or not
  conditions: #[]clusterv1.Condition, optional
```

## Consequences

Let's start for now with a simpler structure.
There it lacks most of the optional fields in the status to verify errors and signal conditions, will be added later on.
All the fields except `controlPlaneEndpoint` are passed to kind in order to create a new cluster
The KinD cluster name will be read from the `metadata.name` field.
