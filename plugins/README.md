# Plugins for Aruba Cloud Provider KOG

This repository contains the source code of a set of plugins for the Aruba Cloud Provider KOG.
Note: currently only one plugin is implemented: the `subnet-plugin` even if the structure allows to easily add more plugins in the future if needed.z

Specialized web services that address some integration issues.
They are designed to work with the [`rest-dynamic-controller`](https://github.com/krateoplatformops/rest-dynamic-controller/).

## Summary

- [Project structure](#project-structure)
- [Subnet plugins](#subnet-plugins)
    - [Get Subnet endpoint](#get-subnet-endpoint)
    - [Create Subnet endpoint](#create-subnet-endpoint)
    - [Update Subnet endpoint](#update-subnet-endpoint)
    - [List Subnets endpoint](#list-subnets-endpoint)
- [Authentication](#authentication)
- [Documentation](#documentation)
- [Testing guide](#testing-guide)
- [Build Instructions](#build-instructions)
  - [Building with ko](#building-with-ko)
  - [Building with Docker](#building-with-docker)

## Project structure

The project is organized as a Go workspace-based monorepo containing multiple, independent plugins. Each plugin resides in its own subdirectory under the `cmd/` directory, while shared code is located in the `pkg/` directory.
Each plugin is a standalone Go module, allowing for independent building and versioning.
Therefore, each plugin can be built into its own container image. Please refer to the [release process section of the main readme](../README.md#release-process) for more details.

## Subnet plugins

### Get Subnet endpoint

**Description**:
This endpoint retrieves a specific subnet by its ID in the specified Aruba Cloud Project and VPC.
It returns the subnet details, including its ID, name, and other metadata and properties.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The endpoint exists to flatten the `metadata` field used in the request and response bodies of the subnet resource in the Aruba Cloud API. This is necessary since the field `metadata.name` is used as resource **identifer** for the subnet by `oasgen-provider` and as a consequence it is put in the status of the generated CRD. Due to a current limitation of the underlying CRD generation library, **nested fields used as identifiers are not fully supported**.
- Therefore the plugin will accept a request body where the entire `metadata` object is flattened. In particular, the request body will have a top-level `name` field (among others), and it will internally map it to the `metadata.name` field expected by the Aruba Cloud API.
- Similarly, after receiving the response from the Aruba Cloud API, the plugin will flatten the `metadata` field in the response body before returning it to the client (client == `rest-dynamic-controller` in this case).

</details>

<details>
<summary><b>Request</b></summary>
<br/>

```http
GET /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets/{id}
```

**Path parameters**:
- `projectId` (string, required): The ID of the Aruba Cloud project.
- `vpcId` (string, required): The ID of the VPC.
- `id` (string, required): The ID of the subnet to retrieve.

**Query parameters**:
- `api-version` (string, required): The version of the Aruba Cloud API to use. For example, `1.0`.
- `ignoreDeletedStatus` (boolean, optional): If set to `true`, the endpoint will ignore subnets with a `Deleted` status.

**Headers**:
- `Authorization` (string, required): The Bearer token for authentication with the Aruba Cloud API.

</details>

<details>
<summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The request was successful and the subnet details are returned.
- `400 Bad Request`: The request is invalid. Ensure that the path parameters are correct.
- `401 Unauthorized`: The request is not authorized.
- `404 Not Found`: The specified subnet does not exist in the given project and VPC.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "category": {
    "name": "Networking",
    "provider": "Aruba.Network",
    "typology": {
      "id": "subnet",
      "name": "SUBNET"
    }
  },
  "createdBy": "<USER_ID>",
  "creationDate": "2025-10-07T15:11:17.005+00:00",
  "id": "<SUBNET_ID>",
  "location": {
    "city": "Bergamo",
    "code": "IT BG",
    "country": "IT",
    "name": "Bergamo - Nord Italia",
    "value": "ITBG-Bergamo"
  },
  "name": "test-subnet-kog-0710",
  "project": {
    "id": "<PROJECT_ID>"
  },
  "properties": {
    "dhcp": {
      "enabled": true
    },
    "network": {
      "address": "192.168.2.0/24",
      "gateway": "192.168.2.1"
    },
    "type": "Basic",
    "vpc": {
      "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>"
    }
  },
  "status": {
    "creationDate": "2025-10-07T15:12:50.694+00:00",
    "state": "Active"
  },
  "tags": [
    "tag1",
    "tag2"
  ],
  "updateDate": "2025-10-07T15:12:50.694+00:00",
  "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>/subnets/<SUBNET_ID>",
  "version": "1.0"
}
```

</details>

---

### Create Subnet endpoint

**Description**:
This endpoint creates a new subnet in the specified Aruba Cloud project and VPC with the provided details in the request body.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The endpoint exists to flatten the `metadata` field used in the request and response bodies of the subnet resource in the Aruba Cloud API. This is necessary since the field `metadata.name` is used as resource **identifer** for the subnet by `oasgen-provider` and as a consequence it is put in the status of the generated CRD. Due to a current limitation of the underlying CRD generation library, **nested fields used as identifiers are not fully supported**.
- Therefore the plugin will accept a request body where the entire `metadata` object is flattened. In particular, the request body will have a top-level `name` field (among others), and it will internally map it to the `metadata.name` field expected by the Aruba Cloud API.
- Similarly, after receiving the response from the Aruba Cloud API, the plugin will flatten the `metadata` field in the response body before returning it to the client (client == `rest-dynamic-controller` in this case).

</details>

<details><summary><b>Request</b></summary>
<br/>

```http
POST /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets
```

**Path parameters**:
- `projectId` (string, required): The ID of the Aruba Cloud project.
- `vpcId` (string, required): The ID of the VPC.

**Headers**:
- `Authorization` (string, required): The Bearer token for authentication with the Aruba Cloud API.

**Request body example**:
```json
{
  "name": "test-subnet-kog-0710",
  "properties": {
    "default": false,
    "type": "Basic"
  },
  "tags": [
    "tag1",
    "tag2"
  ]
}
```

</details>

<details><summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The pipeline was successfully updated.
- `400 Bad Request`: The request is invalid. Ensure that the path parameters are correct and the request body is well-formed.
- `401 Unauthorized`: The request is not authorized.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "category": {
    "name": "Networking",
    "provider": "Aruba.Network",
    "typology": {
      "id": "subnet",
      "name": "SUBNET"
    }
  },
  "createdBy": "<USER_ID>",
  "creationDate": "2025-10-07T15:11:17.0050334+00:00",
  "id": "<SUBNET_ID>",
  "location": {
    "city": "Bergamo",
    "code": "IT BG",
    "country": "IT",
    "name": "Bergamo - Nord Italia",
    "value": "ITBG-Bergamo"
  },
  "name": "test-subnet-kog-0710",
  "project": {
    "id": "<PROJECT_ID>"
  },
  "properties": {
    "dhcp": {
      "enabled": true
    },
    "network": {
      "address": "192.168.2.0/24",
      "gateway": "192.168.2.1"
    },
    "type": "Basic",
    "vpc": {
      "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>"
    }
  },
  "status": {
    "creationDate": "2025-10-07T15:11:16.994605+00:00",
    "state": "InCreation"
  },
  "tags": [
    "tag1",
    "tag2"
  ],
  "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>/subnets/<SUBNET_ID>",
  "version": "1.0"
}
```

</details>

---

### Update Subnet endpoint

**Description**:
This endpoint updates a specific subnet by its ID in the specified Aruba Cloud project and VPC with the provided details in the request body.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The endpoint exists to flatten the `metadata` field used in the request and response bodies of the subnet resource in the Aruba Cloud API. This is necessary since the field `metadata.name` is used as resource **identifer** for the subnet by `oasgen-provider` and as a consequence it is put in the status of the generated CRD. Due to a current limitation of the underlying CRD generation library, **nested fields used as identifiers are not fully supported**.
- Therefore the plugin will accept a request body where the entire `metadata` object is flattened. In particular, the request body will have a top-level `name` field (among others), and it will internally map it to the `metadata.name` field expected by the Aruba Cloud API.
- Similarly, after receiving the response from the Aruba Cloud API, the plugin will flatten the `metadata` field in the response body before returning it to the client (client == `rest-dynamic-controller` in this case).

</details>

<details><summary><b>Request</b></summary>
<br/>

```http
PUT /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets/{id}
```

**Path parameters**:
- `projectId` (string, required): The ID of the Aruba Cloud project.
- `vpcId` (string, required): The ID of the VPC.
- `id` (string, required): The ID of the subnet to update.

**Headers**:
- `Authorization` (string, required): The Bearer token for authentication with the Aruba Cloud API.

**Request body example**:
```json
{
  "name": "test-subnet-kog-0710",
  "properties": {
    "default": false,
    "type": "Basic"
  },
  "tags": [
    "tag1",
    "tag2",
    "tag3"
  ]
}
```

</details>

<details><summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `204 No Content`: The pipeline was successfully deleted.
- `400 Bad Request`: The request is invalid or the pipeline ID does not exist.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `404 Not Found`: The specified pipeline does not exist in the project.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "category": {
    "name": "Networking",
    "provider": "Aruba.Network",
    "typology": {
      "id": "subnet",
      "name": "SUBNET"
    }
  },
  "createdBy": "<USER_ID>",
  "creationDate": "2025-10-07T15:24:26.639+00:00",
  "id": "<SUBNET_ID>",
  "location": {
    "city": "Bergamo",
    "code": "IT BG",
    "country": "IT",
    "name": "Bergamo - Nord Italia",
    "value": "ITBG-Bergamo"
  },
  "name": "test-subnet-kog-0710",
  "project": {
    "id": "<PROJECT_ID>"
  },
  "properties": {
    "dhcp": {
      "enabled": true
    },
    "network": {
      "address": "192.168.2.0/24",
      "gateway": "192.168.2.1"
    },
    "type": "Basic",
    "vpc": {
      "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>"
    }
  },
  "status": {
    "creationDate": "2025-10-07T15:32:39.7550944+00:00",
    "state": "Updating"
  },
  "tags": [
    "tag1",
    "tag2",
    "tag3"
  ],
  "updateDate": "2025-10-07T15:32:39.7551222+00:00",
  "updatedBy": "<USER_ID>",
  "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>/subnets/<SUBNET_ID>",
  "version": "1.0"
}
```

</details>

---

### List Subnets endpoint

**Description**:
This endpoint retrieves the list of all subnets in the specified Aruba Cloud project and VPC.
It returns an array of subnet objects, each containing its ID, name, and other metadata and properties.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The endpoint exists to flatten the `metadata` field used in the request and response bodies of the subnet resource in the Aruba Cloud API. This is necessary since the field `metadata.name` is used as resource **identifer** for the subnet by `oasgen-provider` and as a consequence it is put in the status of the generated CRD. Due to a current limitation of the underlying CRD generation library, **nested fields used as identifiers are not fully supported**.
- Therefore the plugin will accept a request body where the entire `metadata` object is flattened. In particular, the request body will have a top-level `name` field (among others), and it will internally map it to the `metadata.name` field expected by the Aruba Cloud API.
- Similarly, after receiving the response from the Aruba Cloud API, the plugin will flatten the `metadata` field in the response body before returning it to the client (client == `rest-dynamic-controller` in this case).

</details>

<details>
<summary><b>Request</b></summary>
<br/>

```http
GET /projects/{projectId}/providers/Aruba.Network/vpcs/{vpcId}/subnets
```

**Path parameters**:
- `projectId` (string, required): The ID of the Aruba Cloud project.
- `vpcId` (string, required): The ID of the VPC.

**Query parameters**:
- `api-version` (string, required): The version of the Aruba Cloud API to use. For example, `1.0`.
- `filter` (string, optional): Filter expression.
- `sort` (string, optional): Sort expression.
- `projection` (string, optional): Projection expression.
- `offset` (integer, optional): Offset for pagination.
- `limit` (integer, optional): Limit for pagination.

**Headers**:
- `Authorization` (string, required): The Bearer token for authentication with the Aruba Cloud API.

</details>

<details>
<summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The request was successful and the subnet details are returned.
- `400 Bad Request`: The request is invalid. Ensure that the path parameters are correct.
- `401 Unauthorized`: The request is not authorized.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "total": 2,
  "values": [
    {
      "id": "<SUBNET_ID>",
      "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>/subnets/<SUBNET_ID>",
      "name": "automatic-subnet-01",
      "location": {
        "code": "IT BG",
        "country": "IT",
        "city": "Bergamo",
        "name": "Bergamo - Nord Italia",
        "value": "ITBG-Bergamo"
      },
      "project": {
        "id": "<PROJECT_ID>"
      },
      "category": {
        "name": "Networking",
        "provider": "Aruba.Network",
        "typology": {
          "id": "subnet",
          "name": "Subnet"
        }
      },
      "creationDate": "2025-09-30T08:02:30.14+00:00",
      "createdBy": "<USER_ID>",
      "updateDate": "2025-10-05T15:15:29.799+00:00",
      "updatedBy": "<USER_ID>",
      "version": "1.0",
      "status": {
        "state": "Active",
        "creationDate": "2025-10-05T15:15:29.799+00:00",
        "disableStatusInfo": {}
      },
      "properties": {
        "vpc": {
          "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>"
        },
        "type": "Basic",
        "default": true,
        "network": {
          "address": "192.168.1.0/24",
          "gateway": "192.168.1.1"
        },
        "dhcp": {
          "enabled": true
        }
      }
    },
    {
      "id": "<SUBNET_ID>",
      "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>/subnets/<SUBNET_ID>",
      "name": "test-subnet-kog-123",
      "location": {
        "code": "IT BG",
        "country": "IT",
        "city": "Bergamo",
        "name": "Bergamo - Nord Italia",
        "value": "ITBG-Bergamo"
      },
      "project": {
        "id": "<PROJECT_ID>"
      },
      "tags": [
        "tag1",
        "tag2"
      ],
      "category": {
        "name": "Networking",
        "provider": "Aruba.Network",
        "typology": {
          "id": "subnet",
          "name": "SUBNET"
        }
      },
      "creationDate": "2025-10-05T15:32:37.374+00:00",
      "createdBy": "<USER_ID>",
      "updateDate": "2025-10-05T15:39:28.59+00:00",
      "updatedBy": "<USER_ID>",
      "version": "1.0",
      "status": {
        "state": "Active",
        "creationDate": "2025-10-05T15:39:28.59+00:00",
        "disableStatusInfo": {}
      },
      "properties": {
        "vpc": {
          "uri": "/projects/<PROJECT_ID>/providers/Aruba.Network/vpcs/<VPC_ID>"
        },
        "type": "Basic",
        "network": {
          "address": "192.168.0.0/24",
          "gateway": "192.168.0.1"
        },
        "dhcp": {
          "enabled": true
        }
      }
    }
  ]
}
```

</details>

---

## Authentication

The plugin will forward the `Authorization` header passed in the request to this plugin to the Aruba Cloud API.
In particular, it supports the Bearer authentication scheme.

You can get more information in the main [README](../README.md#authentication).

## Documentation

Each plugin serves its own OpenAPI specification. The documentation is generated using the `swag` tool and is stored within each plugin's directory (e.g., `cmd/subnet-plugin/docs`).

To generate or update the documentation for a specific plugin, run the `swag-init.sh` script from this `plugins` directory, passing the plugin's name as an argument.

**Example: Generate docs for `subnet-plugin`**
```sh
./scripts/swag-init.sh subnet-plugin
```

This will generate the necessary `swagger.json`, `swagger.yaml`, and OpenAPI v3 files in the `cmd/subnet-plugin/docs/` directory.

You can then access the Swagger UI for each plugin at `/swagger/index.html`.

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
