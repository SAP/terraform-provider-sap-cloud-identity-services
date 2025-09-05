locals {
    prefix_unique_name = "integration-test"
    prefix_display_name = "Integration Test"

    user_display_name = "${local.prefix_display_name} User"
    group_display_name = "${local.prefix_display_name} Group"
    saml_idp_display_name = "${local.prefix_display_name} SAML IdP"
    oidc_idp_display_name = "${local.prefix_display_name} OIDC IdP"

    user_unique_name = "${local.prefix_unique_name}-user"
    group_unique_name = "${local.prefix_unique_name}-group"
    saml_idp_unique_name = "${local.prefix_unique_name}-saml-idp"
    oidc_idp_unique_name = "${local.prefix_unique_name}-oidc-idp"
    parent_app_unique_name = "${local.prefix_unique_name}-parent-app"
    basic_app_unique_name = "${local.prefix_unique_name}-basic-app"
    saml_app_unique_name = "${local.prefix_unique_name}-saml-app"
    oidc_app_unique_name = "${local.prefix_unique_name}-oidc-app"

    saml_metadata_url = "https://test.com/metadata"
    saml_type = "saml2"
    saml_idp_binding_name = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
    saml_app_binding_name = "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
    saml_sso_endpoint = "https://test.com/sso"
    saml_slo_endpoint = "https://test.com/slo"
    saml_slo_response_endpoint = "https://test.com/slo/response"
    saml_algo = "sha256"

    oidc_type = "openIdConnect"

    idp_logout_url = "https://test.com/logout"
    idp_login_type = "userInput"
    idp_send_method = "authRequest"

    cert_pem = replace(var.certificate, "\\n", "\n")

}

resource "sci_schema" "testSchema" {
  id = "urn:ietf:sap:User"
  name = "User"
  attributes = [
    {
        name = "projectName"
        type = "string"
        mutability = "readWrite"
        returned = "always"
        uniqueness = "none"
        multivalued = false
        required = false
        case_exact = true
        description = "Project name assigned to the employee"
        canonical_values = [
          "SCI Terraform",
          "BTP Terraform"
        ]
    }
  ]
}

resource "sci_user" "testUser" {
    count = 2
    schemas = [
        "urn:ietf:params:scim:schemas:core:2.0:User",
        "urn:ietf:params:scim:schemas:extension:sap:2.0:User",
        sci_schema.testSchema.id
    ]
    user_name = "${local.user_unique_name}-${count.index+1}"
    emails = [
        {
            value = "${local.user_unique_name}-${count.index+1}@sap.com"
            type = "work"
            primary = true
        }
    ]
    name = {
        given_name = local.user_display_name
        family_name = " ${count.index+1}"
        honorific_prefix = "Mr."
    }
    initial_password = var.password
    display_name = "${local.user_display_name} ${count.index+1}"
    user_type = "employee"
    active = true
    sap_extension_user = {
        send_mail = false
        mail_verified = true
        status = "active"
    }
    custom_schemas = jsonencode({
        (sci_schema.testSchema.id) = {
            (sci_schema.testSchema.attributes[0].name) = (sci_schema.testSchema.attributes[0].canonical_values[count.index])
        }
    })
}

resource "sci_group" "testGroup" {

    display_name = local.group_display_name
    group_members = [
      for user in sci_user.testUser :
       {
            value = user.id
            type = "User"
       }    
    ]
    group_extension = {
        name = local.group_unique_name
        description = "Application created for integration tests"
    }
}

resource "sci_corporate_idp" "testSamlIdP" {
    display_name = local.saml_idp_display_name
    name = local.saml_idp_unique_name
    type = local.saml_type
    logout_url = local.idp_logout_url
    forward_all_sso_requests = true
    identity_federation = {
        use_local_user_store = true
        allow_local_users_only = true
        apply_local_idp_auth_and_checks = true
        required_groups = [
            sci_group.testGroup.group_extension.name
        ]
    }
    login_hint_config = {
        login_hint_type = local.idp_login_type
        send_method = local.idp_send_method
    }
    saml2_config = {
        saml_metadata_url = local.saml_metadata_url
        assertion_attributes = [
            {
                name = "uid"
                value = "uid"
            }
        ]
        signing_certificates = [
            {
                base64_certificate = local.cert_pem
                dn = local.prefix_unique_name
                default = true
                valid_from = "2025-08-25T10:30:00Z"
                valid_to = "2026-08-25T10:30:00Z"
            }
        ]
        sso_endpoints = [
            {
                binding_name = local.saml_idp_binding_name
                location = local.saml_sso_endpoint
                default = true
            }
        ]
        slo_endpoints = [
            {
                binding_name = local.saml_idp_binding_name
                location = local.saml_slo_endpoint
                response_location = local.saml_slo_response_endpoint
                default = true
            }
        ]
        digest_algorithm = local.saml_algo
        include_scoping = true
        name_id_format = "unspecified"
        allow_create = "true"
    }
}

resource "sci_corporate_idp" "testOidcIdP" {
    display_name = local.oidc_idp_display_name
    name = local.oidc_idp_unique_name
    type = local.oidc_type
    logout_url = local.idp_logout_url
    forward_all_sso_requests = true
    identity_federation = {
        use_local_user_store = true
        allow_local_users_only = true
        apply_local_idp_auth_and_checks = true
        required_groups = [
            sci_group.testGroup.group_extension.name
        ]
    }
    login_hint_config = {
        login_hint_type = local.idp_login_type
        send_method = local.idp_send_method
    }
    oidc_config = {
        discovery_url = "https://accounts.sap.com"
        client_id = "test-client-id"
        client_secret = "test-client-secret"
        token_endpoint_auth_method = "clientSecretPost"
        subject_name_identifier = "email"
        scopes = [
            "openid",
            "email"
        ]
        enable_pkce = true
		additional_config = {
            enforce_nonce = true
            enforce_issuer_check = true
            disable_logout_id_token_hint = true
        }
    }
}

resource "sci_application" "parentApp" {
    name = local.parent_app_unique_name
    description = "Integration test parent application"
}

resource "sci_application" "testApp" {
    name = local.basic_app_unique_name
    description = "Application created for integration tests"
    parent_application_id = sci_application.parentApp.id
    authentication_schema = {
        subject_name_identifier = {
            source = "Identity Directory"
            value = "uid"
        }
        subject_name_identifier_function = "upperCase"
        assertion_attributes = [
            {
                attribute_name = "Name"
                attribute_value = "firstName"
            }
        ]
        advanced_assertion_attributes = [
            {
                source = "Corporate Identity Provider"
                attribute_name = "Name"
                attribute_value = "user_name"
            },
            {
                source = "Expression"
                attribute_name = "Email"
                attribute_value = "user_email"
            }
        ]
        conditional_authentication = [
            {
                identity_provider_id = sci_corporate_idp.testOidcIdP.id
                user_type = "employee"
                user_group = sci_group.testGroup.id
                user_email_domain = "sap.com"
                ip_network_range = "10.16.12.12/24"
            }
        ]
        
    }
}

resource "sci_application" "testSamlApp" {
    name = local.saml_app_unique_name
    description = "SAML Application created for integration tests"
    authentication_schema = {
        sso_type = local.saml_type
        default_authenticating_idp = sci_corporate_idp.testSamlIdP.id
        saml2_config = {
            saml_metadata_url = local.saml_metadata_url
            acs_endpoints = [
                {
                    binding_name = local.saml_app_binding_name
                    location = local.saml_sso_endpoint
                    index = 1
                    default = true
                }
            ]
            slo_endpoints = [
                {
                    binding_name = local.saml_app_binding_name
                    location = local.saml_slo_endpoint
                    response_location = local.saml_slo_response_endpoint
                }
            ]
            signing_certificates = [
                {
                    base64_certificate = local.cert_pem
                    default = true
                }
            ]
            encryption_certificate = {
                    base64_certificate = local.cert_pem
            }
            response_elements_to_encrypt = "attributes"
            default_name_id_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"
            sign_slo_messages = true
            require_signed_slo_messages = true
            require_signed_auth_requests = true
            sign_assertions = true
            sign_auth_responses = true
            digest_algorithm = local.saml_algo
        }
    }
}

resource "sci_application" "testOidcApp" {
    name = local.oidc_app_unique_name
    description = "OIDC Application created for integration tests"
    authentication_schema = { 
        sso_type = local.oidc_type
        oidc_config = {
            redirect_uris = [
                "https://test.com/redirect"
            ]
            post_logout_redirect_uris = [
                "https://test.com/post-logout-redirect"
            ]
            front_channel_logout_uris = [
                "https://test.com/front-channel-logout"
            ]
            back_channel_logout_uris = [
                "https://test.com/back-channel-logout"
            ]
            token_policy = {
                jwt_validity = 3600
                refresh_validity = 5000
                refresh_parallel = 5
                max_exchange_period = "maxSessionValidity"
                refresh_token_rotation_scenario = "online"
                access_token_format = "jwt"
            }
            restricted_grant_types = [
                "authorizationCode"
            ]
            proxy_config = {
                acrs = [
                    "test"
                ]
            }
        }
    }
} 