locals {
  email_from          = "noreply@${aws_route53_zone.main.name}"
  email_from_friendly = "${var.application} <${local.email_from}>"
  spf_includes = join(" ", [
    for s in concat(["amazonses.com"], var.email_spf_includes) : "include:${s}"
  ])
}

resource "aws_ses_domain_identity" "default" {
  domain = aws_route53_zone.main.name
}

resource "aws_ses_email_identity" "noreply" {
  email = local.email_from
}

resource "aws_route53_record" "ses_spf" {
  zone_id = aws_route53_zone.main.zone_id
  type    = "TXT"
  name    = aws_route53_zone.main.name
  ttl     = 10800
  records = ["v=spf1 ${local.spf_includes} ~all"]
}

# ----------------------------------------------------------------------------------------------------------------------
# Verification
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_route53_record" "ses_verification" {
  zone_id = aws_route53_zone.main.zone_id
  type    = "TXT"
  name    = "_amazonses.${aws_route53_zone.main.name}"
  ttl     = 600
  records = [aws_ses_domain_identity.default.verification_token]
}

resource "aws_ses_domain_identity_verification" "default" {
  domain     = aws_ses_domain_identity.default.id
  depends_on = [aws_route53_record.ses_verification]
}

# ----------------------------------------------------------------------------------------------------------------------
# DKIM
# ----------------------------------------------------------------------------------------------------------------------
resource "aws_ses_domain_dkim" "default" {
  domain = aws_ses_domain_identity.default.domain
}

resource "aws_route53_record" "ses_dkim" {
  count = 3

  zone_id = aws_route53_zone.main.zone_id
  type    = "CNAME"
  ttl     = 600
  name    = "${element(aws_ses_domain_dkim.default.dkim_tokens, count.index)}._domainkey"
  records = ["${element(aws_ses_domain_dkim.default.dkim_tokens, count.index)}.dkim.amazonses.com"]
}
