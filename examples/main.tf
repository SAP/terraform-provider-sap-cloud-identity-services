# data "ias_application" "name" {
#   id = "c5c483bb-62de-444e-b784-ded7d369eabd"
# }

resource "ias_application" "terraform_test" {
  name = "terraform_test"
  description = "sample application"
}