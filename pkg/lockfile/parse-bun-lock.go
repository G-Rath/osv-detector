package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/javascript/bunlock"
)

const BunEcosystem = NpmEcosystem

func ParseBunLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, bunlock.New(), BunEcosystem)
}
