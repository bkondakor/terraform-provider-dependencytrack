package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationRuleTagResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
data "dependencytrack_notification_publisher" "webhook" {
	name = "Outbound Webhook"
}

resource "dependencytrack_tag" "test" {
	name = "test_notification_rule_tag"
}

resource "dependencytrack_notification_rule" "test" {
	name              = "Test_Rule_Tag"
	scope             = "PORTFOLIO"
	notification_level = "INFORMATIONAL"
	publisher_id      = data.dependencytrack_notification_publisher.webhook.id
}

resource "dependencytrack_notification_rule_tag" "test" {
	rule = dependencytrack_notification_rule.test.id
	tag  = dependencytrack_tag.test.name
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("dependencytrack_notification_rule_tag.test", "rule", "dependencytrack_notification_rule.test", "id"),
					resource.TestCheckResourceAttr("dependencytrack_notification_rule_tag.test", "tag", "test_notification_rule_tag"),
				),
			},
		},
	})
}
