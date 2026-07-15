terraform {
  required_providers {
    sci = {
      source  = "sap/sap-cloud-identity-services"
      version = "0.7.0-beta1"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.14.0"
    }
  }
}

provider "sci" {
  tenant_url               = var.tenant_url
  p12_certificate_content  = filebase64(var.certificate_file_path)
  p12_certificate_password = var.certificate_file_password
}