# Create a user in SAP Cloud Identity Services
resource "sci_user" "new_user" {
  user_name   = "jdoe"
  emails = [
    {
      value = "john.doe@sap.com",
      type  = "work"
      primary = true
    }
  ]
  name = {
    family_name = "John"
    given_name  = "Doe"
    honorific_prefix = "Mr."
  }
  initial_password = "1234"       # Must be a valid password 
  display_name = "John Doe"
  user_type = "customer"          # Refer to the documentation for valid values
  active = false
  sap_extension_user = {
    send_mail = false
    mail_verified = true
    status = "active"             # Refer to the documentation for valid values
  }
}


# Create a user in SAP Cloud Identity Services with customSchemas
resource "sci_user" "new_user" {
  schemas = [                    # If custom schemas are to be used, ensure a valid schema ID is provided
    "urn:ietf:params:scim:schemas:core:2.0:User",
    "urn:ietf:params:scim:schemas:extension:sap:2.0:User",
    "urn:custom:SCI:1.0:User"
  ]
  user_name   = "jdoe"
  emails = [
    {
      value = "john.doe@sap.com",
      type  = "work"
    }
  ]
  custom_schemas = jsonencode({
    "urn:custom:SCI:1.0:User" : {
      custom_attr : "custom_val"
    }
  })
}
