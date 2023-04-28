package lockfile

import (
	"github.com/g-rath/osv-detector/internal"
)

type Ecosystem = internal.Ecosystem
type PackageDetails = internal.PackageDetails
type PackageDetailsParser = func(pathToLockfile string) ([]PackageDetails, error)
type PackageDetailsParser2 = func(file ParsableFile) ([]PackageDetails, error)
