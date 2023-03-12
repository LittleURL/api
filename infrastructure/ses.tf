locals {
  email_dmarc         = "dmarcreports@${aws_route53_zone.main.name}"
  email_from          = "noreply@${aws_route53_zone.main.name}"
  email_from_friendly = "${var.application} <${local.email_from}>"
  spf_includes = join(" ", [
    for s in concat(["amazonses.com"], var.email_spf_includes) : "include:${s}"
  ])
  email_domain      = aws_route53_zone.main.name
  dkim_record_count = 3
}

resource "aws_sesv2_email_identity" "domain" {
  email_identity = local.email_domain

  dkim_signing_attributes {
    next_signing_key_length = "RSA_2048_BIT"
  }
}

resource "aws_sesv2_email_identity" "noreply" {
  email_identity = local.email_from

  dkim_signing_attributes {
    next_signing_key_length = "RSA_2048_BIT"
  }

  depends_on = [
    aws_sesv2_email_identity.domain,
    aws_route53_record.ses_dkim_domain[0]
  ]
}

resource "aws_ses_domain_mail_from" "default" {
  domain           = local.email_domain
  mail_from_domain = "bounce.${local.email_domain}"

  depends_on = [aws_sesv2_email_identity.domain]
}

resource "aws_route53_record" "ses_domain_mail_from_mx" {
  zone_id = aws_route53_zone.main.id
  name    = aws_ses_domain_mail_from.default.mail_from_domain
  type    = "MX"
  ttl     = "600"
  records = ["10 feedback-smtp.${var.aws_region}.amazonses.com"]
}

# ----------------------------------------------------------------------------------------------------------------------
# Verification
# ----------------------------------------------------------------------------------------------------------------------
# resource "aws_route53_record" "ses_verification" {
#   zone_id = aws_route53_zone.main.zone_id
#   type    = "TXT"
#   name    = "_amazonses.${aws_route53_zone.main.name}"
#   ttl     = 600
#   records = [aws_ses_domain_identity.default.verification_token]
# }

# resource "aws_ses_domain_identity_verification" "default" {
#   domain     = aws_ses_domain_identity.default.id
#   depends_on = [aws_route53_record.ses_verification]
# }

# ----------------------------------------------------------------------------------------------------------------------
# DKIM
# ----------------------------------------------------------------------------------------------------------------------

resource "aws_route53_record" "ses_dkim_domain" {
  count = local.dkim_record_count

  zone_id = aws_route53_zone.main.zone_id
  type    = "CNAME"
  ttl     = 600
  name = join(".", [
    element(aws_sesv2_email_identity.domain.dkim_signing_attributes[0].tokens, count.index),
    "_domainkey.${local.email_domain}"
  ])
  records = [join(".", [
    element(aws_sesv2_email_identity.domain.dkim_signing_attributes[0].tokens, count.index),
    "dkim.amazonses.com"
  ])]

  depends_on = [aws_sesv2_email_identity.domain]
}

# ----------------------------------------------------------------------------------------------------------------------
# SPF / DMARC
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_route53_record" "ses_spf" {
  zone_id = aws_route53_zone.main.zone_id
  type    = "TXT"
  name    = local.email_domain
  ttl     = 10800
  records = ["v=spf1 ${local.spf_includes} ~all"]
}

resource "aws_route53_record" "ses_dmarc" {
  zone_id = aws_route53_zone.main.zone_id
  type    = "TXT"
  name    = "_dmarc.${local.email_domain}"
  ttl     = 10800
  records = ["v=DMARC1;p=quarantine;pct=25;rua=mailto:${local.email_dmarc}"]
}
