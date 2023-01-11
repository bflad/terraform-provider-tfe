package tfe

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestPluginProvider_providerMeta(t *testing.T) {
	cases := map[string]struct {
		hostname            string
		token               string
		sslSkipVerify       bool
		defaultOrganization string
		err                 error
	}{
		"has none": {},
		"has only hostname": {
			hostname: "terraform.io",
		},
		"has only token": {
			token: "secret",
		},
		"has only ssl_skip_verify": {
			sslSkipVerify: true,
		},
		"has hostname and token": {
			hostname: "terraform.io",
			token:    "secret",
		},
		"has hostname and ssl_skip_verify": {
			hostname:      "terraform.io",
			sslSkipVerify: true,
		},
		"has token and ssl_skip_verify": {
			token:         "secret",
			sslSkipVerify: true,
		},
		"has default_organization": {
			defaultOrganization: "hashicorp",
		},
	}

	for name, tc := range cases {
		config, err := tfprotov5.NewDynamicValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"hostname":             tftypes.String,
				"token":                tftypes.String,
				"ssl_skip_verify":      tftypes.Bool,
				"default_organization": tftypes.String,
			},
		}, tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"hostname":             tftypes.String,
				"token":                tftypes.String,
				"ssl_skip_verify":      tftypes.Bool,
				"default_organization": tftypes.String,
			},
		}, map[string]tftypes.Value{
			"hostname":             tftypes.NewValue(tftypes.String, tc.hostname),
			"token":                tftypes.NewValue(tftypes.String, tc.token),
			"ssl_skip_verify":      tftypes.NewValue(tftypes.Bool, tc.sslSkipVerify),
			"default_organization": tftypes.NewValue(tftypes.String, tc.defaultOrganization),
		}))

		req := &tfprotov5.ConfigureProviderRequest{
			Config: &config,
		}

		meta, err := retrieveProviderMeta(req)
		if err != tc.err {
			t.Fatalf("Test %s: should not be error, got %v", name, err)
		}

		if tc.hostname == "" && meta.hostname != "" {
			t.Fatalf("Test %s: hostname was not set in config and meta hostname should be empty in this moment (in retrieveProviderMeta). It is parsed later in within the `getClient` function", name)
		}

		if tc.hostname != "" && meta.hostname != tc.hostname {
			t.Fatalf("Test %s: hostname was set in config and meta hostname %s  has not been set to what was given %s", name, meta.hostname, tc.hostname)
		}

		if tc.token == "" && meta.token != "" {
			t.Fatalf("Test %s: token was not set in config and meta.token %s has been incorrectly set", name, meta.token)
		}

		if tc.token != "" && meta.token != tc.token {
			t.Fatalf("Test %s: token was set in config and input token %s  does not have the same value in meta %s", name, tc.token, meta.token)
		}

		if tc.sslSkipVerify == false && meta.sslSkipVerify != defaultSSLSkipVerify {
			t.Fatalf("Test %s: ssl_skip_verify was not set in config and has not been set to default", name)
		}

		if tc.defaultOrganization != meta.defaultOrganization {
			t.Fatalf("Test %s: default organization was set in config and input default organization %s does not have the same value in meta %s", name, tc.token, meta.token)
		}

		if tc.sslSkipVerify != false {
			if meta.sslSkipVerify != tc.sslSkipVerify {
				t.Fatalf("Test %s: ssl_skip_verify was set in config but does not have the same value in meta %t", name, meta.sslSkipVerify)
			}
		}
	}
}
