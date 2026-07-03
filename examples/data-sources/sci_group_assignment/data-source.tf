# Read a group assignment by ID and user ID
data "sci_group_assignment" "by_id" {
  group_id = "group_1234567890" # Must be a valid UUID
  group_member = {
    value = "user_1234567890" # Must be a valid UUID
  }
}