# Quick Start Guide

## Introduction

The Terraform provider for SAP Cloud Identity Services enables you to automate the provisioning, management, and configuration of resources on [SAP Cloud Identity Services](https://help.sap.com/docs/cloud-identity-services?locale=en-US). By leveraging this provider, you can simplify and streamline the deployment and maintenance of applications, users, groups and schemas.

## Prerequisites

To follow along with this tutorial, ensure you have access to a [SAP Cloud Identity Services tenant](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/get-your-tenant?locale=en-US) and Terraform installed on your machine. You can download it from the official [Terraform website](https://developer.hashicorp.com/terraform/downloads).

## Authentication

In order to run the scripts, you need the credentials of an [admin](https://help.sap.com/docs/cloud-identity-services/cloud-identity-services/activate-your-account?locale=en-US) on the tenant. Terraform Provider for SAP Cloud Identity Services supports username/password based authentication only, at the moment.

#### Windows

For Windows you have two options to export the environment variables:

If you use Windows CMD, do the export via the following commands:

```Shell
set SCI_USERNAME=<your_username>
set SCI_PASSWORD=<your_password>
```

If you use Powershell, do the export via the following commands:

```Shell
$Env:SCI_USERNAME = '<your_username>'
$Env:SCI_PASSWORD = '<your_password>'
```

#### Mac

For Mac OS export the environment variables via:

```Shell
export SCI_USERNAME=<your_username>
export SCI_PASSWORD=<your_password>
```

#### Linux

For Linux export the environment variables via:

```Shell
export SCI_USERNAME=<your_username>
export SCI_PASSWORD=<your_password>
```

Replace `<your_username>` and `<your_password>` with your admin username and password.

## Documentation

Terraform Provider for SAP Cloud Identity Services [Documentation]((https://registry.terraform.io/providers/SAP/sap-cloud-identity-services/latest/docs))
