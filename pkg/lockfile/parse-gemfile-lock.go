package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/ruby/gemfilelock"
)

const BundlerEcosystem Ecosystem = "RubyGems"

func ParseGemfileLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, gemfilelock.New(), BundlerEcosystem)
}
