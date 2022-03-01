# Open Source Vulnerability Detector

An auditing tool for detecting vulnerabilities using the
[Open Source Vulnerability advisory database provided by GitHub](https://github.com/github/advisory-database)

# NOTE ABOUT CURRENT STATE

This tool is still in an alpha state - it should be usable and stable, but there
are still a few things to be fixed and landed.

In particular, right now the detector is using the standard `semver` package for
comparing which is designed strictly for SemVer itself - this means that version
comparing for ecosystems that don't follow the SemVer specification (and \_only
that specification) won't be compared properly.

A new version parser is being written which will fix this and should land within
a few weeks; until then, some ecosystems (e.g. RubyGems) might have incorrect
results or fail completely.

See [#1](https://github.com/G-Rath/osv-detector/issues/1) for updates

## Usage

The detector accepts a path to a "lockfile" which contains information about the
versions of packages:

```shell
osv-detector path/to/my/package-lock.json
osv-detector path/to/my/composer.lock
```

The detector supports parsing the following lockfiles:

| Lockfile             | Ecosystem   | Tool       |
| -------------------- | ----------- | ---------- |
| `package-lock.json`  | `npm`       | `npm`      |
| `composer.lock`      | `Packagist` | `composer` |
| `Gemfile.lock`       | `RubyGems`  | `bundler`  |
| `requirements.txt`\* | `PyPI`      | `pip`      |

\*: `requirements.txt` support is currently very limited - it ignores anything
that is not a direct requirement (e.g. flags or files) & it assumes the _lowest_
version possible for the constraint (or lack of)

The detector will attempt to automatically determine the parser to use based on
the filename - you can manually specify the parser to use with the `-parse-as`
flag:

```shell
osv-detector --parse-as 'package-lock.json' path/to/my/file.lock
```

The detector maintains a cache of the OSV database it's using locally, which is
updated with any changes at the start of every run.

You can have the detector work solely off this cache with the `--offline` flag:

```shell
osv-detector --offline
```

This requires the detector to have successfully cached an offline copy of the
OSV database at least once.
