resource "sci_application_secret" "example" {
  application_id       = sci_application.example.id
  description          = "My API secret"
  valid_to             = "2029-10-12T10:00:00Z"
  authorization_scopes = ["manageApp", "oAuth"]
}
