# Read a grou's' assignments by group ID
data "sci_group_assignments" "by_id" {
  group_id = "group_1234567890" # Must be a valid UUID
}