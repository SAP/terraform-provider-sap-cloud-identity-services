# Assign a member to a group
resource "sci_group_assignment" "basic_group_assignment" {
  group_id = "valid-uuid"
  group_members = [
    {
      value = "valid-uuid",
      type  = "User" # Refer to the documentation for valid values
    }
  ]
}