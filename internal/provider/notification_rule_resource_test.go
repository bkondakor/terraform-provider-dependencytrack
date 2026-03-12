package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
data "dependencytrack_notification_publisher" "webhook" {
	name = "Outbound Webhook"
}

resource "dependencytrack_notification_rule" "test" {
	name              = "Test_Notification_Rule"
	scope             = "PORTFOLIO"
	notification_level = "INFORMATIONAL"
	publisher_id      = data.dependencytrack_notification_publisher.webhook.id
	notify_on         = ["NEW_VULNERABILITY", "BOM_CONSUMED"]
	publisher_config  = "{\"destination\":\"https://example.com/webhook\"}"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dependencytrack_notification_rule.test", "id"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "name", "Test_Notification_Rule"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "scope", "PORTFOLIO"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "notification_level", "INFORMATIONAL"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "enabled", "true"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "notify_children", "true"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "log_successful_publish", "false"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "trigger_type", "EVENT"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "publisher_config", "{\"destination\":\"https://example.com/webhook\"}"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "dependencytrack_notification_rule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
data "dependencytrack_notification_publisher" "webhook" {
	name = "Outbound Webhook"
}

resource "dependencytrack_notification_rule" "test" {
	name              = "Test_Notification_Rule_Updated"
	scope             = "PORTFOLIO"
	notification_level = "WARNING"
	publisher_id      = data.dependencytrack_notification_publisher.webhook.id
	enabled           = false
	notify_children   = false
	notify_on         = ["NEW_VULNERABILITY", "POLICY_VIOLATION"]
	publisher_config  = "{\"destination\":\"https://example.com/webhook2\"}"
	message           = "Custom notification message"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dependencytrack_notification_rule.test", "id"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "name", "Test_Notification_Rule_Updated"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "notification_level", "WARNING"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "enabled", "false"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "notify_children", "false"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "message", "Custom notification message"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test", "publisher_config", "{\"destination\":\"https://example.com/webhook2\"}"),
				),
			},
		},
	})
}

func TestAccNotificationRuleScheduledResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing for scheduled rule.
			{
				Config: providerConfig + `
data "dependencytrack_notification_publisher" "webhook" {
	name = "Outbound Webhook"
}

resource "dependencytrack_notification_rule" "test_scheduled" {
	name              = "Test_Scheduled_Rule"
	scope             = "PORTFOLIO"
	notification_level = "INFORMATIONAL"
	trigger_type      = "SCHEDULE"
	publisher_id      = data.dependencytrack_notification_publisher.webhook.id
	notify_on         = ["NEW_VULNERABILITIES_SUMMARY"]
	schedule_cron     = "0 0 * * *"
	schedule_skip_unchanged = true
	publisher_config  = "{\"destination\":\"https://example.com/webhook\"}"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dependencytrack_notification_rule.test_scheduled", "id"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test_scheduled", "name", "Test_Scheduled_Rule"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test_scheduled", "trigger_type", "SCHEDULE"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test_scheduled", "schedule_cron", "0 0 * * *"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule.test_scheduled", "schedule_skip_unchanged", "true"),
				),
			},
		},
	})
}
