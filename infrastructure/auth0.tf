resource "auth0_tenant" "default" {
  friendly_name = local.environment == "prod" ? "LittleURL" : "LittleUrl (${local.environment})"
  support_email = "support@${local.domain}"
}

# ----------------------------------------------------------------------------------------------------------------------
# API
# ----------------------------------------------------------------------------------------------------------------------
resource "auth0_resource_server" "api" {
  name        = "Public API"
  identifier  = local.domain_api
  signing_alg = "RS256"

  allow_offline_access                            = true
  token_lifetime                                  = 8600
  skip_consent_for_verifiable_first_party_clients = true
}

# ----------------------------------------------------------------------------------------------------------------------
# Client: Dashboard
# ----------------------------------------------------------------------------------------------------------------------
resource "auth0_client" "website" {
  name                       = "LittleURL"
  description                = "LittleURL Dashboard"
  app_type                   = "spa"
  is_first_party             = true
  allowed_origins            = ["https://${local.domain}"]
  token_endpoint_auth_method = "none"

  # refresh_token {
  #   rotation_type                = "rotating"
  #   expiration_type              = "expiring"
  #   infinite_idle_token_lifetime = false
  #   infinite_token_lifetime      = false
  # }
}

# ----------------------------------------------------------------------------------------------------------------------
# Client: Terraform
# ----------------------------------------------------------------------------------------------------------------------
resource "auth0_client" "api" {
  name                       = "API"
  app_type                   = "non_interactive"
  is_first_party             = true
  oidc_conformant            = true
  allowed_origins            = []
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "client_secret_post"
}

resource "auth0_client_grant" "api" {
  client_id = auth0_client.api.id
  audience  = "https://${var.auth0_domain}/api/v2/"
  scope = [
    "read:users",
    "update:users",
    "read:users_app_metadata",
    "update:users_app_metadata",
    "delete:users_app_metadata",
    "create:users_app_metadata",
  ]
}
