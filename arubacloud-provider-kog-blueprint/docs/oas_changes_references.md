# Aruba Cloud Provider KOG - OpenAPI Specification (OAS) Changes reference

This documents serves as a reference for the changes made to the OpenAPI Specification (OAS) of the resources managed by the Aruba Cloud Provider KOG.
Note that the changes are made to comply with some requirements of the OASGen provider or Rest Dynamic Controller or to fix issues in the original OAS.

OAS source: https://api.arubacloud.com/openapi/network-provider.json

## Security scheme changes

The original security scheme definition in the OAS source is:
```yaml
  securitySchemes:
    Bearer:
      type: apiKey
      description: Insert JWT token
      name: Authorization
      in: header
security:
- Bearer: []
```

In order to actually use Bearer authentication, it has been changed to:
```yaml
  securitySchemes:
    accessToken:
      type: http
      scheme: bearer
security:
- accessToken: []
```
