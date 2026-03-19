# Look up the Outbound Webhook publisher
data "dependencytrack_notification_publisher" "webhook" {
  name = "Outbound Webhook"
}

# Event-driven notification rule for new vulnerabilities
resource "dependencytrack_notification_rule" "vuln_webhook" {
  name               = "Vulnerability Webhook"
  scope              = "PORTFOLIO"
  notification_level = "INFORMATIONAL"
  publisher_id       = data.dependencytrack_notification_publisher.webhook.id
  notify_on          = ["NEW_VULNERABILITY", "NEW_VULNERABLE_DEPENDENCY"]
  publisher_config   = jsonencode({ destination = "https://example.com/webhook" })
}

# Scheduled notification rule for daily vulnerability summary
resource "dependencytrack_notification_rule" "daily_summary" {
  name                    = "Daily Vulnerability Summary"
  scope                   = "PORTFOLIO"
  notification_level      = "INFORMATIONAL"
  trigger_type            = "SCHEDULE"
  publisher_id            = data.dependencytrack_notification_publisher.webhook.id
  notify_on               = ["NEW_VULNERABILITIES_SUMMARY"]
  schedule_cron           = "0 8 * * *"
  schedule_skip_unchanged = true
  publisher_config        = jsonencode({ destination = "https://example.com/webhook" })
}
