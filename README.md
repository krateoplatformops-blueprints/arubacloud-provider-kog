# Aruba Cloud Provider KOG Blueprint

***KOG***: (*Krateo Operator Generator*)

This is a Krateo Blueprint that deploys the Aruba Cloud Provider KOG leveraging the [OASGen Provider](https://github.com/krateoplatformops/oasgen-provider).
This provider allows you to manage Aruba Cloud resources such as subnets in a cloud-native way using the Krateo platform.

## Summary

- [Requirements](#requirements)
- [Project structure](#project-structure)
- [How to install](#how-to-install)
  - [Full provider installation](#full-provider-installation)
  - [Single resource installation](#single-resource-installation)
- [Supported resources](#supported-resources)
  - [Resource details](#resource-details)
    - [Subnet](#subnet)
  - [Resource examples](#resource-examples)
- [Authentication](#authentication)
- [Configuration](#configuration)
  - [Configuration resources](#configuration-resources)
  - [values.yaml](#valuesyaml)
  - [Verbose logging](#verbose-logging)
- [Chart structure](#chart-structure)
- [Troubleshooting](#troubleshooting)

## Requirements

[OASGen Provider](https://github.com/krateoplatformops/oasgen-provider) should be installed in your cluster. Follow the related Helm Chart [README](https://github.com/krateoplatformops/oasgen-provider-chart) for installation instructions.
Note that a standard installation of Krateo contains the OASGen Provider.

## Project structure

This project is composed by the following folders:
- **arubacloud-provider-kog-*-blueprint**: Helm charts that deploys single resources supported by this provider. These charts are useful if you want to deploy only one of the supported resources.
- **arubacloud-provider-kog-blueprint**: a Helm chart that can deploy all resources supported by this provider. It is useful if you want to manage multiple of the supported resources.
- **plugins**: a folder that is a monorepo containing multiple Go plugins. If needed, they are deployed as part of the Helm chart of the specific resource.

## How to install

### Full provider installation

To install the **arubacloud-provider-kog-blueprint** Helm chart (full provider), use the following command:

```sh
helm install arubacloud-provider-kog arubacloud-provider-kog \
  --repo https://marketplace.krateo.io \
  --namespace <release-namespace> \
  --create-namespace \
  --version 1.0.0 \
  --wait
```

> [!NOTE]
> Due to the nature of the providers leveraging the [OASGen Provider](https://github.com/krateoplatformops/oasgen-provider), this chart will install a set of RestDefinitions that will in turn trigger the deployment of a set controllers in the cluster. These controllers need to be up and running before you can create or manage resources using the Custom Resources (CRs) defined by this provider. This may take a few minutes after the chart is installed. The RestDefinitions will reach the condition `Ready` when the related CRDs are installed and the controllers are up and running.

You can check the status of the RestDefinitions with the following commands:

```sh
kubectl get restdefinitions.ogen.krateo.io --all-namespaces | awk 'NR==1 || /arubacloud/'
```
You should see output similar to this:
```sh
NAMESPACE       NAME                                 READY   AGE
krateo-system   arubacloud-provider-kog-subnet       False   59s
```

You can also wait for a specific RestDefinition (`arubacloud-provider-kog-subnet` in this case) to be ready with a command like this:
```sh
kubectl wait restdefinitions.ogen.krateo.io arubacloud-provider-kog-subnet --for condition=Ready=True --namespace krateo-system --timeout=300s
```

Note that the names of the RestDefinitions and the namespace where the RestDefinitions are installed may vary based on your configuration.

### Single resource installation

To manage a single resource, you can install the specific Helm chart for that resource. For example, to install the `arubacloud-provider-kog-subnet` resource, you can use the following command:

```sh
helm install arubacloud-provider-kog-subnet arubacloud-provider-kog-subnet \
  --repo https://marketplace.krateo.io \
  --namespace <release-namespace> \
  --create-namespace \
  --version 1.0.0 \
  --wait
```

## Supported resources

This chart supports the following resources and operations:

| Resource     | Get  | Create | Update | Delete |
|--------------|------|--------|--------|--------|
| Subnet       | ✅   | ✅     | ✅     | ✅     |


The resources listed above are Custom Resources (CRs) defined in the `arubacloud.ogen.krateo.io` API group. They are used to manage Aruba Cloud resources in a Kubernetes-native way, allowing you to create, update, and delete Arubacloud resources using Kubernetes manifests.

### Resource details

#### Subnet

The `Subnet` resource allows you to create, update, and delete Aruba Cloud subnets.
You can specify the subnet name, location, tags, type, and other settings such as DHCP configuration and routes.

An example of a Subnet resource is:
```yaml
apiVersion: arubacloud.ogen.krateo.io/v1alpha1
kind: Subnet
metadata:
  name: test-subnet-kog-123
  namespace: default
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  configurationRef:
    name: example-configuration
    namespace: config-namespace
  projectId: "proj-12345"
  vpcId: "vpc-67890"
  name: "example-subnet-name"
  location:
    value: "ITBG-Bergamo"
  newDefaultSubnet: "" # URI for existing subnet to set as default, if needed during deletion of this subnet
  tags:
  - "tag1"
  - "tag2"
  properties:
    default: false
    type: "Advanced" # allowed values: {Basic, Advanced}
    network:
      address: "10.0.0.0/8" # Address of the network in CIDR Notation. The IP range must be between 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
    dhcp:
      enabled: true
      dns:
      - "8.8.8.8"
      - "8.8.4.4"
      range:
        start: "10.1.0.10"
        count: 200
      routes:
      - address: "192.168.0.0/16"
        gateway: "10.1.0.11"
      - address: "172.16.0.0/12"
        gateway: "10.1.0.12"
```

### Resource examples

You can find example resources for each supported resource type in the `/samples` folder of the main chart.

## Authentication

The authentication to the Aruba Cloud API is managed using 2 kinds of resources (both are required):

- **Kubernetes Secret**: This resource is used to store the Aruba Cloud Token that is used to authenticate with the Aruba Cloud API. 

In order to generate a Aruba Cloud token, follow these instructions: https://api.arubacloud.com/docs/authentication/.

Example of a Kubernetes Secret that you can apply to your cluster:
```sh
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: arubacloud-token
  namespace: default
type: Opaque
stringData:
  token: <YOUR_TOKEN>
EOF
```

Replace `<YOUR_TOKEN>` with your actual Aruba Cloud Token (without quotes and without `Bearer ` prefix).

- **\<Resource\>Configuration**: These resource can reference the Kubernetes Secret and are used to authenticate with the Aruba Cloud API. They must be referenced with the `configurationRef` field of the resources defined in this chart. The configuration resource can be in a different namespace than the resource itself.

Note that the specific configuration resource type depends on the resource you are managing. For instance, in the case of the `Subnet` resource, you would need a `SubnetConfiguration`.

An example of a `SubnetConfiguration` resource that references the Kubernetes Secret, to be applied to your cluster:
```sh
kubectl apply -f - <<EOF
apiVersion: arubacloud.ogen.krateo.io/v1alpha1
kind: SubnetConfiguration
metadata:
  name: my-subnet-config
  namespace: default
spec:
  authentication:
    bearer:
      tokenRef:
        name: arubacloud-token
        namespace: default
        key: token
  configuration:
    query:
      create:
        api-version: "1.0"
      delete:
        api-version: "1.0"
      get:
        api-version: "1.0"
        ignoreDeletedStatus: false
      update:
        api-version: "1.0"
      findby:
        api-version: "1.0"
EOF
```

Then, in the `Subnet` resource, you can reference the `SubnetConfiguration` resource as follows:
```yaml
apiVersion: arubacloud.ogen.krateo.io/v1alpha1
kind: Subnet
metadata:
  name: test-subnet-kog-123
  namespace: default
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  configurationRef:
    name: my-subnet-config
    namespace: default 
  projectId: ABCDEFGHIJKLMN
  vpcId: ABC1234567890
  name: test-subnet-kog-123
```

More details about the configuration resources in the [Configuration resources](#configuration-resources) section below.

## Configuration

### Configuration resources

Each resource type (e.g., `Subnet`) requires a specific configuration resource (e.g., `SubnetConfiguration`) to be created in the cluster.
Currently, the supported configuration resources are:
- `SubnetConfiguration`

These configuration resources are used to store the authentication information (i.e., reference to the Kubernetes Secret containing the Aruba Cloud Token) and other configuration options for the resource type.
You can find examples of these configuration resources in the `/samples/configs` folder of the main chart.
Note that a single configuration resource can be used by multiple resources of the same type.
For example, you can create a single `SubnetConfiguration` resource and reference it in multiple `Subnet` resources.

### values.yaml

You can customize the **arubacloud-provider-kog-blueprint** chart by modifying the `values.yaml` file.
For instance, you can select which resources the provider should support in the oncoming installation.
This may be useful if you want to limit the resources managed by the provider to only those you need, reducing the overhead of managing unnecessary controllers.
The default configuration of the chart enables all resources supported by the chart.

Note: currently `subnet` is the only supported resource.

### Verbose logging

In order to enable verbose logging for the controllers, you can add the `krateo.io/connector-verbose: "true"` annotation to the metadata of the resources you want to manage, as shown in the examples above. 
This will enable verbose logging for those specific resources, which can be useful for debugging and troubleshooting as it will provide more detailed information about the operations performed by the controllers.

## Charts structure

Main components of the charts:

- **RestDefinitions**: These are the core resources needed to manage resources leveraging the OASGen Provider. In this case, they refers to the OpenAPI Specification to be used for the creation of the Custom Resources (CRs) that represent Aruba Cloud resources.
They also define the operations that can be performed on those resources. Once the chart is installed, RestDefinitions will be created and as a result, specific controllers will be deployed in the cluster to manage the resources defined with those RestDefinitions.

- **ConfigMaps**: Refer directly to the OpenAPI Specification content in the `/assets` folder.

- **/assets** folder: Contains the selected OpenAPI Specification files for the Aruba Cloud API.

- **Deployment** (optional): Deploys a plugin that is used as a proxy to resolve some integration issue with Aruba Cloud. The specific endpoins managed by the plugin are described in the [plugins README](./plugins/README.md)

- **Service** (optional): Exposes the plugin described above, allowing the resource controllers to communicate with the Aruba Cloud API through the plugin, only if needed.

## Troubleshooting

For troubleshooting, you can refer to the [Troubleshooting guide](./arubacloud-provider-kog-blueprint/docs/troubleshooting.md) in the `/docs` folder of the main blueprint (chart). 
It contains common issues and solutions related to this chart.

## Release process

Please refer to the [Release guide](./docs/release.md) in the `/docs` folder for detailed instructions on how to release new versions of the chart and its components.

