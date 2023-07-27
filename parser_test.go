package verlib_test

import (
	"reflect"
	"testing"

	"github.com/DwaineSaunderson/go-verlib"
)

func TestParseValidVersion(t *testing.T) {
	type testVersion struct {
		version       string
		major         uint64
		minor         uint64
		patch         uint64
		preRelease    string
		buildMetadata string
	}

	testVersions := []testVersion{
		{"1", 1, 0, 0, "", ""},
		{"1.2", 1, 2, 0, "", ""},
		{"1.2.3-0123.0123", 1, 2, 3, "0123.0123", ""},
		{"1.1.2+.123", 1, 1, 2, "", ".123"},
		{"1.2.3-0123", 1, 2, 3, "0123", ""},
		{"0.0.4", 0, 0, 4, "", ""},
		{"1.2.3", 1, 2, 3, "", ""},
		{"10.20.30", 10, 20, 30, "", ""},
		{"1.1.2-prerelease+meta", 1, 1, 2, "prerelease", "meta"},
		{"1.1.2+meta", 1, 1, 2, "", "meta"},
		{"1.1.2+meta-valid", 1, 1, 2, "", "meta-valid"},
		{"1.0.0-alpha", 1, 0, 0, "alpha", ""},
		{"1.0.0-beta", 1, 0, 0, "beta", ""},
		{"1.0.0-alpha.beta", 1, 0, 0, "alpha.beta", ""},
		{"1.0.0-alpha.beta.1", 1, 0, 0, "alpha.beta.1", ""},
		{"1.0.0-alpha.1", 1, 0, 0, "alpha.1", ""},
		{"1.0.0-alpha0.valid", 1, 0, 0, "alpha0.valid", ""},
		{"1.0.0-alpha.0valid", 1, 0, 0, "alpha.0valid", ""},
		{"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", 1, 0, 0, "alpha-a.b-c-somethinglong", "build.1-aef.1-its-okay"},
		{"1.0.0-rc.1+build.1", 1, 0, 0, "rc.1", "build.1"},
		{"2.0.0-rc.1+build.123", 2, 0, 0, "rc.1", "build.123"},
		{"1.2.3-beta", 1, 2, 3, "beta", ""},
		{"10.2.3-DEV-SNAPSHOT", 10, 2, 3, "DEV-SNAPSHOT", ""},
		{"1.2.3-SNAPSHOT-123", 1, 2, 3, "SNAPSHOT-123", ""},
		{"1.0.0", 1, 0, 0, "", ""},
		{"2.0.0", 2, 0, 0, "", ""},
		{"1.1.7", 1, 1, 7, "", ""},
		{"2.0.0+build.1848", 2, 0, 0, "", "build.1848"},
		{"2.0.1-alpha.1227", 2, 0, 1, "alpha.1227", ""},
		{"1.0.0-alpha+beta", 1, 0, 0, "alpha", "beta"},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12+788", 1, 2, 3, "---RC-SNAPSHOT.12.9.1--.12", "788"},
		{"1.2.3----R-S.12.9.1--.12+meta", 1, 2, 3, "---R-S.12.9.1--.12", "meta"},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12", 1, 2, 3, "---RC-SNAPSHOT.12.9.1--.12", ""},
		{"1.0.0+0.build.1-rc.10000aaa-kk-0.1", 1, 0, 0, "", "0.build.1-rc.10000aaa-kk-0.1"},
		{"18446744073709551615.18446744073709551615.18446744073709551615", 18446744073709551615, 18446744073709551615, 18446744073709551615, "", ""},
		{"1.0.0-0A.is.legal", 1, 0, 0, "0A.is.legal", ""},
		{"1.0.0-alpha_beta", 1, 0, 0, "alpha_beta", ""},
		{"1.0.0-alpha..1", 1, 0, 0, "alpha..1", ""},
		{"1.0.0-alpha...1", 1, 0, 0, "alpha...1", ""},
		{"1.0.0-alpha....1", 1, 0, 0, "alpha....1", ""},
		{"1.0.0-alpha.....1", 1, 0, 0, "alpha.....1", ""},
		{"1.0.0-alpha......1", 1, 0, 0, "alpha......1", ""},
		{"1.0.0-alpha.......1", 1, 0, 0, "alpha.......1", ""},
		{"01.1.1", 1, 1, 1, "", ""},
		{"1.01.1", 1, 1, 1, "", ""},
		{"1.1.01", 1, 1, 1, "", ""},
		{"1.2", 1, 2, 0, "", ""},
		{"1.2.3.DEV", 1, 2, 3, "", ""},
		{"1.2-SNAPSHOT", 1, 2, 0, "SNAPSHOT", ""},
		{"1.2-RC-SNAPSHOT", 1, 2, 0, "RC-SNAPSHOT", ""},
		{"-1.0.3-gamma+b7718", 1, 0, 3, "gamma", "b7718"},
		{"9.8.7+meta+meta", 9, 8, 7, "", "meta+meta"},
		{"9.8.7-whatever+meta+meta", 9, 8, 7, "whatever", "meta+meta"},
	}

	for _, tv := range testVersions {
		t.Run(tv.version, func(t *testing.T) {
			v, err := verlib.ParseVersion(tv.version)
			if err != nil {
				t.Errorf("Failed to parse valid version %q: %v", tv.version, err)
			}
			if v.Major() != tv.major || v.Minor() != tv.minor || v.Patch() != tv.patch ||
				v.PreRelease() != tv.preRelease || v.BuildMetadata() != tv.buildMetadata {
				t.Errorf("ParseVersion(%q) = %v, want major=%v, minor=%v, patch=%v, preRelease=%q, buildMetadata=%q", tv.version, v, tv.major, tv.minor, tv.patch, tv.preRelease, tv.buildMetadata)
			}

			v2 := verlib.MustParseVersion(tv.version)
			if !reflect.DeepEqual(v, v2) {
				t.Errorf("MustParseVersion(%q) != ParseVersion(%q)", tv.version, tv.version)
			}
		})
	}
}

func TestParseInvalidVersion(t *testing.T) {
	invalidVersions := []string{
		"+invalid",
		"-invalid",
		"-invalid+invalid",
		"alpha",
		"alpha.beta",
		"alpha+beta",
		"alpha_beta",
		"alpha.",
		"alpha..",
		"beta",
		"-alpha.",
		"+justmeta",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
	}

	for _, v := range invalidVersions {
		t.Run(v, func(t *testing.T) {
			if _, err := verlib.ParseVersion(v); err == nil {
				t.Errorf("Unexpectedly parsed invalid version %q", v)
			}
			mustPanic(func() { verlib.MustParseVersion(v) }, t)
		})
	}
}

func TestParseValidSemanticVersion(t *testing.T) {
	type testVersion struct {
		version       string
		major         uint64
		minor         uint64
		patch         uint64
		preRelease    string
		buildMetadata string
	}

	testVersions := []testVersion{
		{"0.0.4", 0, 0, 4, "", ""},
		{"1.2.3", 1, 2, 3, "", ""},
		{"10.20.30", 10, 20, 30, "", ""},
		{"1.1.2-prerelease+meta", 1, 1, 2, "prerelease", "meta"},
		{"1.1.2+meta", 1, 1, 2, "", "meta"},
		{"1.1.2+meta-valid", 1, 1, 2, "", "meta-valid"},
		{"1.0.0-alpha", 1, 0, 0, "alpha", ""},
		{"1.0.0-beta", 1, 0, 0, "beta", ""},
		{"1.0.0-alpha.beta", 1, 0, 0, "alpha.beta", ""},
		{"1.0.0-alpha.beta.1", 1, 0, 0, "alpha.beta.1", ""},
		{"1.0.0-alpha.1", 1, 0, 0, "alpha.1", ""},
		{"1.0.0-alpha0.valid", 1, 0, 0, "alpha0.valid", ""},
		{"1.0.0-alpha.0valid", 1, 0, 0, "alpha.0valid", ""},
		{"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", 1, 0, 0, "alpha-a.b-c-somethinglong", "build.1-aef.1-its-okay"},
		{"1.0.0-rc.1+build.1", 1, 0, 0, "rc.1", "build.1"},
		{"2.0.0-rc.1+build.123", 2, 0, 0, "rc.1", "build.123"},
		{"1.2.3-beta", 1, 2, 3, "beta", ""},
		{"10.2.3-DEV-SNAPSHOT", 10, 2, 3, "DEV-SNAPSHOT", ""},
		{"1.2.3-SNAPSHOT-123", 1, 2, 3, "SNAPSHOT-123", ""},
		{"1.0.0", 1, 0, 0, "", ""},
		{"2.0.0", 2, 0, 0, "", ""},
		{"1.1.7", 1, 1, 7, "", ""},
		{"2.0.0+build.1848", 2, 0, 0, "", "build.1848"},
		{"2.0.1-alpha.1227", 2, 0, 1, "alpha.1227", ""},
		{"1.0.0-alpha+beta", 1, 0, 0, "alpha", "beta"},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12+788", 1, 2, 3, "---RC-SNAPSHOT.12.9.1--.12", "788"},
		{"1.2.3----R-S.12.9.1--.12+meta", 1, 2, 3, "---R-S.12.9.1--.12", "meta"},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12", 1, 2, 3, "---RC-SNAPSHOT.12.9.1--.12", ""},
		{"1.0.0+0.build.1-rc.10000aaa-kk-0.1", 1, 0, 0, "", "0.build.1-rc.10000aaa-kk-0.1"},
		{"18446744073709551615.18446744073709551615.18446744073709551615", 18446744073709551615, 18446744073709551615, 18446744073709551615, "", ""},
		{"1.0.0-0A.is.legal", 1, 0, 0, "0A.is.legal", ""},
	}

	for _, tv := range testVersions {
		t.Run(tv.version, func(t *testing.T) {
			v, err := verlib.ParseSemVer(tv.version)
			if err != nil {
				t.Errorf("Failed to parse valid version %q: %v", tv.version, err)
			}
			if v.Major() != tv.major || v.Minor() != tv.minor || v.Patch() != tv.patch ||
				v.PreRelease() != tv.preRelease || v.BuildMetadata() != tv.buildMetadata {
				t.Errorf("ParseVersion(%q) = %v, want major=%v, minor=%v, patch=%v, preRelease=%q, buildMetadata=%q", tv.version, v, tv.major, tv.minor, tv.patch, tv.preRelease, tv.buildMetadata)
			}

			v2 := verlib.MustParseSemVer(tv.version)
			if !reflect.DeepEqual(v, v2) {
				t.Errorf("MustParseSemVer(%q) != ParseSemVer(%q)", tv.version, tv.version)
			}
		})
	}
}

func TestParseInvalidSemanticVersion(t *testing.T) {
	invalidVersions := []string{
		"1",
		"1.2",
		"1.2.3-0123",
		"1.2.3-0123.0123",
		"1.1.2+.123",
		"+invalid",
		"-invalid",
		"-invalid+invalid",
		"-invalid.01",
		"alpha",
		"alpha.beta",
		"alpha.beta.1",
		"alpha.1",
		"alpha+beta",
		"alpha_beta",
		"alpha.",
		"alpha..",
		"beta",
		"1.0.0-alpha_beta",
		"-alpha.",
		"1.0.0-alpha..1",
		"1.0.0-alpha...1",
		"1.0.0-alpha....1",
		"1.0.0-alpha.....1",
		"1.0.0-alpha......1",
		"1.0.0-alpha.......1",
		"01.1.1",
		"1.01.1",
		"1.1.01",
		"1.2",
		"1.2.3.DEV",
		"1.2-SNAPSHOT",
		"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
		"1.2-RC-SNAPSHOT",
		"-1.0.3-gamma+b7718",
		"+justmeta",
		"9.8.7+meta+meta",
		"9.8.7-whatever+meta+meta",
		"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
	}

	for _, v := range invalidVersions {
		t.Run(v, func(t *testing.T) {
			if _, err := verlib.ParseSemVer(v); err == nil {
				t.Errorf("Unexpectedly parsed invalid version %q", v)
			}

			mustPanic(func() { verlib.MustParseSemVer(v) }, t)
		})
	}
}

func TestParseConstraint(t *testing.T) {
	type testCase struct {
		input          string
		expectedOutput verlib.Constraint
		expectError    bool
	}

	testCases := []testCase{
		// Test cases for no operator
		{input: "1", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: "1.2", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: "1.2.3", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: "1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for "=" operator
		{input: "=1", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: "=1.2", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: "=1.2.3", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: "=1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.EQ, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for "!=" operator
		{input: "!=1", expectedOutput: verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: "!=1.2", expectedOutput: verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: "!=1.2.3", expectedOutput: verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: "!=1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.NE, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for "<=" operator
		{input: "<=1", expectedOutput: verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: "<=1.2", expectedOutput: verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: "<=1.2.3", expectedOutput: verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: "<=1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.LE, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for ">=" operator
		{input: ">=1", expectedOutput: verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: ">=1.2", expectedOutput: verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: ">=1.2.3", expectedOutput: verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: ">=1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.GE, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for "<" operator
		{input: "<1", expectedOutput: verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: "<1.2", expectedOutput: verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: "<1.2.3", expectedOutput: verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: "<1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.LT, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for ">" operator
		{input: ">1", expectedOutput: verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: ">1.2", expectedOutput: verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: ">1.2.3", expectedOutput: verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: ">1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.GT, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},

		// Test cases for "~>" operator
		{input: "~>1", expectedOutput: verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 0, 0)), expectError: false},
		{input: "~>1.2", expectedOutput: verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 0)), expectError: false},
		{input: "~>1.2.3", expectedOutput: verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 3)), expectError: false},
		{input: "~>1.2.3-alpha", expectedOutput: verlib.NewConstraint(verlib.GEPessimistic, verlib.NewPreReleaseVersion(1, 2, 3, "alpha")), expectError: false},
	}

	for _, test := range testCases {
		t.Run(test.input, func(t *testing.T) {
			result, err := verlib.ParseConstraint(test.input)
			if (err != nil) != test.expectError {
				t.Errorf("Unexpected error from ParseConstraint(%q): got %v, want %v", test.input, err, test.expectError)
			}

			resultStr, err := result.StrictString()
			if err != nil {
				t.Fatal(err)
			}
			expectedStr, err := test.expectedOutput.StrictString()
			if err != nil {
				t.Fatal(err)
			}

			if resultStr != expectedStr {
				t.Errorf("Unexpected result from ParseConstraint(%q): got %v, want %v", test.input, resultStr, expectedStr)
			}

			mustParseResult := verlib.MustParseConstraint(test.input)

			mustParseResultStr, err := mustParseResult.StrictString()
			if err != nil {
				t.Fatal(err)
			}
			if mustParseResultStr != expectedStr {
				t.Errorf("Unexpected result from MustParseConstraint(%q): got %v, want %v", test.input, mustParseResultStr, expectedStr)
			}
		})
	}
}

func TestParseStrictConstraint(t *testing.T) {
	tests := []struct {
		name        string
		constraint  string
		expectedStr string
		expectErr   bool
	}{
		{
			name:        "EqualMajor",
			constraint:  "= 1.0.0",
			expectedStr: "= 1.0.0",
			expectErr:   false,
		},
		{
			name:        "EqualMajorNoOperator",
			constraint:  "1.0.0",
			expectedStr: "= 1.0.0",
			expectErr:   false,
		},
		{
			name:        "NotEqualMajorMinor",
			constraint:  "!= 2.1.0",
			expectedStr: "!= 2.1.0",
			expectErr:   false,
		},
		{
			name:        "GreaterThanMajorMinorPatch",
			constraint:  "> 3.2.1",
			expectedStr: "> 3.2.1",
			expectErr:   false,
		},
		{
			name:        "GreaterThanOrEqualMajorMinorPatchPre",
			constraint:  ">= 4.3.2-alpha",
			expectedStr: ">= 4.3.2-alpha",
			expectErr:   false,
		},
		{
			name:        "LessThanMajorMinorPatchPreBuild",
			constraint:  "< 5.4.3-beta+20230727",
			expectedStr: "< 5.4.3-beta+20230727",
			expectErr:   false,
		},
		{
			name:        "LessThanOrEqualPatch",
			constraint:  "<= 6.0.0",
			expectedStr: "<= 6.0.0",
			expectErr:   false,
		},
		{
			name:        "PessimisticMajorMinor",
			constraint:  "~> 7.1.0",
			expectedStr: "~> 7.1.0",
			expectErr:   false,
		},
		{
			name:        "InvalidOperator",
			constraint:  "? 1.0.0",
			expectedStr: "",
			expectErr:   true,
		},
		{
			name:        "InvalidVersion",
			constraint:  "> 2.a.b",
			expectedStr: "",
			expectErr:   true,
		},
		{
			name:        "LeadingWhitespace",
			constraint:  "   < 2.0.0",
			expectedStr: "< 2.0.0",
			expectErr:   false,
		},
		{
			name:        "TrailingWhitespace",
			constraint:  ">= 3.0.0   ",
			expectedStr: ">= 3.0.0",
			expectErr:   false,
		},
		{
			name:        "WhitespaceAroundOperator",
			constraint:  "  <=   4.0.0",
			expectedStr: "<= 4.0.0",
			expectErr:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedConstraint, err := verlib.ParseStrictConstraint(tc.constraint)

			if (err != nil) != tc.expectErr {
				t.Errorf("ParseStrictConstraint() error = %v, expectErr %v", err, tc.expectErr)
				return
			}

			if !tc.expectErr {
				str, err := parsedConstraint.StrictString()
				if err != nil {
					t.Errorf("failed to convert parsed constraint back to string: %v", err)
					return
				}
				if str != tc.expectedStr {
					t.Errorf("got %s, want %s", str, tc.expectedStr)
				}
			}
		})
	}
}

func TestParseConstraintSet(t *testing.T) {
	testCases := []struct {
		name             string
		constraintString string
		expected         []string
		expectErr        bool
	}{
		{
			name:             "Single constraint major version only",
			constraintString: "= 1.0.0",
			expected:         []string{"= 1.0.0"},
			expectErr:        false,
		},
		{
			name:             "Single constraint major version only, no operator",
			constraintString: "1.0.0",
			expected:         []string{"= 1.0.0"},
			expectErr:        false,
		},
		{
			name:             "Single constraint major and minor version",
			constraintString: "<=1.2",
			expected:         []string{"<= 1.2.0"},
			expectErr:        false,
		},
		{
			name:             "Single constraint with major, minor and patch version",
			constraintString: ">1.2.3",
			expected:         []string{"> 1.2.3"},
			expectErr:        false,
		},
		{
			name:             "Single constraint with prerelease",
			constraintString: "<1.2.3-alpha",
			expected:         []string{"< 1.2.3-alpha"},
			expectErr:        false,
		},
		{
			name:             "Multiple constraints with different operators and versions",
			constraintString: ">=1.0.0, !=2.0.1-alpha, <1.2.3, >1, <=1.2",
			expected:         []string{">= 1.0.0", "!= 2.0.1-alpha", "< 1.2.3", "> 1.0.0", "<= 1.2.0"},
			expectErr:        false,
		},
		{
			name:             "Invalid constraint",
			constraintString: "abc",
			expectErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			constraints, err := verlib.ParseConstraintSet(tc.constraintString)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				mustPanic(func() { verlib.MustParseConstraintSet(tc.constraintString) }, t)
				return
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			for i, constraint := range constraints {
				str, err := constraint.StrictString()
				if err != nil {
					t.Errorf("Failed to generate strict string: %v", err)
					return
				}

				if str != tc.expected[i] {
					t.Errorf("Expected constraint %s, got %s", tc.expected[i], str)
				}
			}

			constraints2 := verlib.MustParseConstraintSet(tc.constraintString)
			for i, constraint := range constraints2 {
				str, err := constraint.StrictString()
				if err != nil {
					t.Errorf("(MustParse) Failed to generate strict string: %v", err)
					return
				}

				if str != tc.expected[i] {
					t.Errorf("(MustParse) Expected constraint %s, got %s", tc.expected[i], str)
				}
			}
		})
	}
}

func TestParseStrictConstraintSet(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "single strict constraint",
			input:          ">= 2.0.0",
			expectedOutput: ">= 2.0.0",
			expectError:    false,
		},
		{
			name:           "multiple strict constraints",
			input:          ">= 2.0.0, < 3.0.0",
			expectedOutput: ">= 2.0.0, < 3.0.0",
			expectError:    false,
		},
		{
			name:           "invalid strict constraint",
			input:          "!= 2.0.x",
			expectedOutput: "",
			expectError:    true,
		},
		{
			name:           "multiple constraints with one invalid",
			input:          ">= 2.0.0, != 2.0.x",
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			constraints, err := verlib.ParseStrictConstraintSet(tc.input)
			if (err != nil) != tc.expectError {
				t.Fatalf("expected error: %v, got: %v", tc.expectError, err)
			}

			if !tc.expectError {
				strictString, _ := constraints.StrictString()
				if strictString != tc.expectedOutput {
					t.Errorf("expected output: %q, got: %q", tc.expectedOutput, strictString)
				}
			}
		})
	}
}

func TestMustParseStrictConstraintSet(t *testing.T) {
	t.Run("must panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		invalidConstraint := "!= 2.0.x"
		verlib.MustParseStrictConstraintSet(invalidConstraint) // This should panic
	})

	t.Run("should succeed", func(t *testing.T) {
		validConstraint := "!= 2.0.0"
		constraints := verlib.MustParseStrictConstraintSet(validConstraint) // This should panic

		strictString, _ := constraints.StrictString()
		if strictString != validConstraint {
			t.Errorf("expected output: %q, got: %q", validConstraint, strictString)
		}
	})
}

func mustPanic(f func(), t *testing.T) {
	var didPanic bool

	defer func() {
		if !didPanic {
			t.Fatal("expected to panic but ran successfully")
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
		}
	}()
	f()
}
