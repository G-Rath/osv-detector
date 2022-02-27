# Open Source Vulnerability Detector

An auditing tool for detecting vulnerabilities using the Open Source
Vulnerability specification

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
