---
page_title: "SAP Cloud Identity Services Provider"
subcategory: ""
description: |-
  The Terraform provider for SAP Cloud Identity Services enables you to automate the provisioning, management, and configuration of resources in the SAP Cloud Identity Services https://help.sap.com/docs/cloud-identity-services.
---
# Terraform Provider for SAP Cloud Identity Services

The Terraform provider for SAP Cloud Identity Services enables you to automate the provisioning, management, and configuration of resources in the [SAP Cloud Identity Services](https://help.sap.com/docs/cloud-identity-services).

## Example Usage

```terraform
terraform {
  required_providers {
    sci = {
      source  = "SAP/sap-cloud-identity-services"
      version = "0.2.0-beta1"
    }
  }
}

# Configure the SAP Cloud Identity Services Provider
provider "sci" {
  tenant_url = "https://<tenant>.accounts.ondemand.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `tenant_url` (String) The URL of the SCI tenant.

### Optional

- `client_id` (String, Sensitive) The client ID for OAuth2 authentication.
- `client_secret` (String, Sensitive) The client secret for OAuth2 authentication.
- `p12_certificate_content` (String, Sensitive) Base64-encoded content of the `.p12` (PKCS#12) certificate bundle file used for x509 authentication. For example you can use `filebase64("certifiacte.p12")` to load the file content, But any source that provides a valid .p12 certificate base64 string is accepted.
- `p12_certificate_password` (String, Sensitive) Password to decrypt the `.p12` certificate content.
- `password` (String, Sensitive) Your password for Basic Authentication.
- `username` (String) Your user name for Basic Authentication.

## Best Practices

For the best experience using the SAP Cloud Identity Services provider, we recommend applying the common best practices for Terraform adoption as described in the Hashicorp documentation. For example, see [Phases of Terraform Adoption](https://developer.hashicorp.com/well-architected-framework/operational-excellence/operational-excellence-terraform-maturity).

## Authentication

In order to get authenticated, the credentials of an [administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/activate-your-account?locale=en-US) are required. The SAP Cloud Identity Services Provider supports the following authentication flows:

1. [Basic Authentication](#basic-auth) 
2. [X.509 Certificate Authentication](#cert-auth)
3. [OAuth2 Client Authentication](#secret-auth)

<br>

### <u><a id="basic-auth" >Basic Authentication</a></u>

You would require a valid **username** and **password** of a [User Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-user-as-administrator) to get authenticated.
 
You can configure your credentials as part of the provider configuration as shown below:

```hcl
    provider "sci" {
        tenant_url = <your_tenant_url>
        username = <your_username>
        password = <your_password>
    }
```
It is recommended to securely set your credentials as environment variables ```SCI_USERNAME``` and ```SCI_PASSWORD```. In case you want to provide the username and password via variables make sure to follow the guidance given in the [Hashicorp documentation](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables) 
and never commit the values to a source code management system.

<br>

### <u><a id="cert-auth"> X.509 Certificate Authentication </a></u>

You would require a valid **p12 certificate** and the corresponding **password** of a [System Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-system-as-administrator) to get authenticated.
 
```NOTE: Refer to step 6 in the documentation linked above, section Certificates to fetch the required credentials. ```

You can configure the credentials as part of the provider configuration as shown below:

 ```hcl
provider "sci" {
    tenant_url = <your_tenant_url>
    p12_certificate_content = <your_p12_certificate>
    p12_certificate_password = <your_p12_certificate_password>
}
```

Ensure to paste the ***content*** of your p12 certificate rather than the ***file path***.
You can even use the function `filebase64("path_to_certificate.p12")` to load the file content. 

<br>

### <u><a id = "secret-auth">OAuth2 Client Authentication</a></u>

You would require a valid **Client ID** and the corresponding **Client Secret** of a [System Administrator](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/add-administrators?version=Cloud#add-system-as-administrator) to get authenticated.

```NOTE: Refer to step 6 in the documentation linked above, section Secrets to fetch the required credentials. ```
 
You can configure them as part of the provider configuration as shown below:

```hcl
    provider "sci" {
        tenant_url = <your_tenant_url>
        client_id = <your_client_id>
        client_secret = <your_client_secret>
    }
```

It is recommended to securely set your credentials as environment variables ```SCI_UCLIENT_ID``` and ```SCI_CLIENT_SECRET```. In case you want to provide the Client ID and Secret via variables make sure to follow the guidance given in the [Hashicorp documentation](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables) 
and never commit the values to a source code management system.
