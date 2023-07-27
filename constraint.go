package verlib

import (
	"fmt"
	"strings"
)

// Operator represents a comparator between versions. It supports equality, inequality,
// greater than, greater or equal to, less than, less or equal to, and a pessimistic
// greater or equal comparator.
type Operator string

const (
	EQ            Operator = "="  // EQ stands for equality. It allows only one exact version number.
	NE            Operator = "!=" // NE stands for not equal. It excludes an exact version number.
	GT            Operator = ">"  // GT stands for greater than. It allows strictly newer versions.
	GE            Operator = ">=" // GE stands for greater than or equal. It allows newer versions and the exact number specified.
	LT            Operator = "<"  // LT stands for less than. It allows strictly older versions.
	LE            Operator = "<=" // LE stands for less than or equal. It allows older versions and the exact number specified.
	GEPessimistic Operator = "~>" // GEPessimistic stands for pessimistic greater than or equal. It allows only the rightmost version component to increment.
)

// String converts an Operator to its string representation.
func (o Operator) String() string {
	return string(o)
}

// Constraint represents a comparison between a version number and a value. It is used
// to determine whether a version number satisfies a specific condition.
type Constraint struct {
	operator Operator // operator specifies the type of constraint, such as "=", ">", "<", etc.
	version  Version  // version is the version number the constraint is compared to.
}

// Constraints represents a collection of Constraints, all of which must be
// satisfied for a version number to be considered acceptable.
type Constraints []Constraint

// NewConstraint creates a new constraint with the given operator and version.
func NewConstraint(operator Operator, version Version) Constraint {
	return Constraint{
		operator: operator,
		version:  version,
	}
}

// String returns a string representation of a Constraint.
func (c Constraint) String() string {
	return c.operator.String() + " " + c.version.String()
}

// StrictString returns a string representation of the Constraint value.
// The resulting string strictly follows the conventions of version constraints,
// combining the operator and the strict string representation of the version.
// The operator and version are separated by a space.
//
// The version string is generated using the Version type's StrictString method,
// which ensures it's a valid semantic version string according to Semantic Versioning 2.0.0.
//
// If the StrictString method of the version fails to generate a string,
// this method returns an error indicating that the constraint version string
// generation failed.
//
// Returns a string representation of the Constraint and an error which is
// non-nil in case of parsing or validation problems.
func (c Constraint) StrictString() (string, error) {
	versionString, err := c.version.StrictString()
	if err != nil {
		return "", fmt.Errorf("failed to generate strict string for constraint version: %w", err)
	}

	return c.operator.String() + " " + versionString, nil
}

// String returns a string representation of the Constraints slice. Each constraint in the slice
// is converted to a string using the Constraint's String method, and these are all joined together
// with commas. This function is particularly useful for displaying or logging the constraints in a human-readable format.
func (c Constraints) String() string {
	constraintStrings := make([]string, 0, len(c))
	for _, constraint := range c {
		constraintStrings = append(constraintStrings, constraint.String())
	}

	return strings.Join(constraintStrings, ", ")
}

// StrictString returns a strict string representation of the Constraints slice. Each constraint in the slice
// is converted to a strict string using the Constraint's StrictString method, and these are all joined together
// with commas. This function is useful when you need a precise, parseable string representation of the constraints.
func (c Constraints) StrictString() (string, error) {
	constraintStrings := make([]string, 0, len(c))
	for _, constraint := range c {
		constraintString, err := constraint.StrictString()
		if err != nil {
			return "", fmt.Errorf("failed to generate strict string for constraint: %w", err)
		}

		constraintStrings = append(constraintStrings, constraintString)
	}

	return strings.Join(constraintStrings, ", "), nil
}

// Satisfies determines whether a given Version v satisfies a Constraint c.
// It returns true if the version satisfies the constraint, false otherwise.
func (v Version) Satisfies(c Constraint) bool {
	switch c.operator {
	case EQ:
		return v.Equal(c.version)
	case NE:
		return !v.Equal(c.version)
	case GT:
		return v.Greater(c.version)
	case GE:
		return v.GreaterEqual(c.version)
	case LT:
		return v.Less(c.version)
	case LE:
		return v.LessEqual(c.version)
	case GEPessimistic:
		return v.GreaterEqual(c.version) && v.Less(c.version.IncrementPessimistic())
	default:
		return false
	}
}

// Overlaps determines whether two Constraints c1 and c2 overlap, meaning
// there is at least one version that would satisfy both constraints.
func (c Constraint) Overlaps(c2 Constraint) bool {
	switch {
	case c.operator == EQ:
		return c.version.Satisfies(c2)
	case c.operator == NE:
		return !c.version.Satisfies(c2)
	case c.operator == GT:
		return c.version.LessEqual(c2.version)
	case c.operator == GE && c2.operator == GEPessimistic:
		return c.version.GreaterEqual(c2.version) && c.version.Less(c2.version.Increment())
	case c2.operator == GE && c.operator == GEPessimistic:
		return c2.version.GreaterEqual(c.version) && c2.version.Less(c.version.Increment())
	case c.operator == GE:
		return c.version.Less(c2.version)
	case c.operator == LT && c2.operator == LE:
		return c2.version.Less(c.version)
	case c.operator == LE && c2.operator == LT:
		return c.version.Less(c2.version)
	case c.operator == LT && c2.operator == GEPessimistic:
		return c2.version.Less(c.version) || c2.version.Equal(c.version)
	case c2.operator == LT && c.operator == GEPessimistic:
		return c.version.Less(c2.version) || c.version.Equal(c2.version)
	case c.operator == LT:
		return c.version.GreaterEqual(c2.version)
	case c.operator == LE:
		return c.version.Greater(c2.version)
	case c.operator == GEPessimistic:
		return c2.version.GreaterEqual(c.version) && c2.version.Less(c.version.Increment())
	default:
		return false
	}
}
