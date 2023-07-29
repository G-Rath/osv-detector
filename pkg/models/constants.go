package models

type Ecosystem string

const (
	EcosystemGo            Ecosystem = "Go"
	EcosystemNPM           Ecosystem = "npm"
	EcosystemOSSFuzz       Ecosystem = "OSS-Fuzz"
	EcosystemPyPI          Ecosystem = "PyPI"
	EcosystemRubyGems      Ecosystem = "RubyGems"
	EcosystemCratesIO      Ecosystem = "crates.io"
	EcosystemPackagist     Ecosystem = "Packagist"
	EcosystemMaven         Ecosystem = "Maven"
	EcosystemNuGet         Ecosystem = "NuGet"
	EcosystemLinux         Ecosystem = "Linux"
	EcosystemDebian        Ecosystem = "Debian"
	EcosystemAlpine        Ecosystem = "Alpine"
	EcosystemHex           Ecosystem = "Hex"
	EcosystemAndroid       Ecosystem = "Android"
	EcosystemGitHubActions Ecosystem = "GitHub Actions"
	EcosystemPub           Ecosystem = "Pub"
	EcosystemConanCenter   Ecosystem = "ConanCenter"
	EcosystemRockyLinux    Ecosystem = "Rocky Linux"
	EcosystemAlmaLinux     Ecosystem = "AlmaLinux"
	EcosystemBitnami       Ecosystem = "Bitnami"
)

type RangeType string

const (
	RangeSemVer    RangeType = "SEMVER"
	RangeEcosystem RangeType = "ECOSYSTEM"
	RangeGit       RangeType = "GIT"
)
