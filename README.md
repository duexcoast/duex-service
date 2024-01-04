# duex-service

# Project Structure

## zarf 

### k8s 

#### dev

#### base

This is our base configuration that never changes regardless of the environment
that we are running in. We use `kustomize` to patch the base configuration into
different environments (e.g., dev, staging, prod).



