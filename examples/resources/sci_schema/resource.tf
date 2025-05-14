# Create a basic schema in SAP Cloud Identity Services
resource "sci_schema" "basic_schema" {
  id       = "urn:sap:Terraform"               # Must follow the pattern : urn:<namespace-identifier>:<resource-type>
  name     = "Terraform"
  attributes = [
    {
      name = "test_attr"
      type = "string"                          # Refer to the documentation for valid values
      mutability = "writeOnly"                 # Refer to the documentation for valid values
      returned = "default"                     # Refer to the documentation for valid values
      uniqueness = "none"                      # Refer to the documentation for valid values
      canonical_values = [
        "val1",
        "val2"
      ]
      multivalued = true
      description = "Test Attribute"
      required = true
      case_exact = false
    } 
  ]
  description = "Test Schema"
}
