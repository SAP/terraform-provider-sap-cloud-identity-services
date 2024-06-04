terraform {
  required_providers {
    ias = {
        source = "sap/ias"
    }
  }
}

provider "ias" {
  tenant_url = "https://iasprovidertestblr.accounts400.ondemand.com/"
}