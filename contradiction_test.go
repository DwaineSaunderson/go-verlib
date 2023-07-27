package verlib_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DwaineSaunderson/go-verlib"
)

func TestIsContradictory(t *testing.T) {
	tests := []struct {
		name     string
		set      verlib.Constraints
		extraSet verlib.Constraints
		wantErr  bool
	}{
		{
			name: "Test_equals_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
			},
			extraSet: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_equals_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("=", verlib.NewVersion(2, 3, 4)),
			},
			wantErr: true,
		},
		{
			name: "Test_equals_not_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("!=", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_equals_not_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
			},
			wantErr: true,
		},
		{
			name: "Test_equals_greater_than",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 0)),
			},
		},
		{
			name: "Test_equals_greater_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">", verlib.NewVersion(2, 3, 4)),
			},
			wantErr: true,
		},
		{
			name: "Test_equals_greater_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">=", verlib.NewVersion(2, 3, 4)),
			},
			wantErr: true,
		},
		{
			name: "Test_equals_less_than",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_equals_less_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_equals_less_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_equals_less_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_equals_pessimistic_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		// Tests for "!=" operator
		{
			name: "Test_not_equals_not_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("!=", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_not_equals_greater_than",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_not_equals_greater_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_not_equals_greater_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">=", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_not_equals_greater_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
			},
			wantErr: true,
		},
		{
			name: "Test_not_equals_less_than",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_not_equals_less_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 4)),
			},
			wantErr: true,
		},
		{
			name: "Test_not_equals_less_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 2)),
			},
		},
		{
			name: "Test_not_equals_less_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 3)),
			},
			wantErr: true,
		},
		{
			name: "Test_not_equals_pessimistic",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 3, 0)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 2, 0)),
			},
		},
		{
			name: "Test_not_equals_pessimistic_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("!=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 2, 3)),
			},
			wantErr: true,
		},
		// Tests for ">" operator
		{
			name: "Test_greater_than_greater_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">", verlib.NewVersion(2, 3, 4)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_greater_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 0)),
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 1)),
			},
		},
		{
			name: "Test_greater_than_greater_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_less_than",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_greater_than_less_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_less_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_greater_than_less_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_pessimistic",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_greater_than_pessimistic_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 3, 0)),
			},
			wantErr: true,
		},
		// Tests for ">=" operator
		{
			name: "Test_greater_than_equals_greater_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_greater_than_equals_greater_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint(">=", verlib.NewVersion(2, 3, 4)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_equals_less_than",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_greater_than_equals_less_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_equals_less_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_greater_than_equals_less_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_greater_than_equals_pessimistic",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_greater_than_equals_pessimistic_sad",
			set: verlib.Constraints{
				verlib.NewConstraint(">=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 1, 0)),
			},
			wantErr: true,
		},
		// Tests for "<" operator
		{
			name: "Test_less_than_less_than",
			set: verlib.Constraints{
				verlib.NewConstraint("<", verlib.NewVersion(2, 3, 4)),
				verlib.NewConstraint("<", verlib.NewVersion(2, 3, 4)),
			},
		},
		{
			name: "Test_less_than_less_than_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_less_than_less_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 2)),
			},
		},
		{
			name: "Test_less_than_less_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 3)),
			},
			wantErr: true,
		},
		{
			name: "Test_less_than_pessimistic",
			set: verlib.Constraints{
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 0)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 1, 0)),
			},
		},
		{
			name: "Test_less_than_pessimistic_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("<", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 3, 0)),
			},
			wantErr: true,
		},
		// Tests for "<=" operator
		{
			name: "Test_less_than_equals_less_than_equals",
			set: verlib.Constraints{
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 3)),
			},
		},
		{
			name: "Test_less_than_equals_less_than_equals_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 0)),
			},
			wantErr: true,
		},
		{
			name: "Test_less_than_equals_pessimistic_sad",
			set: verlib.Constraints{
				verlib.NewConstraint("<=", verlib.NewVersion(1, 2, 3)),
				verlib.NewConstraint("~>", verlib.NewVersion(1, 3, 0)),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.set.Contradicts(tt.extraSet); (err != nil) != tt.wantErr {
				t.Errorf("Constraints.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsContradictionErr(t *testing.T) {
	testCases := []struct {
		name         string
		input        verlib.Constraints
		expectError  bool
		expectErrMsg string
	}{
		{
			name:         "contradictory constraints",
			input:        verlib.MustParseConstraintSet(">= 2.0.0, < 2.0.0"),
			expectError:  true,
			expectErrMsg: "constraints '>= 2.0.0' and '< 2.0.0' are contradictory",
		},
		{
			name:        "non-contradictory constraints",
			input:       verlib.MustParseConstraintSet(">= 2.0.0, <= 3.0.0"),
			expectError: false,
		},
		{
			name:         "contradictory constraints in separate sets",
			input:        verlib.MustParseConstraintSet(">= 2.0.0, < 2.0.0"),
			expectError:  true,
			expectErrMsg: "constraints '>= 2.0.0' and '< 2.0.0' are contradictory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Contradicts()
			if (err != nil) != tc.expectError {
				t.Fatalf("expected error: %v, got: %v", tc.expectError, err)
			}

			if err != nil {
				var ce verlib.ContradictionErr
				if errors.As(err, &ce) {
					if ce.Error() != tc.expectErrMsg {
						t.Errorf("expected error message: %s, got: %s", tc.expectErrMsg, ce.Error())
					}

					inputConstraint1, inputConstraint2 := tc.input[0], tc.input[1]
					errConstraint1, errConstraint2 := ce.Constraints()
					if !reflect.DeepEqual(inputConstraint1, errConstraint1) {
						t.Errorf("expected error constraint: %s, got: %s", inputConstraint1, errConstraint1)
					}
					if !reflect.DeepEqual(inputConstraint2, errConstraint2) {
						t.Errorf("expected error constraint: %s, got: %s", inputConstraint2, errConstraint2)
					}
				} else {
					t.Errorf("error is of incorrect type: got %T, want %T", err, ce)
				}
			}
		})
	}
}
