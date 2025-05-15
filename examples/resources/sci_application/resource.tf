# Create a basic application in SAP Cloud Identity Services
resource "sci_application" "basic_application" {
  id          = "app_1234567890"                        # Must be a valid and unique UUID
  name        = "My Basic Application"
  description = "A basic application in SAP Cloud Identity Services"
  parent_application_id = "app_0987654321"              # Must be a valid UUID
  multi_tenant_app = false
  authentication_schema = {
    sso_type = "saml2"                                  # Refer to the documentation for valid values
    subject_name_identifier = {
      source = "Identity Directory"                     # Refer to the documentation for valid values
      value  = "uid"                      
    }
    subject_name_identifier_function = "upperCase"      # Refer to the documentation for valid values
    assertion_attributes = [
      {
        attribute_name = "user_attribute"
        attribute_value = "mail"
      }
    ]
    advanced_assertion_attributes = [
      {
        source = "Corporate Identity Provider"          # Refer to the documentation for valid values
        attribute_name = "user_attribute"
        attribute_value = "test"
      }
    ]
    default_authenticating_idp = "idp_1234567890"       # Must be a valid UUID
    conditional_authentication = [
      {
        identity_provider_id = "idp_0987654321"          # Must be a valid UUID
        user_email_domain = "example.com"                
        user_type = "employee"                           # Refer to the documentation for valid values
        user_group = "group_1234567890"                  # Must be a valid UUID
        ip_network_range = "10.0.0.8/16"          
      }
    ]
  }
}
