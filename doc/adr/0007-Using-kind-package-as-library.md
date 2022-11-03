# 7. Using kind package as library

Date: 2022-11-01

## Status

Supercedes [6. Implementing reconcile loop](0006-Implementing-Reconcile-Loop.md)

Accepted

## Context

To drop `kind` binary dependency and the ugly `exec` call we are using kind as libary

## Decision

The KinD public APIs are still in alpha, so we will use the `sigs.k8s.io/kind/pkg/cluster` package to create and delete clusters.
To configure the cluster the package used will be `"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"`.
To implement this logic we levarage a `KindHelper` moving the logic from the controller to a separate package.
The old logic will be deleted.

## Consequences

--
