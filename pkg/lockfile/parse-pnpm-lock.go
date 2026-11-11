package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/javascript/pnpmlock"
)

const PnpmEcosystem = NpmEcosystem

func ParsePnpmLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, pnpmlock.New(), PnpmEcosystem)
}
