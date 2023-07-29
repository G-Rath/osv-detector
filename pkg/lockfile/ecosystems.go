package lockfile

import "github.com/g-rath/osv-detector/pkg/models"

func KnownEcosystems() []Ecosystem {
	return []Ecosystem{
		models.EcosystemNPM,
		models.EcosystemNuGet,
		models.EcosystemCratesIO,
		models.EcosystemRubyGems,
		models.EcosystemPackagist,
		models.EcosystemGo,
		models.EcosystemHex,
		models.EcosystemMaven,
		models.EcosystemPyPI,
		models.EcosystemPub,
	}
}
