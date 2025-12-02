package main

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/g-rath/osv-detector/internal/cachedregexp"
)

// normalizeFilePaths attempts to normalize any file paths in the given `output`
// so that they can be compared reliably regardless of the file path separator
// being used.
//
// Namely, escaped forward slashes are replaced with backslashes.
func normalizeFilePaths(t *testing.T, output string) string {
	t.Helper()

	return strings.ReplaceAll(strings.ReplaceAll(output, "\\\\", "/"), "\\", "/")
}

// normalizeRootDirectory attempts to replace references to the current working
// directory with "<rootdir>", in order to reduce the noise of the cmp diff
func normalizeRootDirectory(t *testing.T, str string) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	cwd = normalizeFilePaths(t, cwd)

	// file uris with Windows end up with three slashes, so we normalize that too
	str = strings.ReplaceAll(str, "file:///"+cwd, "file://<rootdir>")
	str = strings.ReplaceAll(str, cwd, "<rootdir>")

	// Replace versions without the root as well
	str = strings.ReplaceAll(str, pathWithoutRoot(t, cwd), "<rootdir>")

	return str
}

// normalizeUserCacheDirectory attempts to replace references to the current working
// directory with "<tempdir>", in order to reduce the noise of the cmp diff
func normalizeUserCacheDirectory(t *testing.T, str string) string {
	t.Helper()

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("could not get user cache (%v) - results and diff might be inaccurate!", err)
	}

	cacheDir = normalizeFilePaths(t, cacheDir)

	// file uris with Windows end up with three slashes, so we normalize that too
	str = strings.ReplaceAll(str, "file:///"+cacheDir, "file://<tempdir>")

	return strings.ReplaceAll(str, cacheDir, "<tempdir>")
}

// normalizeTempDirectory attempts to replace references to the temp directory
// with "<tempdir>", to ensure tests pass across different OSs
func normalizeTempDirectory(t *testing.T, str string) string {
	t.Helper()

	//nolint:gocritic // ensure that the directory doesn't end with a trailing slash
	tempDir := normalizeFilePaths(t, filepath.Join(os.TempDir()))
	re := cachedregexp.MustCompile(regexp.QuoteMeta(tempDir+`/osv-detector-test-`) + `\d+`)
	str = re.ReplaceAllString(str, "<tempdir>")

	// Replace versions without the root as well
	re = cachedregexp.MustCompile(regexp.QuoteMeta(pathWithoutRoot(t, tempDir)+`/osv-detector-test-`) + `\d+`)

	return re.ReplaceAllString(str, "<tempdir>")
}

// normalizeDatabaseStats attempts to replace references to database stats (such as
// the number of vulnerabilities and the time that the database was last updated)
// in the output with %% wildcards, in order to reduce the noise of the cmp diff
func normalizeDatabaseStats(t *testing.T, str string) string {
	t.Helper()

	re := cachedregexp.MustCompile(`(\w+) \(\d+ vulnerabilities, including withdrawn - last updated \w{3}, \d\d \w{3} \d{4} [012]\d:\d\d:\d\d GMT\)`)

	return re.ReplaceAllString(str, "$1 (%% vulnerabilities, including withdrawn - last updated %%)")
}

// normalizeLocalhostPort attempts to replace references to 127.0.0.1:<port>
// with a placeholder, to ensure tests pass when using httptest.Server
func normalizeLocalhostPort(t *testing.T, str string) string {
	t.Helper()

	re := cachedregexp.MustCompile(`127\.0\.0\.1:\d+`)

	return re.ReplaceAllString(str, "<localhost>:<port>")
}

// normalizeErrors attempts to replace error messages on alternative OSs with their
// known linux equivalents, to ensure tests pass across different OSs
func normalizeErrors(t *testing.T, str string) string {
	t.Helper()

	str = strings.ReplaceAll(str, "The system cannot find the path specified.", "no such file or directory")
	str = strings.ReplaceAll(str, "The system cannot find the file specified.", "no such file or directory")

	return str
}

// normalizeSnapshot applies a series of normalizes to the buffer from a std stream like stdout and stderr
func normalizeSnapshot(t *testing.T, str string) string {
	t.Helper()

	for _, normalizer := range []func(t *testing.T, str string) string{
		normalizeFilePaths,
		normalizeRootDirectory,
		normalizeTempDirectory,
		normalizeUserCacheDirectory,
		normalizeDatabaseStats,
		normalizeLocalhostPort,
		normalizeErrors,
	} {
		str = normalizer(t, str)
	}

	return str
}

func pathWithoutRoot(t *testing.T, str string) string {
	t.Helper()

	// Replace versions without the root as well
	var root string
	if runtime.GOOS == "windows" {
		root = filepath.VolumeName(str) + "\\"
	}

	if strings.HasPrefix(str, "/") {
		root = "/"
	}

	return str[len(root):]
}
