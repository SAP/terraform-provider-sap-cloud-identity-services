terraform {
  required_providers {
    sci = {
      source  = "SAP/sci"
      version = "0.1.0-beta1"
    }
  }
}

# Configure the BTP Provider
provider "sci" {
  tenant_url = "https://<tenant>.authentication.eu10.hana.ondemand.com"
}
