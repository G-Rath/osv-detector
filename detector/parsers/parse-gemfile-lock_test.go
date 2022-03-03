package parsers_test

import (
	"osv-detector/detector/parsers"
	"strings"
	"testing"
)

func TestParseGemfileLock_FileDoesNotExist(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/does-not-exist")

	if err == nil {
		t.Errorf("Expected to get error, but did not")
	}

	if !strings.Contains(err.Error(), "could not read") {
		t.Errorf("Expected to get \"could not read\" error, but got \"%v\"", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGemfileLock_InvalidJson(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/not-json.txt")

	if err == nil {
		t.Errorf("Expected to get error, but did not")
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGemfileLock_NoSpecSection(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/no-spec-section.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGemfileLock_NoGemSection(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/no-gem-section.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGemfileLock_NoGems(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/no-gems.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{})
}

func TestParseGemfileLock_OneGem(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/one-gem.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "ast",
			Version:   "2.4.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
	})
}

func TestParseGemfileLock_SomeGems(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/some-gems.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "coderay",
			Version:   "1.1.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "pry",
			Version:   "0.14.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
	})
}

func TestParseGemfileLock_MultipleGems(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/multiple-gems.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "bundler-audit",
			Version:   "0.9.0.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "coderay",
			Version:   "1.1.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "dotenv",
			Version:   "2.7.6",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "pry",
			Version:   "0.14.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "thor",
			Version:   "1.2.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
	})
}

func TestParseGemfileLock_Rails(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/rails.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "actioncable",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actionmailbox",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actionmailer",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actionpack",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actiontext",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actionview",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "activejob",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "activemodel",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "activerecord",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "activestorage",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "activesupport",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "builder",
			Version:   "3.2.4",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "concurrent-ruby",
			Version:   "1.1.9",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "crass",
			Version:   "1.0.6",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "digest",
			Version:   "3.1.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "erubi",
			Version:   "1.10.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "globalid",
			Version:   "1.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "i18n",
			Version:   "1.10.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "io-wait",
			Version:   "0.2.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "loofah",
			Version:   "2.14.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "mail",
			Version:   "2.7.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "marcel",
			Version:   "1.0.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "mini_mime",
			Version:   "1.1.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "minitest",
			Version:   "5.15.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "net-imap",
			Version:   "0.2.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "net-pop",
			Version:   "0.1.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "net-protocol",
			Version:   "0.1.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "net-smtp",
			Version:   "0.3.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "nio4r",
			Version:   "2.5.8",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "racc",
			Version:   "1.6.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rack",
			Version:   "2.2.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rack-test",
			Version:   "1.1.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rails",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rails-dom-testing",
			Version:   "2.0.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rails-html-sanitizer",
			Version:   "1.4.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "railties",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rake",
			Version:   "13.0.6",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "strscan",
			Version:   "3.0.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "thor",
			Version:   "1.2.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "timeout",
			Version:   "0.2.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "tzinfo",
			Version:   "2.0.4",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "websocket-driver",
			Version:   "0.7.5",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "websocket-extensions",
			Version:   "0.1.5",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "zeitwerk",
			Version:   "2.5.4",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "nokogiri",
			Version:   "1.13.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
	})
}

func TestParseGemfileLock_Rubocop(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/rubocop.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "ast",
			Version:   "2.4.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "parallel",
			Version:   "1.21.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "parser",
			Version:   "3.1.1.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rainbow",
			Version:   "3.1.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "regexp_parser",
			Version:   "2.2.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rexml",
			Version:   "3.2.5",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rubocop",
			Version:   "1.25.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rubocop-ast",
			Version:   "1.16.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "ruby-progressbar",
			Version:   "1.11.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "unicode-display_width",
			Version:   "2.1.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
	})
}

func TestParseGemfileLock_HasLocalGem(t *testing.T) {
	t.Parallel()

	packages, err := parsers.ParseGemfileLock("fixtures/bundler/has-local-gem.lock")

	if err != nil {
		t.Errorf("Got unexpected error: %v", err)
	}

	expectPackages(t, packages, []parsers.PackageDetails{
		{
			Name:      "backbone-on-rails",
			Version:   "1.2.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actionpack",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "actionview",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "activesupport",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "builder",
			Version:   "3.2.4",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "coffee-script",
			Version:   "2.4.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "coffee-script-source",
			Version:   "1.12.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "concurrent-ruby",
			Version:   "1.1.9",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "crass",
			Version:   "1.0.6",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "eco",
			Version:   "1.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "ejs",
			Version:   "1.1.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "erubi",
			Version:   "1.10.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "execjs",
			Version:   "2.8.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "i18n",
			Version:   "1.10.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "jquery-rails",
			Version:   "4.4.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "loofah",
			Version:   "2.14.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "method_source",
			Version:   "1.0.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "minitest",
			Version:   "5.15.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "racc",
			Version:   "1.6.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rack",
			Version:   "2.2.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rack-test",
			Version:   "1.1.0",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rails-dom-testing",
			Version:   "2.0.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rails-html-sanitizer",
			Version:   "1.4.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "railties",
			Version:   "7.0.2.2",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "rake",
			Version:   "13.0.6",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "thor",
			Version:   "1.2.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "tzinfo",
			Version:   "2.0.4",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "zeitwerk",
			Version:   "2.5.4",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "nokogiri",
			Version:   "1.13.3",
			Ecosystem: parsers.BundlerEcosystem,
		},
		{
			Name:      "eco-source",
			Version:   "1.1.0.rc.1",
			Ecosystem: parsers.BundlerEcosystem,
		},
	})
}
