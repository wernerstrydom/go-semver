package semver

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected *Version
		hasError bool
	}{
		// Valid cases
		{"1.0.0-alpha", &Version{1, 0, 0, "alpha", ""}, false},
		{"1.0.0-alpha.1", &Version{1, 0, 0, "alpha.1", ""}, false},
		{"1.0.0-alpha.beta", &Version{1, 0, 0, "alpha.beta", ""}, false},
		{"1.0.0-beta", &Version{1, 0, 0, "beta", ""}, false},
		{"1.0.0-beta.2", &Version{1, 0, 0, "beta.2", ""}, false},
		{"1.0.0-beta.11", &Version{1, 0, 0, "beta.11", ""}, false},
		{"1.0.0-rc.1", &Version{1, 0, 0, "rc.1", ""}, false},
		{"1.0.0", &Version{1, 0, 0, "", ""}, false},
		{"1.0.0-alpha+001", &Version{1, 0, 0, "alpha", "001"}, false},
		{"1.0.0+20130313144700", &Version{1, 0, 0, "", "20130313144700"}, false},
		{"1.0.0-beta+exp.sha.5114f85", &Version{1, 0, 0, "beta", "exp.sha.5114f85"}, false},
		{"1.0.0+21AF26D3----117B344092BD", &Version{1, 0, 0, "", "21AF26D3----117B344092BD"}, false},
		{"1.0.0-0A", &Version{1, 0, 0, "0A", ""}, false},

		// Edge cases
		{"0.0.0", &Version{0, 0, 0, "", ""}, false},                                                 // Minimum version
		{"999999999.999999999.999999999", &Version{999999999, 999999999, 999999999, "", ""}, false}, // Large numbers
		{"1.0.0-01", nil, true},       // Leading zero in pre-release
		{"1.0.0-.", nil, true},        // Invalid pre-release format
		{"1.0.0-..", nil, true},       // Invalid pre-release format
		{"1.0.0-123..456", nil, true}, // Invalid pre-release format
		{"1.0.0-123.0456", nil, true}, // Invalid pre-release format
		{"1.0.0+!@#$%^&*", nil, true}, // Invalid build metadata
		{"1.0.0-", nil, true},         // Trailing hyphen
		{"1.0.0+", nil, true},         // Trailing plus
		{"1.0.0-rc.1+build.1", &Version{1, 0, 0, "rc.1", "build.1"}, false},         // Valid pre-release and build
		{"1.0.0-rc.1+build.1.2.3", &Version{1, 0, 0, "rc.1", "build.1.2.3"}, false}, // Valid pre-release and build with multiple identifiers

		// Invalid cases
		{"invalid.version", nil, true},
		{"1.0", nil, true},     // Missing patch version
		{"1.0.0.0", nil, true}, // Extra version number
		{"1..0", nil, true},    // Missing minor version
		{"1.0.", nil, true},    // Trailing dot
	}

	for _, test := range tests {
		result, err := Parse(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("Parse(%q) error = %v, wantErr %v", test.input, err, test.hasError)
			continue
		}
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Parse(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		version  Version
		expected string
	}{
		// Basic version
		{Version{1, 0, 0, "", ""}, "1.0.0"},
		{Version{2, 1, 3, "", ""}, "2.1.3"},

		// With pre-release
		{Version{1, 0, 0, "alpha", ""}, "1.0.0-alpha"},
		{Version{1, 0, 0, "alpha.1", ""}, "1.0.0-alpha.1"},
		{Version{1, 0, 0, "beta", ""}, "1.0.0-beta"},
		{Version{1, 0, 0, "rc.1", ""}, "1.0.0-rc.1"},

		// With build metadata
		{Version{1, 0, 0, "", "001"}, "1.0.0+001"},
		{Version{1, 0, 0, "", "exp.sha.5114f85"}, "1.0.0+exp.sha.5114f85"},

		// With pre-release and build metadata
		{Version{1, 0, 0, "alpha", "001"}, "1.0.0-alpha+001"},
		{Version{1, 0, 0, "beta", "exp.sha.5114f85"}, "1.0.0-beta+exp.sha.5114f85"},
		{Version{1, 0, 0, "rc.1", "build.1.2.3"}, "1.0.0-rc.1+build.1.2.3"},
	}

	for _, test := range tests {
		result := test.version.String()
		if result != test.expected {
			t.Errorf("Version.String() = %v, want %v", result, test.expected)
		}
	}
}

func TestVersion_CompareTo(t *testing.T) {
	tests := []struct {
		this     Version
		other    Version
		expected int
	}{
		// Major version
		{Version{1, 0, 0, "", ""}, Version{2, 0, 0, "", ""}, -1},
		{Version{2, 0, 0, "", ""}, Version{1, 0, 0, "", ""}, 1},
		{Version{1, 0, 0, "", ""}, Version{1, 0, 0, "", ""}, 0},

		// Minor version
		{Version{1, 0, 0, "", ""}, Version{1, 1, 0, "", ""}, -1},
		{Version{1, 1, 0, "", ""}, Version{1, 0, 0, "", ""}, 1},
		{Version{1, 0, 0, "", ""}, Version{1, 0, 0, "", ""}, 0},

		// Patch version
		{Version{1, 0, 0, "", ""}, Version{1, 0, 1, "", ""}, -1},
		{Version{1, 0, 1, "", ""}, Version{1, 0, 0, "", ""}, 1},
		{Version{1, 0, 0, "", ""}, Version{1, 0, 0, "", ""}, 0},

		// Pre-release
		{Version{1, 0, 0, "alpha", ""}, Version{1, 0, 0, "beta", ""}, -1},
		{Version{1, 0, 0, "beta", ""}, Version{1, 0, 0, "alpha", ""}, 1},
		{Version{1, 0, 0, "alpha", ""}, Version{1, 0, 0, "alpha", ""}, 0},
		{Version{1, 0, 0, "alpha", ""}, Version{1, 0, 0, "alpha.1", ""}, -1},
		{Version{1, 0, 0, "alpha.1", ""}, Version{1, 0, 0, "alpha", ""}, 1},
		{Version{1, 0, 0, "alpha.1", ""}, Version{1, 0, 0, "alpha.1", ""}, 0},
		{Version{1, 0, 0, "alpha.1", ""}, Version{1, 0, 0, "alpha.2", ""}, -1},
		{Version{1, 0, 0, "alpha.2", ""}, Version{1, 0, 0, "alpha.1", ""}, 1},

		// Build metadata (ignored)
		{Version{1, 0, 0, "", "001"}, Version{1, 0, 0, "", "002"}, 0},
	}

	for _, test := range tests {
		result := test.this.CompareTo(&test.other)
		if result != test.expected {
			t.Errorf("Version.CompareTo(\"%s\", \"%s\") = %v, want %v", test.this.String(), test.other.String(), result, test.expected)
		}
	}
}

func TestVersion_IncreaseMajor(t *testing.T) {
	tests := []struct {
		version  Version
		expected Version
	}{
		{Version{1, 0, 0, "", ""}, Version{2, 0, 0, "", ""}},
		{Version{2, 1, 3, "", ""}, Version{3, 0, 0, "", ""}},
	}

	for _, test := range tests {
		test.version.IncreaseMajor()
		if !reflect.DeepEqual(test.version, test.expected) {
			t.Errorf("Version.IncreaseMajor() = %v, want %v", test.version, test.expected)
		}
	}
}

func TestVersion_IncreaseMinor(t *testing.T) {
	tests := []struct {
		version  Version
		expected Version
	}{
		{Version{1, 0, 0, "", ""}, Version{1, 1, 0, "", ""}},
		{Version{2, 1, 3, "", ""}, Version{2, 2, 0, "", ""}},
	}

	for _, test := range tests {
		test.version.IncreaseMinor()
		if !reflect.DeepEqual(test.version, test.expected) {
			t.Errorf("Version.IncreaseMinor() = %v, want %v", test.version, test.expected)
		}
	}
}

func TestVersion_IncreasePatch(t *testing.T) {
	tests := []struct {
		version  Version
		expected Version
	}{
		{Version{1, 0, 0, "", ""}, Version{1, 0, 1, "", ""}},
		{Version{2, 1, 3, "", ""}, Version{2, 1, 4, "", ""}},
	}

	for _, test := range tests {
		test.version.IncreasePatch()
		if !reflect.DeepEqual(test.version, test.expected) {
			t.Errorf("Version.IncreasePatch() = %v, want %v", test.version, test.expected)
		}
	}
}

func TestVersion_IncreasePreRelease(t *testing.T) {
	tests := []struct {
		version  Version
		expected Version
	}{
		{Version{1, 0, 0, "alpha", ""}, Version{1, 0, 0, "alpha.1", ""}},
		{Version{1, 0, 0, "alpha.1", ""}, Version{1, 0, 0, "alpha.2", ""}},
	}

	for _, test := range tests {
		test.version.IncreasePreRelease()
		if !reflect.DeepEqual(test.version, test.expected) {
			t.Errorf("Version.IncreasePreRelease() = %v, want %v", test.version, test.expected)
		}
	}
}

func TestVersion_IsStable(t *testing.T) {
	tests := []struct {
		version  Version
		expected bool
	}{
		{Version{1, 0, 0, "", ""}, true},
		{Version{1, 0, 0, "alpha", ""}, false},
	}

	for _, test := range tests {
		result := test.version.IsStable()
		if result != test.expected {
			t.Errorf("Version.IsStable() = %v, want %v", result, test.expected)
		}
	}
}

func TestVersion_IsPreRelease(t *testing.T) {
	tests := []struct {
		version  Version
		expected bool
	}{
		{Version{1, 0, 0, "", ""}, false},
		{Version{1, 0, 0, "alpha", ""}, true},
	}

	for _, test := range tests {
		result := test.version.IsPreRelease()
		if result != test.expected {
			t.Errorf("Version.IsPreRelease() = %v, want %v", result, test.expected)
		}
	}
}

func TestVersion_Major(t *testing.T) {
	tests := []struct {
		version  Version
		expected int
	}{
		{Version{1, 0, 0, "", ""}, 1},
		{Version{2, 1, 3, "", ""}, 2},
	}

	for _, test := range tests {
		result := test.version.Major()
		if result != test.expected {
			t.Errorf("Version.Major() = %v, want %v", result, test.expected)
		}
	}
}

func TestVersion_Minor(t *testing.T) {
	tests := []struct {
		version  Version
		expected int
	}{
		{Version{1, 0, 0, "", ""}, 0},
		{Version{2, 1, 3, "", ""}, 1},
	}

	for _, test := range tests {
		result := test.version.Minor()
		if result != test.expected {
			t.Errorf("Version.Minor() = %v, want %v", result, test.expected)
		}
	}
}

func TestVersion_Patch(t *testing.T) {
	tests := []struct {
		version  Version
		expected int
	}{
		{Version{1, 0, 0, "", ""}, 0},
		{Version{2, 1, 3, "", ""}, 3},
	}

	for _, test := range tests {
		result := test.version.Patch()
		if result != test.expected {
			t.Errorf("Version.Patch() = %v, want %v", result, test.expected)
		}
	}
}

func Test_compareIdentifiers(t *testing.T) {
	tests := []struct {
		this     string
		other    string
		expected int
	}{
		{"alpha", "beta", -1},
		{"beta", "alpha", 1},
		{"alpha", "alpha", 0},
		{"alpha", "alpha.1", -1},
		{"alpha.1", "alpha", 1},
		{"alpha.1", "alpha.1", 0},
		{"alpha.1", "alpha.2", -1},
		{"alpha.2", "alpha.1", 1},
	}

	for _, test := range tests {
		result := compareIdentifiers(test.this, test.other)
		if result != test.expected {
			t.Errorf("compareIdentifiers(\"%s\", \"%s\") = %v, want %v", test.this, test.other, result, test.expected)
		}
	}
}

func Test_isNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"0", true},
		{"1", true},
		{"123", true},
		{"01", true},
		{"1.0", false},
		{"1.0.0", false},
		{"a", false},
		{"alpha", false},
		{"alpha.1", false},
	}

	for _, test := range tests {
		result := isNumeric(test.input)
		if result != test.expected {
			t.Errorf("isNumeric(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}
