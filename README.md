# go-semver

## Introduction
The `go-semver` project is a Go-based library designed to handle semantic versioning, a versioning scheme that conveys 
meaning about the underlying changes with each release. Semantic versioning uses a three-part version number: 
Major.Minor.Patch, and is a standard approach for versioning in software development.

## Installation

To include `go-semver` in your Go project, use the following command to download and install the package:

```shell
go get github.com/wernerstrydom/go-semver
```

## Usage

After installation, you can integrate go-semver into your project to manage and interpret version numbers. This can be 
particularly useful for managing dependencies and ensuring compatibility across different stages of software development.

### Example 1: Creating a New Version

To create a new semantic version, use the `New` function. This function requires the major, minor, and patch numbers, and optionally, pre-release and build metadata.

```go
package main

import (
	"fmt"
	"log"
    "github.com/wernerstrydom/go-semver"
)

func main() {
	version, err := semver.New(1, 0, 0, "alpha", "001")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New Version:", version.String())
}
```

This code snippet initializes a new version `1.0.0-alpha+001`, and prints it. If any input is invalid, an error is logged.

### Example 2: Parsing a Version String

To parse a semantic version string, use the `Parse` function. This converts a version string into a `Version` struct.

```go
package main

import (
	"fmt"
	"log"
    "github.com/wernerstrydom/go-semver"
)

func main() {
	versionStr := "1.0.0-beta+exp.sha.5114f85"
	version, err := semver.Parse(versionStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Parsed Version:", version.String())
}
```

This example parses the version string `1.0.0-beta+exp.sha.5114f85`. If parsing fails, an error is logged.

### Example 3: Comparing Versions

To compare two versions, use the `CompareTo` method. It returns -1 if the first version is less, 0 if equal, and 1 if greater.

```go
package main

import (
	"fmt"
	"log"
    "github.com/wernerstrydom/go-semver"
)

func main() {
	version1, err := semver.New(1, 0, 0, "", "")
	if err != nil {
		log.Fatal(err)
	}
	version2, err := semver.New(1, 0, 1, "", "")
	if err != nil {
		log.Fatal(err)
	}

	compareResult := version1.CompareTo(version2)
	fmt.Println("Comparison Result:", compareResult)
}
```

In this snippet, `version1` is `1.0.0` and `version2` is `1.0.1`. The result will be -1, indicating `version1` is less than `version2`.

### Example 4: Incrementing Version Numbers

Use methods like `IncreaseMajor`, `IncreaseMinor`, or `IncreasePatch` to increment version numbers.

```go
package main

import (
	"fmt"
	"log"
    "github.com/wernerstrydom/go-semver"
)

func main() {
	version, err := semver.New(1, 0, 0, "", "")
	if err != nil {
		log.Fatal(err)
	}

	version.IncreaseMinor()
	fmt.Println("Increased Minor Version:", version.String())

	version.IncreasePatch()
	fmt.Println("Increased Patch Version:", version.String())
}
```

This example starts with version `1.0.0`, increments the minor version to `1.1.0`, and then increments the patch 
version to `1.1.1`. Each increment is printed.


## Contributing

Contributions to go-semver are welcome. If you find a bug or have a feature request, please open an issue. You can also 
fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License. You can view the full license in the LICENSE file included in the repository.

## Contact

For questions or further information, please contact the repository maintainer through GitHub. Contributions, suggestions, 
and feedback are highly appreciated.