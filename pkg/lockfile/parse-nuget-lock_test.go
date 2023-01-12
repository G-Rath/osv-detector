package lockfile_test

import (
	"github.com/g-rath/osv-detector/pkg/lockfile"
	"testing"
)

func TestParseNuGetLock_InvalidVersion(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseNuGetLock("fixtures/nuget/empty.v0.json")

	expectErrContaining(t, err, "unsupported lockfile version")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}
