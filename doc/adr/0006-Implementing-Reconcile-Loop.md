# 6. Implementing reconcile loop

Date: 2022-10-31

## Status

Superceded by [7. Using kind package as library](0007-Using-kind-package-as-library.md)

## Context

The controller needs a reconcile Loop for the KindCluster custom resources

## Decision

The first iteration will be a simple reconcile loop using `kind` binary directly.
The controller will call the binary to create and delete cluster with specified options

## Consequences

This first iteration add a depenedency on the `kind` binary, which is not ideal, the next iteration will use kind package as a library.
In both cases the controller image will depend on `docker` binary to be available.
