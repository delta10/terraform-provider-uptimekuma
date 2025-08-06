package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/j0r15/uptime-kuma-terraform-provider/internal/provider"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"uptimekuma": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("UPTIMEKUMA_URL"); v == "" {
		t.Fatal("UPTIMEKUMA_URL must be set for acceptance tests")
	}
	
	if v := os.Getenv("UPTIMEKUMA_USERNAME"); v == "" {
		t.Fatal("UPTIMEKUMA_USERNAME must be set for acceptance tests")
	}
	
	if v := os.Getenv("UPTIMEKUMA_PASSWORD"); v == "" {
		t.Fatal("UPTIMEKUMA_PASSWORD must be set for acceptance tests")
	}
}
