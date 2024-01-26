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

- `make dev-update`. If all we did is change some code, this is all we need

To shutdown the project: `make dev-down`

## Foundation

### web

We are using `package httptreemux` as our router, and in package web we create a
small web framework that extends the default funcitonality.     

The key entity is the `App` struct which is the entry point into our
application. We override the default `Handler` function, creating a new
signature which accepts `context.Context` as its first parameter. This allows us
to write our routes in such a way that we can pass in the context object. 

We also override the `Handle` function of the mux, allowing us to inject
middleware before and after the call to the provided `Handler` function.

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

