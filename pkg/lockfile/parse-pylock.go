package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/dart/pubspec"
)

const PylockEcosystem = PipEcosystem

func ParsePylock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, pubspec.New(), PubEcosystem)
}
