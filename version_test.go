package verlib_test

import (
	"fmt"
	"testing"

	"github.com/DwaineSaunderson/go-verlib"
)

func TestNewVersion(t *testing.T) {
	v := verlib.NewVersion(1, 2, 3)

	if v.Major() != 1 {
		t.Errorf("Expected major to be 1, got %d", v.Major())
	}

	if v.Minor() != 2 {
		t.Errorf("Expected minor to be 2, got %d", v.Minor())
	}

	if v.Patch() != 3 {
		t.Errorf("Expected patch to be 3, got %d", v.Patch())
	}

	if v.PreRelease() != "" {
		t.Errorf("Expected pre-release to be empty, got %s", v.PreRelease())
	}

	if v.BuildMetadata() != "" {
		t.Errorf("Expected build metadata to be empty, got %s", v.BuildMetadata())
	}
}

func TestNewPreReleaseVersion(t *testing.T) {
	v := verlib.NewPreReleaseVersion(1, 2, 3, "alpha")

	if v.Major() != 1 {
		t.Errorf("Expected major to be 1, got %d", v.Major())
	}

	if v.Minor() != 2 {
		t.Errorf("Expected minor to be 2, got %d", v.Minor())
	}

	if v.Patch() != 3 {
		t.Errorf("Expected patch to be 3, got %d", v.Patch())
	}

	if v.PreRelease() != "alpha" {
		t.Errorf("Expected pre-release to be alpha, got %s", v.PreRelease())
	}

	if v.BuildMetadata() != "" {
		t.Errorf("Expected build metadata to be empty, got %s", v.BuildMetadata())
	}
}

func TestSetPreRelease(t *testing.T) {
	v := verlib.NewVersion(1, 2, 3)
	v = v.SetPreRelease("beta")

	if v.PreRelease() != "beta" {
		t.Errorf("Expected pre-release to be beta, got %s", v.PreRelease())
	}
}

func TestSetBuildMetadata(t *testing.T) {
	v := verlib.NewVersion(1, 2, 3)
	v = v.SetBuildMetadata("2023.07.27")

	if v.BuildMetadata() != "2023.07.27" {
		t.Errorf("Expected build metadata to be 2023.07.27, got %s", v.BuildMetadata())
	}
}

func TestIncrementMajor(t *testing.T) {
	v := verlib.NewVersion(1, 2, 3)
	v = v.IncrementMajor()

	if v.Major() != 2 {
		t.Errorf("Expected major to be 2 after increment, got %d", v.Major())
	}

	if v.Minor() != 0 {
		t.Errorf("Expected minor to be 0 after major increment, got %d", v.Minor())
	}

	if v.Patch() != 0 {
		t.Errorf("Expected patch to be 0 after major increment, got %d", v.Patch())
	}
}

func TestString(t *testing.T) {
	v := verlib.NewVersion(1, 2, 3)
	v = v.SetPreRelease("alpha")
	v = v.SetBuildMetadata("2023.07.27")

	expected := "1.2.3-alpha+2023.07.27"
	if v.String() != expected {
		t.Errorf("Expected version string to be %s, got %s", expected, v.String())
	}
}

func TestLess(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2", "2.1", true},
		{"2.1", "2.1.0", false},
		{"2.1-alpha", "2.1", true},
		{"2.1", "2.1-alpha", false},
		{"3.4.5", "3.4", false},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%q<%q", tc.v1, tc.v2), func(t *testing.T) {
			v1 := verlib.MustParseVersion(tc.v1)
			v2 := verlib.MustParseVersion(tc.v2)
			if v1.Less(v2) != tc.expected {
				t.Errorf("Expected '%s'.Less(%s) to be %t", tc.v1, tc.v2, tc.expected)
			}
		})
	}
}

func TestStrictString(t *testing.T) {
	cases := []struct {
		v           string
		expected    string
		expectError bool
	}{
		{"1-alpha+foo", "1.0.0-alpha+foo", false},
		{"1.0.0-alpha..", "", true},
		{"1.0.0+_", "", true},
		{"1.2", "1.2.0", false},
		{"2.1-alpha", "2.1.0-alpha", false},
		{"3.4.5", "3.4.5", false},
	}

	for _, tc := range cases {
		t.Run(tc.v, func(t *testing.T) {
			v := verlib.MustParseVersion(tc.v)
			str, err := v.StrictString()

			if tc.expectError && err == nil {
				t.Fatalf("expected error: %v, got: %v", tc.expectError, err)
			}
			if str != tc.expected {
				t.Errorf("Expected %s.StrictString() to be %s, got %s", tc.v, tc.expected, str)
			}
		})
	}
}

func TestIncrement(t *testing.T) {
	cases := []struct {
		v        string
		expected string
	}{
		{"1", "2.0.0"},
		{"1.2", "1.3.0"},
		{"2.1-alpha", "2.2.0"},
		{"3.4.5", "3.4.6"},
	}

	for _, tc := range cases {
		v := verlib.MustParseVersion(tc.v).Increment()
		if str, _ := v.StrictString(); str != tc.expected {
			t.Errorf("Expected %s.Increment().StrictString() to be %s, got %s", tc.v, tc.expected, str)
		}
	}
}

func TestIncrementPatch(t *testing.T) {
	cases := []struct {
		v        string
		expected string
	}{
		{"1", "1.0.1"},
		{"1.2", "1.2.1"},
		{"2.1-alpha", "2.1.1"},
		{"3.4.5", "3.4.6"},
	}

	for _, tc := range cases {
		v := verlib.MustParseVersion(tc.v).IncrementPatch()
		if str, _ := v.StrictString(); str != tc.expected {
			t.Errorf("Expected %s.IncrementPatch().StrictString() to be %s, got %s", tc.v, tc.expected, str)
		}
	}
}

func TestIncrementMinor(t *testing.T) {
	cases := []struct {
		v        string
		expected string
	}{
		{"1", "1.1.0"},
		{"1.2", "1.3.0"},
		{"2.1-alpha", "2.2.0"},
		{"3.4.5", "3.5.0"},
	}

	for _, tc := range cases {
		v := verlib.MustParseVersion(tc.v).IncrementMinor()
		if str, _ := v.StrictString(); str != tc.expected {
			t.Errorf("Expected %s.IncrementMinor().StrictString() to be %s, got %s", tc.v, tc.expected, str)
		}
	}
}

func TestIncrementPessimistic(t *testing.T) {
	testCases := []struct {
		desc       string
		input      string
		wantOutput string
	}{
		{
			desc:       "major version only",
			input:      "1",
			wantOutput: "2.0.0",
		},
		{
			desc:       "major and minor version",
			input:      "2.1",
			wantOutput: "3.0.0",
		},
		{
			desc:       "major, minor, and patch version",
			input:      "2.1.5",
			wantOutput: "2.2.0",
		},
		{
			desc:       "version with pre-release label",
			input:      "4.0.2-alpha.1",
			wantOutput: "4.1.0",
		},
		{
			desc:       "version with pre-release label and build metadata",
			input:      "2.2.1-beta.2+20230101",
			wantOutput: "2.3.0",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			inputVersion := verlib.MustParseVersion(tC.input)
			got := inputVersion.IncrementPessimistic()

			outputString, err := got.StrictString()
			if err != nil {
				t.Fatal(err)
			}

			if outputString != tC.wantOutput {
				t.Errorf("IncrementPessimistic(%s): got %s, want %s", tC.input, outputString, tC.wantOutput)
			}
		})
	}
}
