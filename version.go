package verlib

import (
	"fmt"
	"strconv"
)

// Version represents a version number compliant with Semantic Versioning (SemVer).
// It encapsulates the major, minor and patch version numbers, along with pre-release and build metadata information.
type Version struct {
	major         uint64  // major component of the version
	minor         *uint64 // minor component of the version. Optional.
	patch         *uint64 // patch component of the version. Optional.
	preRelease    string  // preRelease is the pre-release label of the version. Optional.
	buildMetadata string  // Build metadata. Optional
}

// NewVersion creates and returns a Version instance with the given major, minor, and patch numbers.
// Pre-release and build metadata fields are left empty.
func NewVersion(major, minor, patch uint64) Version {
	return Version{
		major: major,
		minor: &minor,
		patch: &patch,
	}
}

// NewPreReleaseVersion creates and returns a Version instance with the given major, minor, and patch numbers and a
// pre-release label. The build metadata field is left empty.
func NewPreReleaseVersion(major, minor, patch uint64, preRelease string) Version {
	return Version{
		major:      major,
		minor:      &minor,
		patch:      &patch,
		preRelease: preRelease,
	}
}

// Major returns the major component of the version.
func (v Version) Major() uint64 {
	return v.major
}

// Minor returns the minor component of the version. If minor is not defined, it returns 0.
func (v Version) Minor() uint64 {
	if v.minor == nil {
		return 0
	}
	return *v.minor
}

// Patch returns the patch component of the version. If patch is not defined, it returns 0.
func (v Version) Patch() uint64 {
	if v.patch == nil {
		return 0
	}
	return *v.patch
}

// PreRelease returns the pre-release label of the version. If no pre-release label exists, it returns an empty string.
func (v Version) PreRelease() string {
	return v.preRelease
}

// SetPreRelease creates a copy of the Version struct, sets the preRelease field to the given string,
// and returns the updated copy. This method can be used to change the pre-release component of a
// version without modifying the original Version struct.
//
// preRelease: The pre-release component of a version to set, according to the Semantic Versioning specification.
// In semantic versioning, the pre-release component is denoted by a hyphen followed by one or more dot-separated
// identifiers immediately following the patch version. Identifiers are comprised of alphanumeric characters and hyphens,
// and must not be empty, though they may contain zeros.
//
// Returns a new Version with the updated preRelease field.
func (v Version) SetPreRelease(preRelease string) Version {
	ver := v.clone()
	ver.preRelease = preRelease
	return ver
}

// BuildMetadata returns the build metadata associated with the version. If no build metadata exists, it returns an
// empty string.
func (v Version) BuildMetadata() string {
	return v.buildMetadata
}

// SetBuildMetadata creates a copy of the Version struct, sets the buildMetadata field to the given string,
// and returns the updated copy. This method can be used to change the build metadata of a version without
// modifying the original Version struct.
//
// buildMetadata: The build metadata to set, according to the Semantic Versioning specification.
// In semantic versioning, the build metadata component is denoted by a plus sign followed by one or
// more dot-separated identifiers immediately following the patch version or a pre-release.
// Identifiers are comprised of alphanumeric characters and hyphens.
//
// Returns a new Version with the updated buildMetadata field.
func (v Version) SetBuildMetadata(buildMetadata string) Version {
	ver := v.clone()
	ver.buildMetadata = buildMetadata
	return ver
}

// String returns a string representation of the Version, adhering to the SemVer 2.0.0 format.
func (v Version) String() string {
	versionStr := strconv.FormatUint(v.major, 10)

	if v.minor != nil {
		versionStr += "." + strconv.FormatUint(*v.minor, 10)
		if v.patch != nil {
			versionStr += "." + strconv.FormatUint(*v.patch, 10)
		}
	}
	if v.preRelease != "" {
		versionStr += "-" + v.preRelease
	}
	if v.buildMetadata != "" {
		versionStr += "+" + v.buildMetadata
	}
	return versionStr
}

// StrictString returns a string representation of the Version value.
// The resulting string strictly follows the Semantic Versioning 2.0.0 specification.
// The version string includes the major, minor and patch version numbers.
// If minor or patch versions are not provided, they default to "0".
// If a pre-release label is present, it is appended after the patch version,
// preceded by a hyphen "-".
// If build metadata is present, it is appended at the end,
// preceded by a plus sign "+".
//
// If the pre-release label or build metadata are not valid according to
// the Semantic Versioning 2.0.0 specification, the function returns an error.
//
// This method ensures that the resulting version string is a valid semantic version.
//
// Returns a string representation of the Version and an error which is
// non-nil in case of parsing or validation problems.
func (v Version) StrictString() (string, error) {
	versionStr := strconv.FormatUint(v.major, 10)

	if v.minor != nil {
		versionStr += "." + strconv.FormatUint(*v.minor, 10)
	} else {
		versionStr += ".0"
	}

	if v.patch != nil {
		versionStr += "." + strconv.FormatUint(*v.patch, 10)
	} else {
		versionStr += ".0"
	}

	if v.preRelease != "" {
		versionStr += "-" + v.preRelease
		if _, err := ParseSemVer(versionStr); err != nil {
			return "", fmt.Errorf("pre-release label invalid for strict version: %w", err)
		}
	}

	if v.buildMetadata != "" {
		versionStr += "+" + v.buildMetadata
		if _, err := ParseSemVer(versionStr); err != nil {
			return "", fmt.Errorf("build metadata invalid for strict version: %w", err)
		}
	}

	return versionStr, nil
}

// Less checks if this Version is less than the other Version. It compares the version fields
// from left to right (major, minor, patch), with precedence given to the leftmost non-equal field.
func (v Version) Less(other Version) bool {
	if v.Major() != other.Major() {
		return v.Major() < other.Major()
	}
	if v.Minor() != other.Minor() {
		return v.Minor() < other.Minor()
	}
	if v.Patch() != other.Patch() {
		return v.Patch() < other.Patch()
	}

	if v.preRelease == "" && other.preRelease != "" {
		return false
	}
	if v.preRelease != "" && other.preRelease == "" {
		return true
	}
	return v.preRelease < other.preRelease
}

// Greater checks if this Version is greater than the other Version. It makes use of the Less method defined for the
// Version type.
func (v Version) Greater(other Version) bool {
	return other.Less(v) && !v.Less(other)
}

// Equal checks if this Version is semantically equivalent to the other Version. This comparison disregards any
// build metadata.
func (v Version) Equal(other Version) bool {
	return !v.Less(other) && !other.Less(v)
}

// GreaterEqual checks if this Version is greater than or equal to the other Version. It's a logical OR operation on the
// results of Equal and Greater methods.
func (v Version) GreaterEqual(other Version) bool {
	return !v.Less(other)
}

// LessEqual checks if this Version is less than or equal to the other Version. It's a logical OR operation on the
// results of Equal and Less methods.
func (v Version) LessEqual(other Version) bool {
	return !v.Greater(other)
}

// Increment generates a new Version which represents the next possible version after the current one.
// It increments the rightmost version field that exists (patch, minor, or major),
// and unsets preRelease and buildMetadata.
func (v Version) Increment() Version {
	newVersion := v.clone()
	newVersion.preRelease, newVersion.buildMetadata = "", ""

	if v.patch != nil {
		newVersion.patch = new(uint64)
		*newVersion.patch = *v.patch + 1
		return newVersion
	}

	if v.minor != nil {
		newVersion.minor = new(uint64)
		*newVersion.minor = *v.minor + 1
		return newVersion
	}

	newVersion.major++
	return newVersion
}

// IncrementMajor generates a new Version by incrementing the major field of the current version,
// and sets minor and patch to zero. It unsets preRelease and buildMetadata.
func (v Version) IncrementMajor() Version {
	return Version{
		major: v.major + 1,
		minor: new(uint64),
		patch: new(uint64),
	}
}

// IncrementMinor generates a new Version by incrementing the minor field of the current version, and setting
// patch to zero. If minor does not exist, it sets it to 1. It unsets preRelease and buildMetadata.
func (v Version) IncrementMinor() Version {
	newVersion := v.clone()
	newVersion.preRelease = ""
	newVersion.buildMetadata = ""
	newVersion.patch = new(uint64)

	if newVersion.minor == nil {
		newVersion.minor = new(uint64)
	}

	*newVersion.minor++
	return newVersion
}

// IncrementPatch generates a new Version by incrementing the patch field of the current version.
// If patch does not exist, it sets it to 1. It unsets preRelease and buildMetadata.
func (v Version) IncrementPatch() Version {
	newVersion := v.clone()
	newVersion.preRelease, newVersion.buildMetadata = "", ""

	if newVersion.minor == nil {
		newVersion.minor = new(uint64)
	}
	if newVersion.patch == nil {
		newVersion.patch = new(uint64)
	}

	*newVersion.patch++
	return newVersion
}

// IncrementPessimistic creates a new Version instance that represents the
// next version number in a pessimistic versioning strategy. Pessimistic
// versioning, also known as "strict" versioning, assumes that any change in
// the software may break compatibility and thus increments the version
// components more aggressively than other strategies.
//
// The method proceeds as follows:
//
// 1. If the Version has a pre-release label or build metadata, they are removed.
//
//  2. If the Version has a patch version, it resets the patch version to 0 and
//     increments the minor version.
//
//  3. If the Version does not have a patch version but does have a minor version,
//     it resets the minor version to 0 and increments the major version.
//
// 4. If the Version does not have a minor version, it increments the major version.
//
// The function does not modify the original Version instance; instead, it
// operates on a clone to ensure that Version instances are immutable.
//
// This method is useful when you want to increment the version of your software
// in a way that communicates the potential for breaking changes, even if those
// changes are minor or patch-level changes.
//
// Usage:
//
//	v := verlib.MustParseVersion("1.2.3-alpha")
//	newV := v.IncrementPessimistic()  // newV is "1.3.0"
func (v Version) IncrementPessimistic() Version {
	newVersion := v.clone()
	newVersion.preRelease = ""
	newVersion.buildMetadata = ""

	if newVersion.patch != nil {
		*newVersion.patch = 0
		*newVersion.minor++
	} else if newVersion.minor != nil {
		newVersion.minor = new(uint64)
		newVersion.major++
	} else {
		newVersion.major++
	}
	return newVersion
}

// clone returns a deep copy of the Version. This is used in the Increment method to avoid mutating the original Version.
func (v Version) clone() Version {
	var minor, patch *uint64

	if v.minor != nil {
		minorVal := *v.minor
		minor = &minorVal
	}
	if v.patch != nil {
		patchVal := *v.patch
		patch = &patchVal
	}

	return Version{
		major:         v.major,
		minor:         minor,
		patch:         patch,
		preRelease:    v.preRelease,
		buildMetadata: v.buildMetadata,
	}
}
