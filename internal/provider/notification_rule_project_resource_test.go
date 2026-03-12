package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationRuleProjectResource(t *testing.T) {
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
	name              = "Test_Rule_Project"
	scope             = "PORTFOLIO"
	notification_level = "INFORMATIONAL"
	publisher_id      = data.dependencytrack_notification_publisher.webhook.id
}

resource "dependencytrack_project" "test" {
	name = "Test_Notification_Project"
}

resource "dependencytrack_notification_rule_project" "test" {
	rule    = dependencytrack_notification_rule.test.id
	project = dependencytrack_project.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("dependencytrack_notification_rule_project.test", "rule", "dependencytrack_notification_rule.test", "id"),
					resource.TestCheckResourceAttrPair("dependencytrack_notification_rule_project.test", "project", "dependencytrack_project.test", "id"),
				),
			},
		},
	})
}
