# Create a basic group in SAP Cloud Identity Services without any members
resource "sci_group_base" "basic_group" {
  display_name = "My Basic Group"
  group_extension = {
    name        = "Terraform"
    description = "Group for terraform users"
  }
}