package lockfile_test

import (
	"osv-detector/internal/lockfile"
	"testing"
)

func TestParseMavenLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseMavenLock("fixtures/maven/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseMavenLock_Invalid(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseMavenLock("fixtures/maven/not-pom.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseMavenLock_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseMavenLock("fixtures/maven/empty.xml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseMavenLock_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseMavenLock("fixtures/maven/one-package.xml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "maven-artifact",
			Version:   "1.0.0",
			Ecosystem: lockfile.MavenEcosystem,
		},
	})
}

func TestParseMavenLock_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseMavenLock("fixtures/maven/two-packages.xml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "netty-all",
			Version:   "4.1.42.Final",
			Ecosystem: lockfile.MavenEcosystem,
		},
		{
			Name:      "slf4j-log4j12",
			Version:   "1.7.25",
			Ecosystem: lockfile.MavenEcosystem,
		},
	})
}

func TestParseMavenLock_Interpolation(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseMavenLock("fixtures/maven/interpolation.xml")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "mypackage",
			Version:   "1.0.0",
			Ecosystem: lockfile.MavenEcosystem,
		},
		{
			Name:      "my.package",
			Version:   "2.3.4",
			Ecosystem: lockfile.MavenEcosystem,
		},
	})
}
