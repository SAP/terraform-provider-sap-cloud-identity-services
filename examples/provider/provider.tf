terraform {
  required_providers {
    ias = {
      source  = "SAP/ias"
      version = "0.1.0-beta1"
    }
  }
}

# Configure the BTP Provider
provider "ias" {
  tenant_url = "https://<tenant>.authentication.eu10.hana.ondemand.com"
}
