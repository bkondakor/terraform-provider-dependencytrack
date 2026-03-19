# Look up a built-in notification publisher by name
data "dependencytrack_notification_publisher" "webhook" {
  name = "Outbound Webhook"
}

# Other available publishers:
# "Slack", "Microsoft Teams", "Mattermost", "Email",
# "Console", "Cisco Webex", "Jira"
