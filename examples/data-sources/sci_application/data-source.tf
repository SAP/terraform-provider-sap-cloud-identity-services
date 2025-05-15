# Read an application by ID
data "sci_application" "by_id" {
  id = "app_1234567890"     # Must be a valid UUID
}
