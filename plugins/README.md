# Plugins for Aruba Cloud Provider KOG

This repository contains the source code of a set of plugins for the Aruba Cloud Provider KOG.

## Testing guide

For detailed instructions on building and testing the plugins, please refer to the [Testing Guide](./docs/testing.md).

## Build Instructions

This project is a Go workspace-based monorepo containing multiple, independent plugins. The build system is centralized in this directory (`plugins/`).

### Building with ko

The primary way to build and publish the container images is using Google's `ko` tool. The configuration is in the `.ko.yaml` file in this directory. This is useful for local development.

The `.ko.yaml` file defines a unique image name for each plugin. To build and publish all plugins, simply run:

```sh
ko publish .
```

`ko` will read the `.ko.yaml` file, build each plugin specified, and push them to the container registry defined in your `KO_DOCKER_REPO` environment variable.

Example published images:
- `KO_DOCKER_REPO`/subnet-plugin

### Building with Docker

A generic, multi-stage `Dockerfile` is located in this directory. It is optimized for granular layer caching in the Go workspace, which means that changes to one plugin will not invalidate the build cache for other, unrelated plugins.

It can be used to build a container image for any specific plugin by passing a build argument.

**To build the `subnet-plugin`:**
```sh
docker build --build-arg PLUGIN_NAME=subnet-plugin -t subnet-plugin:latest .
```

This command should be run from this `plugins` directory, as it sets the necessary build context to access the shared `pkg` directory.
Note that the `PLUGIN_NAME` argument must match the name of the subdirectory containing the `main.go` file for the desired plugin therefore it is case-sensitive and must be exact.

Note that usually this process is automated as part of the [release process](../docs/release.md).
