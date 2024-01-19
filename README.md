# duex-service

# Project Structure

## Makefile

To get the project up and running:
- `make dev-up`. This will get our kubernetes cluster up and running.

- `make dev-apply`. Apply changes made to the yaml defining our kubernetes
cluster

- `make dev-update-apply`. This will build our docker image, use kind to load
the docker image into the k8s cluster, use kustomize to load our base k8s
config, and then wait until the cluster has applied the changes. 

Applying code changes to the kubernetes cluster: `make dev-update`

To shutdown the project: `make dev-down`

## zarf 

### k8s 

#### dev

#### base

This is our base configuration that never changes regardless of the environment
that we are running in. We use `kustomize` to patch the base configuration into
different environments (e.g., dev, staging, prod).


## Logging

We use a single logger.

## Configuration 

Our guiding principle is that all configuration must happen in `main.go`. This
is to enforce simplicity. When all configuration happens in one place, then we
don't have to worry about losing track of it.

