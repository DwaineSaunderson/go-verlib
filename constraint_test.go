package verlib_test

import (
	"testing"

	"github.com/DwaineSaunderson/go-verlib"
)

func TestCheckSatisfy(t *testing.T) {
	testCases := []struct {
		v        verlib.Version
		c        verlib.Constraint
		expected bool
	}{
		// EQ: Equal
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 4)), false},
		// NE: Not Equal
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 3)), false},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 4)), true},
		// GT: Greater Than
		{verlib.NewVersion(2, 0, 0), verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 9, 9)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 2, 3)), false},
		// GE: Greater Than or Equal
		{verlib.NewVersion(2, 0, 0), verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 9, 9)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 4)), false},
		// LT: Less Than
		{verlib.NewVersion(1, 9, 9), verlib.NewConstraint(verlib.LT, verlib.NewVersion(2, 0, 0)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 2, 3)), false},
		// LE: Less Than or Equal
		{verlib.NewVersion(1, 9, 9), verlib.NewConstraint(verlib.LE, verlib.NewVersion(2, 0, 0)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 2)), false},
		// GEPessimistic: Pessimistic Greater Than or Equal
		{verlib.NewVersion(1, 8, 999), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 8, 0)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewVersion(1, 2, 5), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 4)), true},
		{verlib.NewVersion(1, 2, 3), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 3, 0)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.c.String()+"#"+tc.v.String(), func(t *testing.T) {
			if result := tc.v.Satisfies(tc.c); result != tc.expected {
				t.Errorf("For version %s and constraint %s, expected %t, but got %t", tc.v.String(), tc.c.String(), tc.expected, result)
			}
		})
	}
}

func TestCheckOverlap(t *testing.T) {
	testCases := []struct {
		c1       verlib.Constraint
		c2       verlib.Constraint
		expected bool
	}{
		// EQ: Equal
		{verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.EQ, verlib.NewVersion(1, 2, 4)), false},
		// NE: Not Equal
		{verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.NE, verlib.NewVersion(1, 2, 4)), false},
		// GT: Greater Than
		{verlib.NewConstraint(verlib.GT, verlib.NewVersion(2, 0, 0)), verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 9, 9)), false},
		{verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.GT, verlib.NewVersion(1, 2, 3)), true},
		// GE: Greater Than or Equal
		{verlib.NewConstraint(verlib.GE, verlib.NewVersion(2, 0, 0)), verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 9, 9)), false},
		{verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 3)), false},
		{verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.GE, verlib.NewVersion(1, 2, 4)), true},
		// LT: Less Than
		{verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 9, 9)), verlib.NewConstraint(verlib.LT, verlib.NewVersion(2, 0, 0)), false},
		{verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.LT, verlib.NewVersion(1, 2, 3)), true},
		// LE: Less Than or Equal
		{verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 9, 9)), verlib.NewConstraint(verlib.LE, verlib.NewVersion(2, 0, 0)), false},
		{verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 3)), false},
		{verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.LE, verlib.NewVersion(1, 2, 2)), true},
		// GEPessimistic: Pessimistic Greater Than or Equal
		{verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 9, 0)), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 8, 0)), false},
		{verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 3)), true},
		{verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 4)), false},
		{verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 2, 3)), verlib.NewConstraint(verlib.GEPessimistic, verlib.NewVersion(1, 3, 0)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.c1.String()+","+tc.c2.String(), func(t *testing.T) {
			if result := tc.c1.Overlaps(tc.c2); result != tc.expected {
				t.Errorf("For constraints %s and %s, expected %t, but got %t", tc.c1.String(), tc.c2.String(), tc.expected, result)
			}
		})
	}
}

func TestConstraintString(t *testing.T) {
	testCases := []struct {
		desc       string
		constraint verlib.Constraint
		wantOutput string
	}{
		{
			desc: "Operator and version are present",
			constraint: verlib.NewConstraint(
				verlib.GE,
				verlib.MustParseVersion("1.2.3"),
			),
			wantOutput: ">= 1.2.3",
		},
		// Add more test cases as required
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.constraint.String()
			if got != tC.wantOutput {
				t.Errorf("Constraint.String(): got %s, want %s", got, tC.wantOutput)
			}
		})
	}
}

func TestConstraintStrictString(t *testing.T) {
	testCases := []struct {
		desc       string
		constraint verlib.Constraint
		wantOutput string
		wantErr    bool
	}{
		{
			desc: "Operator and valid version are present",
			constraint: verlib.NewConstraint(
				verlib.GE,
				verlib.MustParseVersion("1.2.3"),
			),
			wantOutput: ">= 1.2.3",
			wantErr:    false,
		},
		{
			desc: "Constraint with loose version",
			constraint: verlib.NewConstraint(
				verlib.GE,
				verlib.MustParseVersion("1.0.0-alpha.."),
			),
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := tC.constraint.StrictString()
			if (err != nil) != tC.wantErr {
				t.Errorf("Constraint.StrictString() error = %v, wantErr %v", err, tC.wantErr)
				return
			}
			if got != tC.wantOutput {
				t.Errorf("Constraint.StrictString(): got %s, want %s", got, tC.wantOutput)
			}
		})
	}
}

func TestConstraintsString(t *testing.T) {
	testCases := []struct {
		desc        string
		constraints verlib.Constraints
		wantOutput  string
		wantErr     bool
	}{
		{
			desc: "Multiple constraints present",
			constraints: verlib.Constraints{
				verlib.NewConstraint(
					verlib.GE,
					verlib.MustParseVersion("1.2.3"),
				),
				verlib.NewConstraint(
					verlib.LT,
					verlib.MustParseVersion("2.0.0"),
				),
			},
			wantOutput: ">= 1.2.3, < 2.0.0",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := tC.constraints.String()
			if got != tC.wantOutput {
				t.Errorf("Constraints.String(): got %s, want %s", got, tC.wantOutput)
			}
		})
	}
}

func TestStrictConstraintsString(t *testing.T) {
	testCases := []struct {
		desc        string
		constraints verlib.Constraints
		wantOutput  string
		wantErr     bool
	}{
		{
			desc: "Multiple constraints present",
			constraints: verlib.Constraints{
				verlib.NewConstraint(
					verlib.GE,
					verlib.MustParseVersion("1.2.3"),
				),
				verlib.NewConstraint(
					verlib.LT,
					verlib.MustParseVersion("2.0.0"),
				),
			},
			wantOutput: ">= 1.2.3, < 2.0.0",
		},
		{
			desc: "Constraint with loose version",
			constraints: verlib.Constraints{
				verlib.NewConstraint(
					verlib.GE,
					verlib.MustParseVersion("1.0.0-alpha.."),
				),
				verlib.NewConstraint(
					verlib.LT,
					verlib.MustParseVersion("2.0.0"),
				),
			},
			wantErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := tC.constraints.StrictString()
			if err != nil && !tC.wantErr {
				t.Errorf("Constraints.StrictString() error = %v, wantErr %v", err, tC.wantErr)
			}

			if got != tC.wantOutput {
				t.Errorf("Constraints.StrictString(): got %s, want %s", got, tC.wantOutput)
			}
		})
	}
}
