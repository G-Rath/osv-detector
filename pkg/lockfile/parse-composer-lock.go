package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/php/composerlock"
)

const ComposerEcosystem Ecosystem = "Packagist"

func ParseComposerLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, composerlock.New(), ComposerEcosystem)
}
