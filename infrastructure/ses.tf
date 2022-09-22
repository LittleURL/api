locals {
  email_from = "noreply@${local.domain}"
}

resource "aws_ses_domain_identity" "default" {
  domain = local.domain
}

resource "aws_ses_email_identity" "noreply" {
  email = local.email_from
}

# ----------------------------------------------------------------------------------------------------------------------
# Verification
# ----------------------------------------------------------------------------------------------------------------------
resource "cloudflare_record" "ses_verification" {
  zone_id = local.zone_id
  name    = "_amazonses"
  type    = "TXT"
  value   = aws_ses_domain_identity.default.verification_token
  ttl     = 600
}

resource "aws_ses_domain_identity_verification" "default" {
  domain     = aws_ses_domain_identity.default.id
  depends_on = [cloudflare_record.ses_verification]
}

# ----------------------------------------------------------------------------------------------------------------------
# DKIM
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_ses_domain_dkim" "default" {
  domain = aws_ses_domain_identity.default.domain
}

resource "cloudflare_record" "ses_dkim" {
  count = 3 // TODO: fix using for_each (TF race condition)

  zone_id = local.zone_id
  name    = "${element(aws_ses_domain_dkim.default.dkim_tokens, count.index)}._domainkey"
  type    = "CNAME"
  value   = "${element(aws_ses_domain_dkim.default.dkim_tokens, count.index)}.dkim.amazonses.com"
  ttl     = 600
}
