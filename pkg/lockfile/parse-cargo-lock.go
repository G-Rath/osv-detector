package lockfile

import (
	"github.com/google/osv-scalibr/extractor/filesystem/language/rust/cargolock"
)

const CargoEcosystem Ecosystem = "crates.io"

func ParseCargoLock(pathToLockfile string) ([]PackageDetails, error) {
	return extract(pathToLockfile, cargolock.New(), CargoEcosystem)
}
