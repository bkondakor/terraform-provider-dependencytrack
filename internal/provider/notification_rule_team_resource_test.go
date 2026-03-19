package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationRuleTeamResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
data "dependencytrack_notification_publisher" "email" {
	name = "Email"
}

resource "dependencytrack_notification_rule" "test" {
	name              = "Test_Rule_Team"
	scope             = "PORTFOLIO"
	notification_level = "INFORMATIONAL"
	publisher_id      = data.dependencytrack_notification_publisher.email.id
}

resource "dependencytrack_team" "test" {
	name = "Test_Notification_Team"
}

resource "dependencytrack_notification_rule_team" "test" {
	rule = dependencytrack_notification_rule.test.id
	team = dependencytrack_team.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("dependencytrack_notification_rule_team.test", "rule", "dependencytrack_notification_rule.test", "id"),
					resource.TestCheckResourceAttrPair("dependencytrack_notification_rule_team.test", "team", "dependencytrack_team.test", "id"),
				),
			},
		},
	})
}
