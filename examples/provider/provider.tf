terraform {
  required_providers {
    sci = {
      source  = "SAP/sap-cloud-identity-services"
      version = "0.4.0-beta1"
    }
  }
}

# Configure the SAP Cloud Identity Services Provider
provider "sci" {
  tenant_url = "https://<tenant>.accounts.ondemand.com"
}
