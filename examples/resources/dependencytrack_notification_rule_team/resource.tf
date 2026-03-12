# Associate a team with an email notification rule
resource "dependencytrack_notification_rule_team" "example" {
  rule = dependencytrack_notification_rule.email_alerts.id
  team = dependencytrack_team.security.id
}
