# Create a user in SAP Cloud Identity Services
resource "ias_user" "new_user" {
  user_name   = "TO BE DONE"
  name        = "TO BE DONE"
}


# Create a user in SAP Cloud Identity Services with customSchemas
resource "ias_user" "new_user" {
  user_name   = "TO BE DONE"
  name        = "TO BE DONE"
  custom_schemas = jsonencoded({
    "schema_id" : {
      "attr1" : value
    }
  })
}
