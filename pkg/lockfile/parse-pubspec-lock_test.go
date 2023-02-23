package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
)

func TestParsePubspecLockFile_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePubspecLockFile_InvalidYaml(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/not-yaml.txt")

	expectErrContaining(t, err, "could not parse")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePubspecLockFile_Empty(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/empty.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePubspecLockFile_NoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/no-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParsePubspecLockFile_OnePackage(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/one-package.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "back_button_interceptor",
			Version:   "6.0.1",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
	})
}

func TestParsePubspecLockFile_OnePackageDev(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/one-package-dev.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "build_runner",
			Version:   "2.2.1",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
	})
}

func TestParsePubspecLockFile_TwoPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/two-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "shelf",
			Version:   "1.3.2",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
		{
			Name:      "shelf_web_socket",
			Version:   "1.0.2",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
	})
}

func TestParsePubspecLockFile_MixedPackages(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/mixed-packages.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "back_button_interceptor",
			Version:   "6.0.1",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
		{
			Name:      "build_runner",
			Version:   "2.2.1",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
		{
			Name:      "shelf",
			Version:   "1.3.2",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
		{
			Name:      "shelf_web_socket",
			Version:   "1.0.2",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
		},
	})
}

func TestParsePubspecLockFile_PackageWithGitSource(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/source-git.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flutter_rust_bridge",
			Version:   "1.32.0",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "e5adce55eea0b74d3680e66a2c5252edf17b07e1",
		},
		{
			Name:      "screen_retriever",
			Version:   "0.1.2",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "406b9b038b2c1d779f1e7bf609c8c248be247372",
		},
		{
			Name:      "tray_manager",
			Version:   "0.1.8",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "3aa37c86e47ea748e7b5507cbe59f2c54ebdb23a",
		},
		{
			Name:      "window_manager",
			Version:   "0.2.7",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "88487257cbafc501599ab4f82ec343b46acec020",
		},
		{
			Name:      "toggle_switch",
			Version:   "1.4.0",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "",
		},
	})
}

func TestParsePubspecLockFile_PackageWithSdkSource(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/source-sdk.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "flutter_web_plugins",
			Version:   "0.0.0",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "",
		},
	})
}

func TestParsePubspecLockFile_PackageWithPathSource(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParsePubspecLockFile("fixtures/pub/source-path.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "maa_core",
			Version:   "0.0.1",
			Ecosystem: lockfile.PubEcosystem,
			CompareAs: lockfile.PubEcosystem,
			Commit:    "",
		},
	})
}
