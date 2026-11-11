package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/javascript/yarnlock"
)

const YarnEcosystem = NpmEcosystem

func ParseYarnLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, yarnlock.New(), YarnEcosystem)
}
