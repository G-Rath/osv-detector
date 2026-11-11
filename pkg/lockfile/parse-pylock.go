package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/python/pylock"
)

const PylockEcosystem = PipEcosystem

func ParsePylock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, pylock.New(), PylockEcosystem)
}
