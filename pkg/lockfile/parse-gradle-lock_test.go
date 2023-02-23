package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParseGradleLockFile_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGradleLockFile("fixtures/gradle/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGradleLockFile_OnlyComments(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGradleLockFile("fixtures/gradle/only-comments")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGradleLockFile_EmptyStatement(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGradleLockFile("fixtures/gradle/only-empty")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGradleLockFile_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGradleLockFile("fixtures/gradle/one-pkg")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "org.springframework.security:spring-security-crypto",
			Version:   "5.7.3",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
	})
}

func TestParseGradleLockFile_MultiplePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGradleLockFile("fixtures/gradle/5-pkg")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "org.springframework.boot:spring-boot-autoconfigure",
			Version:   "2.7.4",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
		{
			Name:      "org.springframework.boot:spring-boot-configuration-processor",
			Version:   "2.7.5",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
		{
			Name:      "org.springframework.boot:spring-boot-devtools",
			Version:   "2.7.6",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},

		{
			Name:      "org.springframework.boot:spring-boot-starter-aop",
			Version:   "2.7.7",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
		{
			Name:      "org.springframework.boot:spring-boot-starter-data-jpa",
			Version:   "2.7.8",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
	})
}

func TestParseGradleLockFile_WithInvalidLines(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGradleLockFile("fixtures/gradle/with-bad-pkg")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "org.springframework.boot:spring-boot-autoconfigure",
			Version:   "2.7.4",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
		{
			Name:      "org.springframework.boot:spring-boot-configuration-processor",
			Version:   "2.7.5",
			Ecosystem: lockfile.MavenEcosystem,
			CompareAs: lockfile.MavenEcosystem,
		},
	})
}
