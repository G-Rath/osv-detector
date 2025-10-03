package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/python/poetrylock"
)

const PoetryEcosystem = PipEcosystem

func ParsePoetryLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, poetrylock.New(), PoetryEcosystem)
}
