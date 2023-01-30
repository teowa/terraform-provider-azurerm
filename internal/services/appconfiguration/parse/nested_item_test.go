package parse

import "testing"

func TestNewNestedItemID(t *testing.T) {
	cases := []struct {
		Scenario                  string
		ConfigurationStoreBaseURL string
		Key                       string
		Label                     string
		Expected                  string
		ExpectError               bool
	}{
		{
			ConfigurationStoreBaseURL: "",
			Key:                       "testKey",
			Label:                     "testLabel",
			Expected:                  "",
			ExpectError:               true,
		},
		{
			ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
			Key:                       "testKey",
			Label:                     "testLabel",
			Expected:                  "https://testappconf1.azconfig.io/kv/testKey?label=testLabel",
			ExpectError:               false,
		},
		{
			ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
			Key:                       "test+/123",
			Label:                     "testLabel",
			Expected:                  "https://testappconf1.azconfig.io/kv/test+%2F123?label=testLabel",
			ExpectError:               false,
		},
		{
			ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
			Key:                       "testKey",
			Label:                     "test+/123",
			Expected:                  "https://testappconf1.azconfig.io/kv/testKey?label=test%2B%2F123",
			ExpectError:               false,
		},
		{
			ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
			Key:                       "testKey",
			Label:                     "",
			Expected:                  "https://testappconf1.azconfig.io/kv/testKey",
			ExpectError:               false,
		},
		{
			ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
			Key:                       "testKey",
			Label:                     "%00",
			Expected:                  "https://testappconf1.azconfig.io/kv/testKey?label=%2500",
			ExpectError:               false,
		},
	}
	for _, tc := range cases {
		id, err := NewNestedItemID(tc.ConfigurationStoreBaseURL, tc.Key, tc.Label)
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for New Nested Item ID (BaseURL:%q, Key:%q, Label:%q): %+v", tc.ConfigurationStoreBaseURL, tc.Key, tc.Label, err)
				return
			}
			continue
		}
		if id.ID() != tc.Expected {
			t.Fatalf("Expected id for (BaseURL:%q, Key:%q, Label:%q) to be %q, got %q", tc.ConfigurationStoreBaseURL, tc.Key, tc.Label, tc.Expected, id.ID())
		}
	}
}

func TestParseNestedItemID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    NestedItemId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/?label=testLabel",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=testLabel",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "testLabel",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/test+%2F123?label=test%2B%2F123",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "test+/123",
				Label:                     "test+/123",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?b=",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=a&b=c",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=a&label=b",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=a%2Cb%2Cc",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "a,b,c",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=%00",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=%2500",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       "testKey",
				Label:                     "%00",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/.appconfig.featureflag%2FtestKey?label=testLabel",
			ExpectError: false,
			Expected: NestedItemId{
				ConfigurationStoreBaseURL: "https://testappconf1.azconfig.io",
				Key:                       ".appconfig.featureflag/testKey",
				Label:                     "testLabel",
			},
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/.appconfig.featureflag/testKey?label=testLabel",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		nestedItemId, err := ParseNestedItemID(tc.Input)
		if err != nil {
			if tc.ExpectError {
				continue
			}

			t.Fatalf("Got error for ID %q: %+v", tc.Input, err)
		}

		if nestedItemId == nil {
			t.Fatalf("Expected a nestedItemID to be parsed for ID %q, got nil.", tc.Input)
		}

		if tc.Expected.ConfigurationStoreBaseURL != nestedItemId.ConfigurationStoreBaseURL {
			t.Fatalf("Expected ConfigurationStoreBaseURL to be %q, got %q for ID %q", tc.Expected.ConfigurationStoreBaseURL, nestedItemId.ConfigurationStoreBaseURL, tc.Input)
		}

		if tc.Expected.Key != nestedItemId.Key {
			t.Fatalf("Expected Key to be %q, got %q for ID %q", tc.Expected.Key, nestedItemId.Key, tc.Input)
		}

		if tc.Expected.Label != nestedItemId.Label {
			t.Fatalf("Expected Label to be %q, got %q for ID %q", tc.Expected.Label, nestedItemId.Label, tc.Input)
		}
	}
}
