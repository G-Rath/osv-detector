package parsers

import "osv-detector/detector"

type Ecosystem = detector.Ecosystem
type PackageDetails = detector.PackageDetails
type PackageDetailsParser = func(pathToLockfile string) ([]PackageDetails, error)
