package verlib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionRegex matches versions which are semantic-version like, in a less strict manner, capturing the major,
// minor, and patch versions, the pre-release label, and the build metadata.
//
// ^ : Matches the start of the line.
// \D*? : Matches zero or more non-digit characters (not greedily), e.g., "v" in "v1.0.0".
// This part allows version numbers to be preceded by non-digit characters.
//
// (?P<major>\d+) : Named capture group "major" that matches one or more digits, e.g., "1" in "v1.0.0".
// This represents the major version.
//
// (\.(?P<minor>\d+))? : Optionally matches a period followed by one or more digits in the "minor" group,
// e.g., ".0" in "v1.0.0". This represents the minor version.
//
// (\.(?P<patch>\d+))? : Optionally matches a period followed by one or more digits in the "patch" group,
// e.g., ".0" in "v1.0.0". This represents the patch version.
//
// (?:-(?P<prerelease>[^+\n]*))? : Optionally matches a dash followed by any characters except '+' or newline in the
// "prerelease" group, e.g., "-alpha" in "v1.0.0-alpha". This represents the pre-release label.
//
// (\+(?P<buildmetadata>.*))? : Optionally matches a plus sign followed by any characters in the "buildmetadata"
// group, e.g., "+20130313144700" in "v1.0.0+20130313144700". This represents the build metadata.
//
// [^\n]*$ : Matches any characters except newlines until the end of the line. This part is intended to discard any
// non-digit data after the version number.
var versionRegex = regexp.MustCompile(`^\D*?(?P<major>\d+)(\.(?P<minor>\d+))?(\.(?P<patch>\d+))?(?:-(?P<prerelease>[^+\n]*))?(\+(?P<buildmetadata>.*))?[^\n]*$`)

// semVerRegex adheres to the Semantic Versioning specification (SemVer 2.0.0). It enforces the absence of
// leading zeros in numeric identifiers, and allows for the optional prerelease and build metadata components of a
// version, where the identifiers must meet certain criteria.
//
// (?P<major>0|[1-9]\d*) : Named capture group "major" that matches a zero or a non-zero number without leading
// zeros. This represents the major version.
//
// \.(?P<minor>0|[1-9]\d*) : A period followed by the "minor" group that matches a zero or a non-zero number without
// leading zeros. This represents the minor version.
//
// \.(?P<patch>0|[1-9]\d*) : A period followed by the "patch" group that matches a zero or a non-zero number without
// leading zeros. This represents the patch version.
//
// (?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))? :
// Optionally matches a dash followed by the "prerelease" group. The prerelease group can be either a numeric
// identifier, or an alphanumeric identifier that can include hyphen (-) but must not be fully numerical.
// Pre-release versions can have multiple identifiers separated by periods.
//
// (?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))? : Optionally matches a plus sign followed by the
// "buildmetadata" group. The build metadata is a dot-separated series of alphanumeric identifiers, and can include
// the hyphen (-).
//
// $ : Matches the end of the line.
var semVerRegex = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// parseVersion takes a regular expression (r) and a version string as input.
// It attempts to parse the version string according to the structure of the regular expression,
// and returns a populated Version struct and nil on successful parsing.
//
// The function uses named capture groups in the regular expression to extract the 'major', 'minor',
// 'patch', 'prerelease', and 'buildmetadata' parts of the version string.
//
// If the version string does not match the regular expression, the function will return an empty
// Version struct and an error stating that the regex failed to match.
//
// If there are more submatches than subexpression names in the regular expression, the function will panic.
// This is a programming error that should not happen with correct usage.
//
// If any of the 'major', 'minor', or 'patch' parts are not valid unsigned integers, the function will
// return an empty Version struct and an error indicating the failure to parse that part.
//
// Parameters:
//   - r: The regular expression used to parse the version string. It should include named capture
//     groups for 'major', 'minor', 'patch', 'prerelease', and 'buildmetadata'.
//   - versionString: The version string to parse.
//
// Returns:
//   - A populated Version struct on successful parsing.
//   - An error if the version string does not match the regular expression, or if any of the 'major',
//     'minor', or 'patch' parts are not valid unsigned integers.
func parseVersion(r *regexp.Regexp, versionString string) (Version, error) {
	matches := r.FindStringSubmatch(versionString)
	if len(matches) == 0 {
		return Version{}, fmt.Errorf("regex failed to match version %q", versionString)
	}

	subExpNames := r.SubexpNames()
	if len(matches) > len(subExpNames) {
		panic("length of matches > length of sub-expression names")
	}

	ver := Version{}

	for i, subMatch := range matches {
		switch subExpNames[i] {
		case "major":
			major, err := strconv.ParseUint(subMatch, 10, 64)
			if err != nil {
				return Version{}, fmt.Errorf("failed to parse major version component %q: %w", subMatch, err)
			}
			ver.major = major
		case "minor":
			if subMatch == "" {
				continue
			}
			minor, err := strconv.ParseUint(subMatch, 10, 64)
			if err != nil {
				return Version{}, fmt.Errorf("failed to parse minor version component %q: %w", subMatch, err)
			}
			ver.minor = &minor
		case "patch":
			if subMatch == "" {
				continue
			}
			patch, err := strconv.ParseUint(subMatch, 10, 64)
			if err != nil {
				return Version{}, fmt.Errorf("failed to patch version component %q: %w", subMatch, err)
			}
			ver.patch = &patch
		case "prerelease":
			ver.preRelease = subMatch
		case "buildmetadata":
			ver.buildMetadata = subMatch
		}
	}

	return ver, nil
}

// ParseVersion takes a version string as input and attempts to parse it according to a
// pre-defined regular expression (versionRegex).
// Returns a populated Version struct and nil on successful parsing.
// If the parsing fails, it returns an empty Version struct and an error.
func ParseVersion(version string) (Version, error) {
	return parseVersion(versionRegex, version)
}

// MustParseVersion is similar to ParseVersion, but it panics if the parsing fails.
// It's useful when you're certain the input version string is valid, and any failure
// is a programming error that should stop the program execution.
// Returns a populated Version struct on successful parsing.
func MustParseVersion(version string) Version {
	v, err := ParseVersion(version)
	if err != nil {
		panic(fmt.Errorf("failed to parse version: %w", err))
	}
	return v
}

// ParseSemVer takes a semantic version string as input and attempts to parse it
// according to a pre-defined regular expression (semVerRegex) designed to parse
// Semantic Versioning 2.0.0 compliant versions.
// Returns a populated Version struct and nil on successful parsing.
// If the parsing fails, it returns an empty Version struct and an error.
func ParseSemVer(semanticVersion string) (Version, error) {
	return parseVersion(semVerRegex, semanticVersion)
}

// MustParseSemVer is similar to ParseSemVer, but it panics if the parsing fails.
// It's useful when you're certain the input semantic version string is valid, and any failure
// is a programming error that should stop the program execution.
// Returns a populated Version struct on successful parsing.
func MustParseSemVer(version string) Version {
	v, err := ParseSemVer(version)
	if err != nil {
		panic(fmt.Errorf("failed to parse semantic version: %w", err))
	}
	return v
}

// constraintRegex is a regular expression used to parse version constraints.
// The regular expression is broken down as follows:
//
// `^` : Matches the start of the line.
//
// `(!=|=|>=|>|<=|<|~>)?` : This group matches an optional comparison operator. The comparison operator can be one of
// the following: "!=" (not equal), "=" (equal), ">=" (greater than or equal to), ">" (greater than), "<=" (less than
// or equal to), "<" (less than), or "~>" (approximately greater than).
//
// `[^\d\n]*` : This group matches zero or more characters that are neither digits nor newlines. This part is intended
// to allow any optional text before the version number.
//
// `(\d+\S*)` : This group matches one or more digits followed by zero or more non-whitespace characters. This part
// is intended to capture a version number that may include minor and patch versions, along with any additional
// alphanumeric or symbolic data (like "-alpha", "+20130313144700" in semantic versioning).
//
// `$` : Matches the end of the line.
var constraintRegex = regexp.MustCompile(`^(!=|=|>=|>|<=|<|~>)?[^\d\n]*(\d+\S*)$`)

// strictConstraintRegex is a regular expression used to parse strict version constraints.
// The regular expression is broken down as follows:
//
// `^` : Matches the start of the line.
//
// `(P<operator>!=|=|>=|>|<=|<|~>)?` : This named capture group matches an optional comparison operator.
// The comparison operator can be one of the following: "!=" (not equal), "=" (equal), ">=" (greater than or equal to),
// ">" (greater than), "<=" (less than or equal to), "<" (less than), or "~>" (approximately greater than).
//
// `\s*` : Matches any whitespace character between the operator and the semver.
//
// `(?P<semver>` : Named capture group "semver" that matches a semver adhering to the Semantic Versioning specification (SemVer 2.0.0).
//
// `$` : Matches the end of the line.
var strictConstraintRegex = regexp.MustCompile(`^(?P<operator>!=|=|>=|>|<=|<|~>)?\s*(?P<semver>(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)(?:-(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(?:\+[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*)?)?$`)

// ParseConstraint takes a version constraint string as input and attempts to parse it
// according to a pre-defined regular expression (constraintRegex).
// The version constraint string should contain an operator and a version.
// Returns a populated Constraint struct and nil on successful parsing.
// If the parsing fails or if the operator is invalid, it returns an empty Constraint struct and an error.
func ParseConstraint(verConstraint string) (Constraint, error) {
	verConstraint = strings.TrimSpace(verConstraint)
	constraintMatches := constraintRegex.FindStringSubmatch(verConstraint)
	if len(constraintMatches) != 3 {
		return Constraint{}, fmt.Errorf("failed to parse version constraint %q", verConstraint)
	}

	parsedOperator, rawVersion := Operator(constraintMatches[1]), constraintMatches[2]
	if parsedOperator == "" {
		parsedOperator = EQ
	}

	switch parsedOperator {
	case EQ, NE, GT, GE, LT, LE, GEPessimistic:
	default:
		return Constraint{}, fmt.Errorf("invalid operator %q in constraint %q", parsedOperator, verConstraint)
	}

	version, err := ParseVersion(rawVersion)
	if err != nil {
		return Constraint{}, fmt.Errorf("failed to parse version in constraint: %w", err)
	}

	return Constraint{
		operator: parsedOperator,
		version:  version,
	}, nil
}

// MustParseConstraint is similar to ParseConstraint, but it panics if the parsing fails.
// It's useful when you're certain the input constraint string is valid, and any failure
// is a programming error that should stop the program execution.
// Returns a populated Constraint struct on successful parsing.
func MustParseConstraint(verConstraint string) Constraint {
	constraint, err := ParseConstraint(verConstraint)
	if err != nil {
		panic(err)
	}
	return constraint
}

// ParseStrictConstraint takes a strict version constraint string as input and attempts to parse it
// according to a pre-defined regular expression (strictConstraintRegex).
// The version constraint string should contain an operator (or no operator), any amount of whitespace and a strict semver v2.
// Returns a populated Constraint struct and nil on successful parsing.
// If the parsing fails or if the operator is invalid, it returns an empty Constraint struct and an error.
func ParseStrictConstraint(verConstraint string) (Constraint, error) {
	verConstraint = strings.TrimSpace(verConstraint)
	matches := strictConstraintRegex.FindStringSubmatch(verConstraint)

	subExpNames := strictConstraintRegex.SubexpNames()
	if len(matches) > len(subExpNames) {
		panic("length of matches > length of sub-expression names")
	}

	var (
		parsedOperator Operator
		rawVersion     string
	)

	for i, subMatch := range matches {
		if subExpNames[i] == "operator" {
			if subMatch == "" {
				subMatch = "="
			}
			parsedOperator = Operator(subMatch)
			continue
		}
		if subExpNames[i] == "semver" {
			rawVersion = subMatch
			continue
		}
	}

	if parsedOperator == "" || rawVersion == "" {
		return Constraint{}, fmt.Errorf("failed to parse strict version constraint %q", verConstraint)
	}

	switch parsedOperator {
	case EQ, NE, GT, GE, LT, LE, GEPessimistic:
	default:
		return Constraint{}, fmt.Errorf("invalid operator %q in strict constraint %q", parsedOperator, verConstraint)
	}

	version, err := ParseSemVer(rawVersion)
	if err != nil {
		return Constraint{}, fmt.Errorf("failed to parse semver in strict constraint: %w", err)
	}

	return Constraint{
		operator: parsedOperator,
		version:  version,
	}, nil
}

// ParseConstraintSet takes a constraint string, splits it by commas and attempts to parse each
// resulting constraint according to the ParseConstraint function.
// Returns a populated Constraints (slice of Constraint) and nil on successful parsing.
// If the parsing fails for any constraint, it returns an empty Constraints and an error.
func ParseConstraintSet(constraintString string) (Constraints, error) {
	var result Constraints

	versionConstraints := strings.Split(constraintString, ",")

	for _, rawConstraint := range versionConstraints {
		constraint, err := ParseConstraint(rawConstraint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse constraint %q: %w", rawConstraint, err)
		}

		result = append(result, constraint)
	}

	return result, nil
}

// MustParseConstraintSet is similar to ParseConstraintSet, but it panics if the parsing fails.
// It's useful when you're certain the input constraint set string is valid, and any failure
// is a programming error that should stop the program execution.
// Returns a populated Constraints (slice of Constraint) on successful parsing.
func MustParseConstraintSet(constraintString string) Constraints {
	constraintSet, err := ParseConstraintSet(constraintString)
	if err != nil {
		panic(err)
	}
	return constraintSet
}

// ParseStrictConstraintSet takes a constraint string, splits it by commas and attempts to parse each
// resulting constraint according to the ParseStrictConstraint function.
// Returns a populated Constraints (slice of Constraint) and nil on successful parsing.
// If the parsing fails for any constraint, it returns an empty Constraints and an error.
func ParseStrictConstraintSet(constraintString string) (Constraints, error) {
	var result Constraints

	versionConstraints := strings.Split(constraintString, ",")

	for _, rawConstraint := range versionConstraints {
		constraint, err := ParseStrictConstraint(rawConstraint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse strict constraint %q: %w", rawConstraint, err)
		}

		result = append(result, constraint)
	}

	return result, nil
}

// MustParseStrictConstraintSet is a utility function that wraps ParseStrictConstraintSet,
// but instead of returning an error, it panics if an error occurs.
// This is useful when initializing global variables or in other cases where you want to
// ensure that a constraint set is parsed without checking errors at each usage.
// If successful, returns a populated Constraints (slice of Constraint).
func MustParseStrictConstraintSet(constraintString string) Constraints {
	constraintSet, err := ParseStrictConstraintSet(constraintString)
	if err != nil {
		panic(err)
	}
	return constraintSet
}
