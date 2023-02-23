package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseNuGetLockFile_InvalidVersion(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNuGetLockFile("fixtures/nuget/empty.v0.json")

	expectErrContaining(t, err, "unsupported lockfile version")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}
