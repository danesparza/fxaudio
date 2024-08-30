// The version package provides a location to set the release versions for all
// packages to consume, without creating import cycles.
//
// This package should not import any other fxaudio packages.
package version

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

// Version is the main version number that is being run at the moment.
var Version = "1.1"

// BuildNumber is the build number.  Set during build.  Empty for local dev
var BuildNumber = "0"

// CommitID is the commit information.  Set during build.  Empty for local dev
var CommitID string

// Prerelease is a pre-release marker for the version. If this is "-" (dash)
// then it means that it is a final release. Otherwise, this is a pre-release
// such as "dev" (in development), "beta", "rc1", etc.
var Prerelease = "dev"

// SemVer is an instance of version.Version. This has the secondary
// benefit of verifying during tests and init time that our version is a
// proper semantic version, which should always be the case.
var SemVer *version.Version

func getFormattedVersion() string {
	return fmt.Sprintf("%s.%s", Version, BuildNumber)
}

func init() {
	SemVer = version.Must(version.NewVersion(getFormattedVersion()))
}

// Header is the header name used to send the current version
// in http requests.
const Header = "fxaudio-service-version"

// String returns the complete version string, including prerelease
func String() string {
	if Prerelease != "-" {
		return fmt.Sprintf("%s-%s", getFormattedVersion(), Prerelease)
	}
	return getFormattedVersion()
}
