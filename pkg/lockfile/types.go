package lockfile

import (
	"github.com/g-rath/osv-detector/internal"
	"io"
)

type Ecosystem = internal.Ecosystem
type PackageDetails = internal.PackageDetails
type PackageDetailsParser = func(pathToLockfile string) ([]PackageDetails, error)
type PackageDetailsParserWithReader = func(r io.Reader) ([]PackageDetails, error)
