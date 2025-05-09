# Create a user in SAP Cloud Identity Services
resource "sci_user" "new_user" {
  user_name   = "jdoe"
  name = {
    family_name = "John"
    given_name  = "Doe"
  }
  emails = [
    {
      value = "john.doe@sap.com",
      type  = "work"
    }
  ]
}


# Create a user in SAP Cloud Identity Services with customSchemas
resource "sci_user" "new_user" {
  user_name   = "TO BE DONE"
  name = {
    family_name = "John"
    given_name  = "Doe"
  }
  emails = [
    {
      value = "john.doe@sap.com",
      type  = "work"
    }
  ]
  custom_schemas = jsonencode({
    "schema_id" : ["attr1"]}
    )
}
