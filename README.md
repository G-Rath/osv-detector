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
| `cargo.lock`         | `crates.io` | `cargo`    |
| `package-lock.json`  | `npm`       | `npm`      |
| `yarn.lock`          | `npm`       | `yarn`     |
| `pnpm-lock.yaml`     | `npm`       | `pnpm`     |
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

### Auxiliary output commands

The detector supports a few auxiliary commands that have it output information
which can be useful for debugging issues and general exploring.

#### `--list-ecosystems`

Lists all the ecosystems that exist in the loaded OSV database. This can be
useful when exploring new parsers, or building wrappers around the detector
since a valid ecosystem is required to determine if a package has a
vulnerability, and the ecosystem names are case-sensitive:

```
$ osv-detector --list-ecosystems
Loaded 6532 vulnerabilities (including withdrawn, last updated Fri, 04 Mar 2022 00:11:50 GMT)
The loaded OSV has vulnerabilities for the following ecosystems:
  Packagist
  Go
  crates.io
  RubyGems
  npm
  PyPI
  Maven
  NuGet
```

#### `--list-packages`

Lists all the packages that the detector was able to parse out of the given
lockfile. This can be useful for sense-checking parsers and can also be used for
building other tools.

Each package is outputted on its own line, in the format of
`<ecosystem>: <name>@<version>`:

```
$ osv-detector --list-packages /path/to/my/Gemfile.lock
Loaded 6532 vulnerabilities (including withdrawn, last updated Fri, 04 Mar 2022 00:11:50 GMT)
The following packages were found in /path/to/my/Gemfile.lock:
  RubyGems: ast@2.4.2
  RubyGems: parallel@1.21.0
  RubyGems: parser@3.1.1.0
  RubyGems: rainbow@3.1.1
  RubyGems: regexp_parser@2.2.1
  RubyGems: rexml@3.2.5
  RubyGems: rubocop@1.25.1
  RubyGems: rubocop-ast@1.16.0
  RubyGems: ruby-progressbar@1.11.0
  RubyGems: unicode-display_width@2.1.0
```

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
