package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/python/uvlock"
)

const UvEcosystem = PipEcosystem

func ParseUvLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, uvlock.New(), UvEcosystem)
}
