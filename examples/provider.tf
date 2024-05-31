terraform {
  required_providers {
    cloudidentityservices = {
        source = "sap/cloudidentityservices"
    }
  }
}

provider "cloudidentityservices" {
  tenant_url = "https://iasprovidertestblr.accounts400.ondemand.com/"
}