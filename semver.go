package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version with private fields.
type Version struct {
	major      int
	minor      int
	patch      int
	preRelease string
	build      string
}

// New creates a new Version struct.
func New(major, minor, patch int, preRelease, build string) (*Version, error) {
	// Check for negative numbers
	if major < 0 || minor < 0 || patch < 0 {
		return nil, fmt.Errorf("version numbers cannot be negative")
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
func Parse(version string) (*Version, error) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([\w\.-]+))?(?:\+([\w\.-]+))?$`)
	matches := re.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid version format")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	// Check for negative numbers
	if major < 0 || minor < 0 || patch < 0 {
		return nil, fmt.Errorf("version numbers cannot be negative")
	}

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		preRelease: matches[4],
		build:      matches[5],
	}, nil
}

// String returns the string representation of the version.
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

// Compare compares two versions. Returns -1 if v < other, 1 if v > other, 0 if equal.
func (v *Version) Compare(other *Version) int {
	if v.major != other.major {
		return compareInts(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInts(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInts(v.patch, other.patch)
	}
	return comparePreRelease(v.preRelease, other.preRelease)
}

func compareInts(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func comparePreRelease(a, b string) int {
	if a == "" && b != "" {
		return 1
	} else if a != "" && b == "" {
		return -1
	} else if a == "" && b == "" {
		return 0
	}

	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aNum, aErr := strconv.Atoi(aParts[i])
		bNum, bErr := strconv.Atoi(bParts[i])

		if aErr == nil && bErr == nil {
			if aNum != bNum {
				return compareInts(aNum, bNum)
			}
		} else if aErr == nil {
			return -1
		} else if bErr == nil {
			return 1
		} else {
			if aParts[i] != bParts[i] {
				return strings.Compare(aParts[i], bParts[i])
			}
		}
	}

	return compareInts(len(aParts), len(bParts))
}

// IncrementMajor increments the major version and resets minor and patch.
func (v *Version) IncrementMajor() {
	v.major++
	v.minor = 0
	v.patch = 0
	v.preRelease = ""
	v.build = ""
}

// IncrementMinor increments the minor version and resets patch.
func (v *Version) IncrementMinor() {
	v.minor++
	v.patch = 0
	v.preRelease = ""
	v.build = ""
}

// IncrementPatch increments the patch version.
func (v *Version) IncrementPatch() {
	v.patch++
	v.preRelease = ""
	v.build = ""
}

// Major returns the major version.
func (v *Version) Major() int { return v.major }

// Minor returns the minor version.
func (v *Version) Minor() int { return v.minor }

// Patch returns the patch version.
func (v *Version) Patch() int { return v.patch }

// PreRelease returns the pre-release version.
func (v *Version) PreRelease() string { return v.preRelease }

// Build returns the build version.
func (v *Version) Build() string { return v.build }

// SetMajor sets the major version.
func (v *Version) SetMajor(major int) { v.major = major }

// SetMinor sets the minor version.
func (v *Version) SetMinor(minor int) { v.minor = minor }

// SetPatch sets the patch version.
func (v *Version) SetPatch(patch int) { v.patch = patch }

// SetPreRelease sets the pre-release version.
func (v *Version) SetPreRelease(preRelease string) { v.preRelease = preRelease }

// SetBuild sets the build version.
func (v *Version) SetBuild(build string) { v.build = build }
