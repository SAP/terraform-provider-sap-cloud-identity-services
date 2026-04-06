# Create a basic application in SAP Cloud Identity Services
resource "sci_application" "basic_application" {
  name                  = "My Basic Application"
  description           = "A basic application in SAP Cloud Identity Services"
  parent_application_id = "app_0987654321" # Must be a valid UUID
  multi_tenant_app      = false
  authentication_schema = {
    sso_type = "saml2" # Refer to the documentation for valid values
    subject_name_identifier = {
      source = "Identity Directory" # Refer to the documentation for valid values
      value  = "uid"
    }
    subject_name_identifier_function = "upperCase" # Refer to the documentation for valid values
    assertion_attributes = [
      {
        attribute_name  = "user_attribute"
        attribute_value = "mail"
      }
    ]
    advanced_assertion_attributes = [
      {
        source          = "Corporate Identity Provider" # Refer to the documentation for valid values
        attribute_name  = "user_attribute"
        attribute_value = "test"
      }
    ]
    default_authenticating_idp = "idp_1234567890" # Must be a valid UUID or the internal ID of the IdP
    conditional_authentication = [
      {
        identity_provider_id = "idp_0987654321" # Must be a valid UUID or the internal ID of the IdP
        user_email_domain    = "example.com"
        user_type            = "employee"         # Refer to the documentation for valid values
        user_group           = "group_1234567890" # Must be a valid UUID
        ip_network_range     = "10.0.0.8/16"
      }
    ]
    rest_api_authentication = {
      allow_public_client_flows = true
      all_apis_access          = false
      allow_locking           = true
      unlock                   = false
    }
  }
}

# Create a SAML2 application in SAP Cloud Identity Services
resource "sci_application" "saml2_application" {
  name                  = "My Basic SAML2 Application"
  description           = "A basic saml2 application in SAP Cloud Identity Services"
  authentication_schema = {
    sso_type = "saml2"
    saml2_config = {
      saml_metadata_url = "https://example.com/saml/metadata"
      acs_endpoints = [
        {
          binding = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"  # Refer to the documentation for valid values
          location = "https://example.com/saml/acs"
          index = 0
          default = true
        }
      ]
      slo_endpoints = [
        {
          binding_name = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" # Refer to the documentation for valid values
          location = "https://example.com/saml/slo"
          response_location = "https://example.com/saml/slo/response"
        }
      ]
      signing_certificates = [
        {
          base64_certificate = "-----BEGIN CERTIFICATE-----<vali-base64-encoded-certificate>-----END CERTIFICATE-----"
          default = true
        }
      ]
      encryption_certificate = {
        base64_certificate = "-----BEGIN CERTIFICATE-----<vali-base64-encoded-certificate>-----END CERTIFICATE-----"
      }
      response_elements_to_encrypt = ["attributes"]         # Refer to the documentation for valid values
      default_name_id_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"     # Refer to the documentation for valid values 
      sign_slo_messages = true
      require_signed_slo_messages = true
      sign_assertions = false
      sign_auth_responses = false
      digest_algorithm = "sha256"     # Refer to the documentation for valid values
    }
  }
}

# Create a OIDC application in SAP Cloud Identity Services
resource "sci_application" "oidc_application" {
  name                  = "My Basic OIDC Application"
  description           = "A basic oidc application in SAP Cloud Identity Services"
  authentication_schema = {
    sso_type = "openIdConnect"
    oidc_config = {
      redirect_uris = [
        "https://example.com/oidc/callback"
      ]
      post_logout_redirect_uris = [
        "https://example.com/oidc/logout/callback"
      ]
      front_channel_logout_uris = [
        "https://example.com/oidc/frontchannel-logout"
      ]
      back_channel_logout_uris = [
        "https://example.com/oidc/backchannel-logout"
      ]
      token_policy = {
        jwt_validity = 3600
        refresh_validity = 7200
        refresh_parallel = 5
        max_exchange_period = "unlimited"   # Refer to the documentation for valid values
        refresh_token_rotation_scenario = "online"  # Refer to the documentation for valid values
        access_token_format = "default"  # Refer to the documentation for valid values
      }
      restricted_grant_types = [
        "authorizationCode"                 # Refer to the documentation for valid values
      ]
      proxy_config = {
        acrs = [
          "example_value"
        ]
      }
    }
  }
}