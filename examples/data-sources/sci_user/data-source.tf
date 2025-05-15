# Read a user by ID
data "sci_user" "by_id" {
  id = "user_1234567890"      # Must be a valid UUID
}
