package lockfile

import (
	"github.com/g-rath/osv-detector/pkg/models"
)

type Ecosystem = models.Ecosystem
type PackageDetails = models.PackageInfo
type PackageDetailsParser = func(pathToLockfile string) ([]PackageDetails, error)
