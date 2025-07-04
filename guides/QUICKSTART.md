# Quick Start Guide

## Introduction

The Terraform provider for SAP Cloud Identity Services enables you to automate the provisioning, management, and configuration of resources on [SAP Cloud Identity Services](https://help.sap.com/docs/cloud-identity-services?locale=en-US). By leveraging this provider, you can simplify and streamline the deployment and maintenance of applications, users, groups and schemas.

## Prerequisites

To follow along with this tutorial, ensure you have access to a [SAP Cloud Identity Services tenant](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/get-your-tenant?locale=en-US) and Terraform installed on your machine. You can download it from the official [Terraform website](https://developer.hashicorp.com/terraform/downloads).

## Authentication

In order to run the scripts, you need the credentials of an [admin](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/activate-your-account?locale=en-US) on the tenant. Terraform Provider for SAP Cloud Identity Services supports the following authentication methods:

1. [Basic Authentication](./basic_auth.md) 
2. [X.509 Certificate Authentication](cert_auth.md)
3. [OAuth2 Client Authentication](./secret_auth.md)

Refer to the link corresponding to the chosen authentication method.


## Documentation

Terraform Provider for SAP Cloud Identity Services [Documentation](https://registry.terraform.io/providers/SAP/sap-cloud-identity-services/latest/docs)
