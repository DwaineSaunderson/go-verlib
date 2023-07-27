package verlib

import (
	"errors"
)

// ContradictionErr is an error type that represents a contradiction
// between two Constraint instances. It embeds two Constraint structs: c1 and c2.
type ContradictionErr struct {
	c1 Constraint
	c2 Constraint
}

// Error implements the error interface for the ContradictionErr type.
// It returns a string indicating that the two constraints are contradictory.
func (ce ContradictionErr) Error() string {
	return "constraints '" + ce.c1.String() + "' and '" + ce.c2.String() + "' are contradictory"
}

// Constraints returns the two Constraint instances which are contradictory and caused the error.
func (ce ContradictionErr) Constraints() (Constraint, Constraint) {
	return ce.c1, ce.c2
}

// checkContradict checks if two Constraint structures contradict each other.
// This happens when no version can satisfy both constraints at the same time.
func checkContradict(c1, c2 Constraint) bool {
	switch {
	case c1.operator == c2.operator && c1.version.Equal(c2.version):
		return false
	case c1.operator == ">" && c2.operator == ">=" && c1.version.LessEqual(c2.version):
		return false
	case c1.operator == ">=" && c2.operator == ">" && c1.version.Greater(c2.version):
		return false
	case c1.operator == "!=" && c2.operator == "!=":
		return c1.version.Equal(c2.version)
	case c1.operator == "!=":
		return c1.version.Satisfies(c2)
	case c2.operator == "!=":
		return c2.version.Satisfies(c1)
	default:
		return !c1.Overlaps(c2)
	}
}

// Contradicts checks if any constraints in the Constraints and additional contradict each other.
// If contradiction exists, an error containing all contradictory pairs is returned.
func (c Constraints) Contradicts(additional ...Constraints) error {
	allConstraints := append(make(Constraints, 0, len(c)), c...)

	for _, additionalSet := range additional {
		allConstraints = append(allConstraints, additionalSet...)
	}

	var err error
	for i := 0; i < len(allConstraints); i++ {
		for j := i + 1; j < len(allConstraints); j++ {
			c1, c2 := allConstraints[i], allConstraints[j]
			if checkContradict(c1, c2) || checkContradict(c2, c1) {
				err = errors.Join(err, ContradictionErr{c1, c2})
			}
		}
	}

	return err
}
