package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationPublisherResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: providerConfig + `
resource "dependencytrack_notification_publisher" "test" {
	name               = "Test_Publisher"
	description        = "A test publisher"
	publisher_class    = "org.dependencytrack.notification.publisher.WebhookPublisher"
	template           = "{}"
	template_mime_type = "application/json"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dependencytrack_notification_publisher.test", "id"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "name", "Test_Publisher"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "description", "A test publisher"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "publisher_class", "org.dependencytrack.notification.publisher.WebhookPublisher"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "template", "{}"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "template_mime_type", "application/json"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "default_publisher", "false"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      "dependencytrack_notification_publisher.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing.
			{
				Config: providerConfig + `
resource "dependencytrack_notification_publisher" "test" {
	name               = "Test_Publisher_Updated"
	description        = "An updated test publisher"
	publisher_class    = "org.dependencytrack.notification.publisher.WebhookPublisher"
	template           = "{\"updated\": true}"
	template_mime_type = "application/json"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dependencytrack_notification_publisher.test", "id"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "name", "Test_Publisher_Updated"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "description", "An updated test publisher"),
					resource.TestCheckResourceAttr("dependencytrack_notification_publisher.test", "template", "{\"updated\": true}"),
				),
			},
		},
	})
}
