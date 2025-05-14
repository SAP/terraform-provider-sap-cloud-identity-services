# Create a basic group in SAP Cloud Identity Services
resource "sci_group" "basic_group" {
    display_name = "My Basic Group"
    group_members = [
      {
        value = "valid-uuid",
        type = "User"                       # Refer to the documentation for valid values
      }
    ]
    group_extension = {
      name = "Terraform"
      description = "Group for terraform users"
    }
}
