# 8. Create loggerwrapper

Date: 2022-11-03

## Status

Accepted

## Context

Kind lib is using its own logger, it is configurable but is different from the one used by controller-runtime.

## Decision

Create a logger wrapper that exposes the same interface as the kind logger but uses the controller-runtime logger to log messages

## Consequences

Log messages from kind will be totally trasnparent and we will have an idea of the various phases of kind creation and deletion on the logs
