terraform {
  required_providers {
    sci = {
      source  = "SAP/sap-cloud-identity-services"
    }
  }
}

# Configure the SAP Cloud Identity Services Provider
provider "sci" {
  tenant_url = "https://iasprovidertestblr.accounts400.ondemand.com"
  certificate_path = "terraform-eds-cert.pem"
  private_key_path = "terraform-eds-key.pem"
}

# List all groups
data "sci_groups" "all" {
}

output "test_groups" {
  value = data.sci_groups.all
  
}