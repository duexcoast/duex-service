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

### web/v1/mid

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


#### Logging Middleware

Simple enough. We log both before and after the call to `handler()`, when the
request begins and when the request completes.

Some values are taken from the `request` object, while others are stored in the
context at the foundation layer, such as the `traceID`.


#### Errors Middleware

This middleware accepts a logger as a parameter, because part of handling an
error is to log it. 

- The first course of action is to log the error.
- We use a switch-case block to determine if the error is trusted or untrusted.
We use this determination to fill up the `ErrorResponse` struct appropriately.
- We respond with the appropriate error message and status code.
- Finally, we check if the error was a shutdown error, allowing the app to
gracefully shutdown after having logged the error and responded to the request.

#### Panics Middleware 

This middleware calls `recover()` (within a deffered function, otherwise the
call will do nothing). If there is a `panic!` then we want to log it, we also
use of the setter functions in `package metrics` to safely update the panic
metrics value held in the context of the request.  

We want to return an error here, so we need to utilize *named return values*,
allowing us to update the value of the `err` variable from within the defer.
This error is then handled by the errors middleware, which will respond with a
`500 Internal Server Error`.

#### Metrics Middleware

This middleware sets the metrics data on the context of the request. After the
handler is called we increment requests, and increment goroutine metrics (which
uses a modulo operation to only do the work every 100 requests).

## Foundation

The foundation layer is meant to be usable across many applications, and other
developers. Each package in this layer should be as unopinionated as possible.
To make sure we maintain this high-degree of reusability, we have strong
policies in place for these packages:

- *No logging*. Logging should be determined by the user of these packages, we
don't want these packages to prescribe a specific logging package or
implementation pattern.

### web

We are using `package httptreemux` as our router, and in package web we create a
small web framework that extends the default funcitonality.     

The key entity is the `App` struct which is the entry point into our
application. We override the default `Handler` function, creating a new
signature which accepts the `context` as its first parameter.

```go
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
```

This allows us to write our routes in such a way that we can pass in the context
what object object.

We also override the `Handle` function of the mux, allowing us to inject
middleware before and after the call to the provided `Handler` function. We
differentiate between two layers of middleware:

- Application-layer middleware that is applied to all routes. This middleware is
passed to the `NewApp()` function. 

- Local middleware that is applied selectively. We supply this middleware in the
`Handle()` function.

#### Middleware Handled in Foundation Layer

While the majority of our middleware is implemented in the business-layer, we
are choosing to implement some middleware-like functionality in the foundation
layer, as part of our web framework. 

The business layer middleware is implemented in a mostly-typical, if not
somewhat unique pattern:

- Write a middleware function that can accept any type of argument, and which
returns a type of `web.Middleware`: This is the more familiar pattern: a
function that takes in a `Handler` and returns a `Handler`.

- The `wrapMiddleware()` function in our framework will apply all middlewares,
ultimately returning a single wrapped `Handler` super-charged with all of the
provided functionality.

Meanwhile, the functionality that we are implementing in the foundation layer
uses the `context` instead. In our `Handle()` function, we create a struct
containing a trace ID, the current time and store it in the `context`. We can
then make use of these in other middleware, by retrieving them from the request
context (e.g., as in the logger middleware).

#### Functions for Responses and Requests

`web.Respond`

 


#### Middleware Concerns
- *Local Layer Middleware*. This is for middleware that does not need to be
applied to all routes. For example, not all routes will require authentication.
This middleware is applied when we register a route via the `Handle` function.

- *Application Layer Middleware*. This gets applied to all routes attached to an
`App`. We use this for logging, error handling, metrics, and panic handling.

If we think of each route servicing incoming requests as an onion with the
`Handler` itself at the center, then we should think of the local layer
middleware as closest to the center, as it is applied first, with the
application layer middleware on the outside, as it is applied last.

#### context

Our `httptreemux` router will create a context for every incoming request, and
every outgoing response. Keeping this information in the foundation layer is a
debatable decision.

For incoming requests, we are keeping track of three things in the context:

- `TraceID`. A unique string that allows us to differentiate requests. We use
UUID for this.

- `Now`. We track the time and duration.

- `Status Code`. By keeping the status code in the context, we can use it in our
foundation layer `Respond()` function.

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

## Error handling

## Metrics

Our `metrics.go` file uses a singleton, `var m *metrics` is a package level
variable. In this instance, the API of the `expvar` package has forced our hand,
because it is using a singleton internally. We can get away with a package level
variable when we tick the following boxes:

- *The order of initialization does not matter*. We must not have an
order-dependency on our package level variables, our only need is that they are
intialized before `main()`.

- *They don't rely on anything from the configuration system*. If our package
level variable requires configuration, then it must be intialized in `main()`
(or `run()`, our nominal `main`).

- *The only code touching the variable is the code in the same source code file
as the variable*. Even if it's in the same package, we don't want code in other
files touching the package level variable. We won't to keep it centralized
because we want to have an easy mental model of how the variable is being
mutated.

We store our metrics data in the context of each request, updating individual
metrics on every request. Because we are using [package
expvar](https://pkg.go.dev/expvar), we can safely update these values
concurrently.

We make use of the `init()` function, which runs before `main()`, to initialize
the metrics singleton. We are keeping track of the number of:

- Goroutines
- Requests
- Errors
- Panics

Because we are storing metrics in the context, we want to use getters and
setters, because of the lack of type safety involved in the context.

The goroutine metric will slow down any given request because we have to make
use of the runtime package which will internally inspect the runtime state,
adding additional goroutines to do this work. We use a modulo operation on the
request value (`if v.requests.Value()%100 == 0`), so that we only have to
actually check how many goroutines are running on every 100th request (this
number is likely way too small for most services, but is just a placeholder).

### Categories of Errors

#### Untrusted Errors

This is the default. We do not want to expose information about the state of our
application beyond the boundaries of our API. This is a security risk. By
default, unless we can justify a reason otherwise, errors should be returned to
the client as `500 Internal Server Error`.

#### Trusted Errors 

We trust the messaging of the error to be secure and not leak any sensitive
information, such that we can respond to the client with the messaging.

We place this error type in the business layer, and we version it as well. This
is because our policies may change in future versions of our API in regards to
how we want to handle trusted errors.

#### Shutdown Error

If our service is having data integrity issues it should be shut down. When our
service is corrupting databases, file systems, etc., then we want to gracefully
shutdown. But we do not want any code outside *actually* shutting down our app,
and so instead these errors serve as a _suggestion_., our application can then
determine the appropriate action after inspecting the error.

We place this error type in the foundation layer. A Shutdown Error is a special
case of an Untrusted Error. Like all other errors, it gets handled by the
`errors` middleware, because it is untrusted it is logged, an an `Internal
Server Error` gets returned to the client. The middleware then checks if the
error is indeed a Shutdown Error, and if so, returns the error. This sends the
control flow back to the `(*App).Handle` function in the foundation layer.

At this point, we run a function `validateShutdown`. The reason we need this is
because it's possible that we may have network problems causing our code to
reach this line, but that are not errors worthy of shutting the app down.

When we do validate the shutdown error, we then use the `shutdown` channel in
the `App` struct to signal to the application that we need to gracefully
shutdown. Graceful shutdown will begin.


## Configuration 

Our guiding principle is that all configuration must happen in `main.go`. This
is to enforce simplicity. When all configuration happens in one place, then we
don't have to worry about losing track of it.

