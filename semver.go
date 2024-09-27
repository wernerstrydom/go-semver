package semver

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version according to Semantic Versioning 2.0.0.
type Version struct {
	major      int    // Major version
	minor      int    // Minor version
	patch      int    // Patch version
	preRelease string // Pre-release version (optional)
	build      string // Build metadata (optional)
}

var (
	ErrorInvalidMajorVersion   = errors.New("major version must be numeric, and non-negative")
	ErrorInvalidMinorVersion   = errors.New("minor version must be numeric, and non-negative")
	ErrorInvalidPatchVersion   = errors.New("patch version must be numeric, and non-negative")
	ErrorInvalidPreRelease     = errors.New("pre-release identifiers must not be empty")
	ErrorPreReleaseLeadingZero = errors.New("numeric pre-release identifiers must not have leading zeros")
	ErrorInvalidVersion        = errors.New("invalid version format: must be in the form X.Y.Z[-PRERELEASE][+BUILD]")
)

// New creates a new Version struct.
func New(major, minor, patch int, preRelease, build string) (*Version, error) {
	if major < 0 {
		return nil, ErrorInvalidMajorVersion
	}
	if minor < 0 {
		return nil, ErrorInvalidMinorVersion
	}
	if patch < 0 {
		return nil, ErrorInvalidPatchVersion
	}
	if preRelease != "" {
		identifiers := strings.Split(preRelease, ".")
		for _, id := range identifiers {
			if id == "" {
				return nil, ErrorInvalidPreRelease
			}
			if isNumeric(id) && hasLeadingZero(id) {
				return nil, ErrorPreReleaseLeadingZero
			}
		}
	}

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		preRelease: preRelease,
		build:      build,
	}, nil
}

// Parse parses a version string and returns a Version struct.
func Parse(versionStr string) (*Version, error) {
	// Regular expression to match semantic versioning
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([\da-zA-Z-]+(?:\.[\da-zA-Z-]+)*))?(?:\+([\da-zA-Z-]+(?:\.[\da-zA-Z-]+)*))?$`)
	matches := re.FindStringSubmatch(versionStr)

	if matches == nil {
		return nil, ErrorInvalidVersion
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, ErrorInvalidMajorVersion
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, ErrorInvalidMinorVersion
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, ErrorInvalidPatchVersion
	}

	preRelease := matches[4]
	build := matches[5]

	return New(major, minor, patch, preRelease, build)
}

// isNumeric checks if a string is numeric.
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// hasLeadingZero checks if a numeric string has a leading zero.
func hasLeadingZero(s string) bool {
	return len(s) > 1 && s[0] == '0'
}

// String converts a Version struct to its string representation.
func (v *Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.preRelease != "" {
		version += "-" + v.preRelease
	}
	if v.build != "" {
		version += "+" + v.build
	}
	return version
}

// CompareTo compares two versions and returns -1, 0, or 1 if v is less than, equal to, or greater than w.
func (v *Version) CompareTo(other *Version) int {
	if v.major < other.major {
		return -1
	}
	if v.major > other.major {
		return 1
	}
	if v.minor < other.minor {
		return -1
	}
	if v.minor > other.minor {
		return 1
	}
	if v.patch < other.patch {
		return -1
	}
	if v.patch > other.patch {
		return 1
	}
	if v.preRelease == "" && other.preRelease == "" {
		return 0
	}
	if v.preRelease == "" {
		return 1
	}
	if other.preRelease == "" {
		return -1
	}
	return compareIdentifiers(v.preRelease, other.preRelease)
}

func compareIdentifiers(this string, other string) int {
	thisIdentifiers := strings.Split(this, ".")
	otherIdentifiers := strings.Split(other, ".")
	for i := 0; i < len(thisIdentifiers) && i < len(otherIdentifiers); i++ {
		thisID := thisIdentifiers[i]
		otherID := otherIdentifiers[i]
		if isNumeric(thisID) && isNumeric(otherID) {
			thisNum, _ := strconv.Atoi(thisID)
			otherNum, _ := strconv.Atoi(otherID)
			if thisNum < otherNum {
				return -1
			}
			if thisNum > otherNum {
				return 1
			}
		} else {
			if thisID < otherID {
				return -1
			}
			if thisID > otherID {
				return 1
			}
		}
	}
	if len(thisIdentifiers) < len(otherIdentifiers) {
		return -1
	}
	if len(thisIdentifiers) > len(otherIdentifiers) {
		return 1
	}
	return 0
}

// Major returns the major version.
func (v *Version) Major() int {
	return v.major
}

// Minor returns the minor version.
func (v *Version) Minor() int {
	return v.minor
}

// Patch returns the patch version.
func (v *Version) Patch() int {
	return v.patch
}

// PreRelease returns the pre-release version.
func (v *Version) PreRelease() string {
	return v.preRelease
}

// Build returns the build metadata.
func (v *Version) Build() string {
	return v.build
}

// IsStable returns true if the version is stable (i.e., no pre-release version).
func (v *Version) IsStable() bool {
	return v.preRelease == ""
}

// IsPreRelease returns true if the version is a pre-release version.
func (v *Version) IsPreRelease() bool {
	return v.preRelease != ""
}

// IncreaseMajor increments the major version by 1 and resets the minor and patch versions to 0.
func (v *Version) IncreaseMajor() {
	v.major++
	v.minor = 0
	v.patch = 0
}

// IncreaseMinor increments the minor version by 1 and resets the patch version to 0.
func (v *Version) IncreaseMinor() {
	v.minor++
	v.patch = 0
}

// IncreasePatch increments the patch version by 1.
func (v *Version) IncreasePatch() {
	v.patch++
}

// IncreasePreRelease increments the pre-release version by 1.
func (v *Version) IncreasePreRelease() {
	identifiers := strings.Split(v.preRelease, ".")
	if len(identifiers) == 0 {
		v.preRelease = "1"
		return
	}
	last := identifiers[len(identifiers)-1]
	if isNumeric(last) {
		identifiers[len(identifiers)-1] = strconv.Itoa(1 + toInt(last))
	} else {
		identifiers = append(identifiers, "1")
	}
	v.preRelease = strings.Join(identifiers, ".")
}

func toInt(last string) int {
	i, _ := strconv.Atoi(last)
	return i
}
