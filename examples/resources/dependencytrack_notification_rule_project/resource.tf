# Scope a notification rule to a specific project
resource "dependencytrack_notification_rule_project" "example" {
  rule    = dependencytrack_notification_rule.vuln_webhook.id
  project = dependencytrack_project.example.id
}
