# 2. Implementing cluster api infrastructure provider

Date: 2022-10-35

## Status

Accepted

## Context

This project aim to implement a cluster-api [infrascture provider](https://cluster-api.sigs.k8s.io/developer/providers/cluster-infrastructure.html) for [KinD](https://kind.sigs.k8s.io/)

Requirements are:

- when a KindCluster custom resource is submitted, the provider should ensure that a matching KinD cluster is created
- The custom resources should specify the KinD cluster name and a ready boolean in the status which should be set to `true` when the KinD cluster is ready.
- When the Custom Resource is removed the provider should ensure that the KinD cluster is deleted.
- Utilize finalizrs to ensure that the no resources are left behind
- (OPTIONAL) Implements as much as [Infrastructure Prodiver Specification](https://cluster-api.sigs.k8s.io/developer/providers/cluster-infrastructure.html) as possible

## Decision

Going to:

- use [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) to scaffold a kubernetes controller
- use the following standards:
  - [Provider contracts](https://cluster-api.sigs.k8s.io/developer/providers/contracts.html) in the controller yaml manifests
  - [Naming](https://cluster-api.sigs.k8s.io/developer/providers/implementers-guide/naming.html)
  - [CRDs](https://cluster-api.sigs.k8s.io/developer/providers/implementers-guide/generate_crds.html)
- using [goreleaser](https://goreleaser.com/) to help build and release the provider
- using github actions for the CI pipeline

## Consequences

--
