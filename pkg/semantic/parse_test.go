package semantic_test

import (
	"errors"
	"testing"

	"github.com/g-rath/osv-detector/pkg/models"
	"github.com/g-rath/osv-detector/pkg/semantic"
)

func supportedEcosystems(t *testing.T) []models.Ecosystem {
	t.Helper()

	return []models.Ecosystem{
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

func TestParse(t *testing.T) {
	t.Parallel()

	ecosystems := supportedEcosystems(t)

	for _, ecosystem := range ecosystems {
		_, err := semantic.Parse("", ecosystem)

		if errors.Is(err, semantic.ErrUnsupportedEcosystem) {
			t.Errorf("'%s' is not a supported ecosystem", ecosystem)
		}
	}
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic - '%s'", r)
		}
	}()

	ecosystems := supportedEcosystems(t)

	for _, ecosystem := range ecosystems {
		semantic.MustParse("", ecosystem)
	}
}

func TestMustParse_Panic(t *testing.T) {
	t.Parallel()

	defer func() { _ = recover() }()

	semantic.MustParse("", "<unknown>")

	// if we reached here, then we can't have panicked
	t.Errorf("function did not panic when given an unknown ecosystem")
}
