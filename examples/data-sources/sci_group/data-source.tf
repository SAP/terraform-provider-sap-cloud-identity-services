# Read a group by ID
data "sci_group" "by_id" {
  id = "group_1234567890" # Must be a valid UUID
}
