package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMonitorResourceConfig("test-monitor"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptimekuma_monitor.test", "name", "test-monitor"),
					resource.TestCheckResourceAttr("uptimekuma_monitor.test", "type", "http"),
					resource.TestCheckResourceAttr("uptimekuma_monitor.test", "url", "https://example.com"),
					resource.TestCheckResourceAttrSet("uptimekuma_monitor.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "uptimekuma_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccMonitorResourceConfig("test-monitor-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptimekuma_monitor.test", "name", "test-monitor-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMonitorResourceConfig(name string) string {
	return `
provider "uptimekuma" {
  url      = "http://localhost:3001"
  username = "admin"
  password = "test123"
}

resource "uptimekuma_monitor" "test" {
  name = "` + name + `"
  type = "http"
  url  = "https://example.com"
  
  interval = 60
  timeout  = 30
  active   = true
}
`
}
