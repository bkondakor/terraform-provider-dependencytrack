package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationPublisherDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "dependencytrack_notification_publisher" "webhook" {
	name = "Outbound Webhook"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.dependencytrack_notification_publisher.webhook", "id"),
					resource.TestCheckResourceAttr("data.dependencytrack_notification_publisher.webhook", "name", "Outbound Webhook"),
					resource.TestCheckResourceAttrSet("data.dependencytrack_notification_publisher.webhook", "publisher_class"),
					resource.TestCheckResourceAttr("data.dependencytrack_notification_publisher.webhook", "default_publisher", "true"),
				),
			},
		},
	})
}
