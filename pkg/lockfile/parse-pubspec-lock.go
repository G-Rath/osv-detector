package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/dart/pubspec"
)

const PubEcosystem Ecosystem = "Pub"

func ParsePubspecLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, pubspec.New(), PubEcosystem)
}
