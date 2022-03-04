# Open Source Vulnerability Detector

An auditing tool for detecting vulnerabilities using the
[Open Source Vulnerability advisory database provided by GitHub](https://github.com/github/advisory-database)

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
| `yarn.lock`          | `npm`       | `yarn      |
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

## Version parsing and comparing

Versions are compared using an internal `semver` package which aims to support
any number of components followed by a build string.

Components are numbers broken up by dots, e.g. `1.2.3` has the components
`1, 2, 3`. Anything that is not a number or a dot is considered to be the start
of a build string, and anything afterwards (including numbers and dots) are
likewise considered to be part of the build string.

Versions are compared by their components first, in order. Versions are not
required to have the same number of components to be comparable.

If all components are equal, then the build string is compared (if present).

Build string comparison is not guaranteed to be correct, since they can be in
any format. Generally, the comparer attempts to extract numbers from the build
strings which are then compared.

Here are examples of versions with build strings that _can_ be accurately
compared:

```
1.0.0.beta.2
1.0.0-rc.0
1.0.0.v3
1.0.0a1
```

Currently, characters & words in the build string are not factored into the
comparison - this means e.g. `1.0.0a2` will be considered _greater than_
`1.0.0b1`. Ideally this will be supported in the future.

Versions without a build string are considered to be higher than those with
(provided they have the same components).

Improvements to the build string comparor are welcome!
