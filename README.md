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

- `make dev-down`. This will shut down the project, shutting down the container,
deleting the cluster, and the telepresence daemon.

## Business

### web

#### v1

##### mid

We have our middleware functions contained here functions contained here. These
functions can take in arguments, and they return type `web.Middleware`, for
example, our Logging middleware takes the logger as a parameter. This is useful
so we can avoid having to use the singleton pattern for the logger.

To achieve this, we utilize a closure as can be seen below:

```go 

func Logger(log *zap.SugaredLogger) web.Middleware { 

    m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

            // Logging

			err := handler(ctx, w, r)

            // Some more logging

			return err
		}

		return h
	}

	return m
}

```


## Foundation

### web

We are using `package httptreemux` as our router, and in package web we create a
small web framework that extends the default funcitonality.     

The key entity is the `App` struct which is the entry point into our
application. We override the default `Handler` function, creating a new
signature which accepts `context.Context` as its first parameter. This allows us
to write our routes in such a way that we can pass in the context object. 

We also override the `Handle` function of the mux, allowing us to inject
middleware before and after the call to the provided `Handler` function. We
differentiate between two layers of middleware:


- *Local Layer Middleware*. This is for middleware that does not need to be
applied to all routes. For example, not all routes will require authentication.
This middleware is applied when we register a route via the `Handle` function.

- *Application Layer Middleware*. This gets applied to all routes attached to an
`App`. We use this for logging, error handling, metrics, and panic handling.

If we think of each route servicing incoming requests as an onion with the
`Handler` itself at the center, then we should think of the local layer
middleware as closest to the center, as it is applied first, with the
application layer middleware on the outside, as it is applied last.

## zarf 

### k8s 

#### dev

#### base

This is our base k8s configuration that never changes regardless of the
environment that we are running in. We use `kustomize` to patch the base
configuration into different environments (e.g., dev, staging, prod).


## Logging

We are using [package zap](https://github.com/uber-go/zap) for logging. We pass
around a single logger that is intialized in the `main()` function of the
server. This single logger is then passed around through the app using various
`cfg` structs.

## Configuration 

Our guiding principle is that all configuration must happen in `main.go`. This
is to enforce simplicity. When all configuration happens in one place, then we
don't have to worry about losing track of it.

