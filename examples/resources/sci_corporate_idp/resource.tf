# Example 1: SAML2 Corporate IdP
resource "sci_corporate_idp" "saml2_example" {
  # Required
  display_name = "My SAML2 Corporate IdP"

  # Optional top-level
  name                     = "my-saml2-idp"
  type                     = "saml2"
  logout_url               = "https://idp.example.com/logout"
  forward_all_sso_requests = false

  identity_federation = {
    use_local_user_store            = true
    allow_local_users_only          = true
    apply_local_idp_auth_and_checks = false
    required_groups                 = ["group-a", "group-b"]
  }

  login_hint_config = {
    login_hint_type = "mail"
    send_method     = "urlParam"
  }

  saml2_config = {
    saml_metadata_url = "https://idp.example.com/metadata"
    digest_algorithm  = "sha256"
    include_scoping   = true
    name_id_format    = "email"
    allow_create      = "default"

    assertion_attributes = [
      {
        name  = "firstName"
        value = "first_name"
      },
      {
        name  = "lastName"
        value = "last_name"
      }
    ]

    signing_certificates = [
      {
        base64_certificate = file("${path.module}/cert.pem")
        default            = true
        dn                 = "CN=my-saml2-cert"
        valid_from         = "2024-01-01T00:00:00Z"
        valid_to           = "2026-01-01T00:00:00Z"
      }
    ]

    sso_endpoints = [
      {
        binding_name = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
        location     = "https://idp.example.com/sso/post"
      },
      {
        binding_name = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
        location     = "https://idp.example.com/sso/redirect"
      }
    ]

    slo_endpoints = [
      {
        binding_name      = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
        location          = "https://idp.example.com/slo/post"
        response_location = "https://idp.example.com/slo/response"
      }
    ]
  }
}

# Example 2: OIDC Corporate IdP
resource "sci_corporate_idp" "oidc_example" {
  # Required
  display_name = "My OIDC Corporate IdP"

  # Optional top-level
  name                     = "my-oidc-idp"
  type                     = "openIdConnect"
  logout_url               = "https://idp.example.com/oidc/logout"
  forward_all_sso_requests = false

  identity_federation = {
    use_local_user_store            = false
    allow_local_users_only          = false
    apply_local_idp_auth_and_checks = false
    required_groups                 = []
  }

  login_hint_config = {
    login_hint_type = "mail"
    send_method     = "urlParam"
  }

  oidc_config = {
    discovery_url              = "https://idp.example.com/.well-known/openid-configuration"
    client_id                  = "my-client-id"
    client_secret              = "my-client-secret"
    subject_name_identifier    = "email"
    token_endpoint_auth_method = "clientSecretPost"
    scopes                     = ["openid", "email", "profile"]
    enable_pkce                = true

    additional_config = {
      enforce_nonce                = true
      enforce_issuer_check         = true
      disable_logout_id_token_hint = false
    }
  }
}
