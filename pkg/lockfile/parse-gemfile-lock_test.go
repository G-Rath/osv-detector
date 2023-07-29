package lockfile_test

import (
	"testing"

	"github.com/g-rath/osv-detector/pkg/lockfile"
	"github.com/g-rath/osv-detector/pkg/models"
)

func TestParseGemfileLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/does-not-exist")

	expectErrContaining(t, err, "could not read")
	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGemfileLock_NoSpecSection(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/no-spec-section.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGemfileLock_NoGemSection(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/no-gem-section.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGemfileLock_NoGems(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/no-gems.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{})
}

func TestParseGemfileLock_OneGem(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/one-gem.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "ast",
			Version:   "2.4.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
	})
}

func TestParseGemfileLock_SomeGems(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/some-gems.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "coderay",
			Version:   "1.1.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "pry",
			Version:   "0.14.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
	})
}

func TestParseGemfileLock_MultipleGems(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/multiple-gems.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "bundler-audit",
			Version:   "0.9.0.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "coderay",
			Version:   "1.1.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "dotenv",
			Version:   "2.7.6",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "pry",
			Version:   "0.14.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "thor",
			Version:   "1.2.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
	})
}

func TestParseGemfileLock_Rails(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/rails.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "actioncable",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actionmailbox",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actionmailer",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actionpack",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actiontext",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actionview",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "activejob",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "activemodel",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "activerecord",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "activestorage",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "activesupport",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "builder",
			Version:   "3.2.4",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "concurrent-ruby",
			Version:   "1.1.9",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "crass",
			Version:   "1.0.6",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "digest",
			Version:   "3.1.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "erubi",
			Version:   "1.10.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "globalid",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "i18n",
			Version:   "1.10.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "io-wait",
			Version:   "0.2.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "loofah",
			Version:   "2.14.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "mail",
			Version:   "2.7.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "marcel",
			Version:   "1.0.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "mini_mime",
			Version:   "1.1.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "minitest",
			Version:   "5.15.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "net-imap",
			Version:   "0.2.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "net-pop",
			Version:   "0.1.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "net-protocol",
			Version:   "0.1.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "net-smtp",
			Version:   "0.3.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "nio4r",
			Version:   "2.5.8",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "racc",
			Version:   "1.6.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rack",
			Version:   "2.2.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rack-test",
			Version:   "1.1.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rails",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rails-dom-testing",
			Version:   "2.0.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rails-html-sanitizer",
			Version:   "1.4.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "railties",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rake",
			Version:   "13.0.6",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "strscan",
			Version:   "3.0.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "thor",
			Version:   "1.2.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "timeout",
			Version:   "0.2.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "tzinfo",
			Version:   "2.0.4",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "websocket-driver",
			Version:   "0.7.5",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "websocket-extensions",
			Version:   "0.1.5",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "zeitwerk",
			Version:   "2.5.4",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "nokogiri",
			Version:   "1.13.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
	})
}

func TestParseGemfileLock_Rubocop(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/rubocop.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "ast",
			Version:   "2.4.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "parallel",
			Version:   "1.21.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "parser",
			Version:   "3.1.1.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rainbow",
			Version:   "3.1.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "regexp_parser",
			Version:   "2.2.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rexml",
			Version:   "3.2.5",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rubocop",
			Version:   "1.25.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rubocop-ast",
			Version:   "1.16.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "ruby-progressbar",
			Version:   "1.11.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "unicode-display_width",
			Version:   "2.1.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
	})
}

func TestParseGemfileLock_HasLocalGem(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/has-local-gem.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "backbone-on-rails",
			Version:   "1.2.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actionpack",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "actionview",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "activesupport",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "builder",
			Version:   "3.2.4",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "coffee-script",
			Version:   "2.4.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "coffee-script-source",
			Version:   "1.12.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "concurrent-ruby",
			Version:   "1.1.9",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "crass",
			Version:   "1.0.6",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "eco",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "ejs",
			Version:   "1.1.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "erubi",
			Version:   "1.10.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "execjs",
			Version:   "2.8.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "i18n",
			Version:   "1.10.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "jquery-rails",
			Version:   "4.4.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "loofah",
			Version:   "2.14.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "minitest",
			Version:   "5.15.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "racc",
			Version:   "1.6.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rack",
			Version:   "2.2.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rack-test",
			Version:   "1.1.0",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rails-dom-testing",
			Version:   "2.0.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rails-html-sanitizer",
			Version:   "1.4.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "railties",
			Version:   "7.0.2.2",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "rake",
			Version:   "13.0.6",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "thor",
			Version:   "1.2.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "tzinfo",
			Version:   "2.0.4",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "zeitwerk",
			Version:   "2.5.4",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "nokogiri",
			Version:   "1.13.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
		{
			Name:      "eco-source",
			Version:   "1.1.0.rc.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
		},
	})
}

func TestParseGemfileLock_HasGitGem(t *testing.T) {
	t.Parallel()

	packages, err := lockfile.ParseGemfileLock("fixtures/bundler/has-git-gem.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []lockfile.PackageDetails{
		{
			Name:      "hanami-controller",
			Version:   "2.0.0.alpha1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
			Commit:    "027dbe2e56397b534e859fc283990cad1b6addd6",
		},
		{
			Name:      "hanami-utils",
			Version:   "2.0.0.alpha1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
			Commit:    "5904fc9a70683b8749aa2861257d0c8c01eae4aa",
		},
		{
			Name:      "concurrent-ruby",
			Version:   "1.1.7",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
			Commit:    "",
		},
		{
			Name:      "rack",
			Version:   "2.2.3",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
			Commit:    "",
		},
		{
			Name:      "transproc",
			Version:   "1.1.1",
			Ecosystem: models.EcosystemRubyGems,
			CompareAs: models.EcosystemRubyGems,
			Commit:    "",
		},
	})
}
