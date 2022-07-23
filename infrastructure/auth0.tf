resource "auth0_tenant" "default" {
  friendly_name = local.environment == "prod" ? "LittleURL" : "LittleUrl (${local.environment})"
  support_email = "support@${local.domain}"

  idle_session_lifetime = 72
  sandbox_version       = "12"
  session_lifetime      = 168

  session_cookie {
    mode = "non-persistent"
  }

  change_password {
    enabled = false
    html    = ""
  }

  error_page {
    html          = ""
    url           = ""
    show_log_link = false
  }

  guardian_mfa_page {
    enabled = true
    html    = ""
  }

  flags {
    allow_legacy_delegation_grant_types    = false
    allow_legacy_ro_grant_types            = false
    allow_legacy_tokeninfo_endpoint        = false
    dashboard_insights_view                = false
    dashboard_log_streams_next             = false
    disable_clickjack_protection_headers   = false
    disable_fields_map_fix                 = false
    disable_management_api_sms_obfuscation = false
    enable_adfs_waad_email_verification    = false
    enable_apis_section                    = false
    enable_client_connections              = false
    enable_custom_domain_in_emails         = false
    enable_dynamic_client_registration     = false
    enable_idtoken_api2                    = false
    enable_legacy_logs_search_v2           = false
    enable_legacy_profile                  = false
    enable_pipeline2                       = false
    enable_public_signup_user_exists_error = false
    no_disclose_enterprise_connections     = false
    revoke_refresh_token_grant             = false
    universal_login                        = true
    use_scope_descriptions_for_consent     = false
  }

  lifecycle {
    prevent_destroy = true
  }
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
# Client: API
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

# ----------------------------------------------------------------------------------------------------------------------
# Client: Terraform
# ----------------------------------------------------------------------------------------------------------------------
resource "auth0_client" "terraform" {
  name                       = "Terraform"
  logo_uri                   = "https://www.terraform.io/favicon.ico"
  app_type                   = "non_interactive"
  is_first_party             = true
  oidc_conformant            = true
  allowed_origins            = []
  grant_types                = ["client_credentials"]
  token_endpoint_auth_method = "client_secret_post"
}

# TODO: reduce scopes to minimum needed
resource "auth0_client_grant" "terraform" {
  client_id = auth0_client.terraform.id
  audience  = "https://${var.auth0_domain}/api/v2/"
  scope = [
    "read:client_grants",
    "create:client_grants",
    "delete:client_grants",
    "update:client_grants",
    "read:users",
    "update:users",
    "delete:users",
    "create:users",
    "read:users_app_metadata",
    "update:users_app_metadata",
    "delete:users_app_metadata",
    "create:users_app_metadata",
    "read:user_custom_blocks",
    "create:user_custom_blocks",
    "delete:user_custom_blocks",
    "create:user_tickets",
    "read:clients",
    "update:clients",
    "delete:clients",
    "create:clients",
    "read:client_keys",
    "update:client_keys",
    "delete:client_keys",
    "create:client_keys",
    "read:connections",
    "update:connections",
    "delete:connections",
    "create:connections",
    "read:resource_servers",
    "update:resource_servers",
    "delete:resource_servers",
    "create:resource_servers",
    "read:device_credentials",
    "update:device_credentials",
    "delete:device_credentials",
    "create:device_credentials",
    "read:rules",
    "update:rules",
    "delete:rules",
    "create:rules",
    "read:rules_configs",
    "update:rules_configs",
    "delete:rules_configs",
    "read:hooks",
    "update:hooks",
    "delete:hooks",
    "create:hooks",
    "read:actions",
    "update:actions",
    "delete:actions",
    "create:actions",
    "read:email_provider",
    "update:email_provider",
    "delete:email_provider",
    "create:email_provider",
    "blacklist:tokens",
    "read:stats",
    "read:insights",
    "read:tenant_settings",
    "update:tenant_settings",
    "read:logs",
    "read:logs_users",
    "read:shields",
    "create:shields",
    "update:shields",
    "delete:shields",
    "read:anomaly_blocks",
    "delete:anomaly_blocks",
    "update:triggers",
    "read:triggers",
    "read:grants",
    "delete:grants",
    "read:guardian_factors",
    "update:guardian_factors",
    "read:guardian_enrollments",
    "delete:guardian_enrollments",
    "create:guardian_enrollment_tickets",
    "read:user_idp_tokens",
    "create:passwords_checking_job",
    "delete:passwords_checking_job",
    "read:custom_domains",
    "delete:custom_domains",
    "create:custom_domains",
    "update:custom_domains",
    "read:email_templates",
    "create:email_templates",
    "update:email_templates",
    "read:mfa_policies",
    "update:mfa_policies",
    "read:roles",
    "create:roles",
    "delete:roles",
    "update:roles",
    "read:prompts",
    "update:prompts",
    "read:branding",
    "update:branding",
    "delete:branding",
    "read:log_streams",
    "create:log_streams",
    "delete:log_streams",
    "update:log_streams",
    "create:signing_keys",
    "read:signing_keys",
    "update:signing_keys",
    "read:limits",
    "update:limits",
    "create:role_members",
    "read:role_members",
    "delete:role_members",
    "read:entitlements",
    "read:attack_protection",
    "update:attack_protection",
    "read:organizations",
    "update:organizations",
    "create:organizations",
    "delete:organizations",
    "create:organization_members",
    "read:organization_members",
    "delete:organization_members",
    "create:organization_connections",
    "read:organization_connections",
    "update:organization_connections",
    "delete:organization_connections",
    "create:organization_member_roles",
    "read:organization_member_roles",
    "delete:organization_member_roles",
    "create:organization_invitations",
    "read:organization_invitations",
    "delete:organization_invitations",
  ]
}
