// Copyright (c) eHealth.co.id as PT Aksara Digital Indonesia
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"uptimekuma": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Check for required environment variables for acceptance tests
	// Skip the test if these are not set (standard pattern for optional acceptance tests)
	requiredEnvVars := []string{
		"UPTIMEKUMA_BASE_URL",
		"UPTIMEKUMA_USERNAME",
		"UPTIMEKUMA_PASSWORD",
	}

	for _, env := range requiredEnvVars {
		if v := os.Getenv(env); v == "" {
			t.Skipf("%s environment variable must be set for acceptance tests", env)
		}
	}
}
