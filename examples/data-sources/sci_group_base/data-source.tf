# Read a group by ID 
data "sci_group_base" "by_id" {
  group_id = "group_1234567890" # Must be a valid UUID
}