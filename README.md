# Open Source Vulnerability Detector

An auditing tool for detecting vulnerabilities, powered by advisory databases
that follow the [OSV specification](https://ossf.github.io/osv-schema/).

Currently, this uses the ecosystem databases provided by
[osv.dev](https://osv.dev/).

## Usage

The detector accepts a path to a "lockfile" which contains information about the
versions of packages:

```shell
osv-detector path/to/my/package-lock.json
osv-detector path/to/my/composer.lock

# you can also pass multiple files
osv-detector path/to/my/package-lock.json path/to/my/composer.lock

# or a directory which is expected to contain at least one supported lockfile
osv-detector path/to/my/
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
| `go.mod`             | `Go`        | `go mod`   |
| `pom.xml`\*          | `Maven`     | `maven`    |
| `requirements.txt`\* | `PyPI`      | `pip`      |

\*: `pom.xml` and `requirements.txt` are technically not lockfiles, as they
don't have to specify the complete dependency tree and can have version
constraints/ranges. When parsing these files, the detector will assume the
_lowest_ version possible for non-exact dependencies, and will ignore anything
that is not a dependency specification (e.g. flags or files in the case of
`requirements.txt`, though `<properties>` _is_ supported for `pom.xml`)

The detector will attempt to automatically determine the parser to use for each
file based on the filename - you can manually specify the parser to use for all
files with the `-parse-as` flag:

```shell
osv-detector --parse-as 'package-lock.json' path/to/my/file.lock
```

By default, the detector attempts to detect known vulnerabilities by checking
the versions of packages specified by the parsed lockfile against the versions
specified by the OSVs in the loaded OSV databases, using an internal
semver-based package that aims to minimize false negatives (see
[this section](#version-parsing-and-comparing) for more details about version
handling).

This allows the detector to be very fast and work offline, but does not support
commits which means the detector can report false positives when using git-based
dependencies.

You can disable dynamically loading the ecosystem databases by passing
`--use-dbs=false`.

You can also have the detector use the `osv.dev` API to check for known
vulnerabilities by supplying the `--use-api` flag. The API is very fast,
typically a few hours ahead of the offline databases, and supports commits;
however it currently can produce false negatives for some ecosystems.

> While the API supports commits, the detector currently has limited support for
> extracting them - only the `composer.lock`, `Gemfile.lock`,
> `package-lock.json`, `yarn.lock`, & `pnpm.yaml` parsers include commit details
>
> See [this section](#passing-arbitrary-package-details-advanced-usage) for how
> you can provide the detector with arbitrary commits to check

You cannot use the API in `--offline` mode, but you can use both the offline
databases and the API together; the detector will remove any duplicate results.

Once all the lockfiles have been pared, the detector will then determine all the
databases to load - if `--use-dbs` is `true` (which it is by default) then this
will include ecosystem specific databases based on the parsed packages.

See [this section](#extra-databases) for details on how to configure extra
databases for the detector to use.

> Remotely sourced databases will be cached along with their etag and
> last-modified date for future checks, to determine if those databases need to
> be updated.

By default, the detector will output the results to `stdout` as plain text, and
exit with an error code of `1` if at least one vulnerability is found. See
[here](#ignoring-certain-vulnerabilities) for how to configure the detector to
ignore certain vulnerabilities.

You can use the `--json` flag to have the detector output its results as JSON:

```shell
osv-detector --json path/to/my/package-lock.json
```

This will result in a JSON object being printed to `stdout` with a `results`
property that has an array containing the results of each lockfile that was
passed:

```json
{
  "results": [
    {
      "filePath": "path/to/my/go.mod",
      "parsedAs": "go.mod",
      "packages": [
        {
          "name": "github.com/BurntSushi/toml",
          "version": "1.0.0",
          "ecosystem": "Go",
          "vulnerabilities": [],
          "ignored": []
        }
      ]
    },
    {
      "filePath": "path/to/my/package-lock.json",
      "parsedAs": "package-lock.json",
      "packages": [
        {
          "name": "wrappy",
          "version": "1.0.2",
          "ecosystem": "npm",
          "vulnerabilities": [],
          "ignored": []
        }
      ]
    }
  ]
}
```

Errors are always sent to `stderr` as plain text, even if the `--json` flag is
passed.

### Config files

The detector supports loading configuration details from a YAML file, which
makes it easy to provide advanced settings (such as extra databases), provide a
consistent results whenever the detector is run on a project, and provide an
audit trail of ignored vulnerabilities (through version control).

By default, the detector will look for a `.osv-detector.yaml` or
`.osv-detector.yml` in the same folder as the current lockfile it's checking,
and will _merge_ the config with any flags being passed.

You can also provide a path to a specific config file that will be used for all
lockfiles being checked with the `--config` flag:

```shell
osv-detector --config ruby-ignores.yml path/to/my/first-ruby-project path/to/my/second-ruby-project
```

You can have the detector ignore specific parts of the config with the
`--no-config-ignores` and `--no-config-databases` flags, or ignore any configs
all together with the `--no-config` flag.

#### Ignoring certain vulnerabilities

> Ignored vulnerabilities won't be included in the text output, and won't be
> counted when determined the code to exit with.

You can provide the detector with a list of IDs for OSVs to ignore when checking
lockfiles with the `ignore` property:

```yaml
ignore:
  - GHSA-4 # "Prototype pollution in xyz"
  - GHSA-5 # "RegExp DDoS in abc"
  - GHSA-6 # "Command injection in hjk"
```

You can also use the `--ignore` flag:

```
osv-detector --ignore GHSA-896r-f27r-55mw package-lock.json

# you can pass multiple ignores
osv-detector --ignore GHSA-896r-f27r-55mw --ignore GHSA-74fj-2j2h-c42q package-lock.json
```

Ignores provided via the flag will be combined with any ignores specified in the
loaded config file.

You can use `jq` to generate a list of OSV ids if you want to ignore all current
known vulnerabilities found by the detector:

```shell
osv-detector-t --json . | jq -r  '.results[].packages | map("- " + .vulnerabilities[].id) | unique | sort | .[]'
```

#### Extra Databases

You can configure the detector to use extra databases with the `extra-databases`
property:

```yaml
extra-databases:
  - url: https://github.com/github/advisory-database/archive/refs/heads/main.zip
    name: GitHub Advisory Database
    working-directory: 'advisory-database-main/advisories/github-reviewed' # only load the reviewed advisories
```

Each extra database must have a `url` property which specifies the source of the
database (for local databases this must begin with `file:`), but all other
properties are optional.

The `url` should be either:

- the path to a local directory, in which case the `url` must start with `file:`
- a url for a remote zip archive; if the url does not end with `.zip`, you must
  specify the `type` as `zip`
  - if you host your OSV database as a repository on GitHub, it can be consumed
    as a zip archive
- a url for a rest API that implements the [osv.dev](https://osv.dev/docs/) API
  - the `url` _should_ include the `/v1` e.g. if you wanted to use the `osv.dev`
    staging API you would specify `https://api-staging.osv.dev/v1` as the `url`

> The detector will attempt to detect the type of each database based on the
> above, however you can explicitly provide the type if needed with the `type`
> property (such as if you have a zip archive whose url does not end with
> `.zip`)

For the file based database sources (`dir` and `zip`), the detector will
recursively load all `.json` files from the `working-directory` relative to the
root of the database as OSVs.

> The detector assumes that you trust the source of the configuration file and
> thus the databases that it points to.
>
> If you are using the detector in a way that allows users to provide arbitrary
> config files that you don't trust, you can use the `--no-config-databases`
> flag to have the detector load the rest of the config without any extra
> databases it may define

This is a very powerful feature as it enables you to create custom OSVs that can
be easily consumed by multiple projects and that cover anything you want - for
example you could write OSVs to check if you're using versions of packages that
are now considered end of life.

This can also be useful for drafting new OSVs or modifications to existing ones,
and becomes even more powerful when combined with the ability to pass arbitrary
package details, as you can have custom ecosystems and write custom tools to
handle extracting the package details.

Here are some further examples of `extra-databases`:

```yaml
extra-databases:
  # include a specific osv.dev ecosystem database for any lockfiles being checked
  - url: https://osv-vulnerabilities.storage.googleapis.com/OSS-Fuzz/all.zip
    name: OSS-Fuzz

  # include all the unreviewed advisories
  - url: https://github.com/github/advisory-database/archive/refs/heads/main.zip
    name: GitHub Advisory Database (unreviewed)
    working-directory: 'advisory-database-main/advisories/unreviewed'

  # include the osv staging api
  - url: https://api-staging.osv.dev/v1
    name: GitHub Advisory Database (unreviewed)

  # include a local directory database (relative)
  - url: file:/../relative/path/to/dir
    name: My local database (relative)

  # include a local directory database (root)
  - url: file:////root/path/to/dir
    name: My local database (root)
```

### Offline mode

You can have the detector run purely in offline mode with the `--offline` flag:

```shell
osv-detector --offline path/to/my/file.lock
```

Remotely sourced databases can only be used in offline mode if they have been
cached by the detector as part of a previous run, and API-based databases will
be skipped entirely.

You can have the detector cache the databases for all known ecosystems supported
by the detector for later offline use with the `--cache-all-databases`:

```shell
osv-detector --cache-all-databases
```

This can be useful if you're planning to run the detector over a number of
lockfiles in bulk.

### Passing arbitrary package details (advanced usage)

The detector supports being passed arbitrary package details in CSV form to
check for known vulnerabilities.

This is useful for one-off manual checks (such as when deciding on a new
library), or if you have packages that are not specified in a lock (such as
vendored dependencies). It also allows you to check packages in ecosystems that
the detector doesn't know about, such as `NuGet`.

You can either pass in CSV rows:

```
osv-detector --parse-as csv-row 'npm,@typescript-eslint/types,5.13.0' 'Packagist,sentry/sdk,2.0.4'
```

or you can specify paths to csv files:

```
osv-detector --parse-as csv-file path/to/my/first-csv path/to/my/second-csv
```

Each CSV row must have at least three fields which hold the ecosystem, package
name, and version (or commit) respectively, and CSV files cannot contain a
header.

The `ecosystem` does _not_ have to be one listed by the detector as known,
meaning you can use any ecosystem that [osv.dev](https://osv.dev/) provides.

If the ecosystem field is empty, then the `version` field is expected to be a
commit. In this case, the `package` column is decorative as only the commit is
passed to the API.

> Remember to tell the detector to use the `osv.dev` API via the `--use-api`
> flag if you're wanting to check commits!

You can also omit the version to have the detector list all known
vulnerabilities in the loaded database that apply to the given package:

```
osv-detector --parse-as csv-row 'NuGet,Yarp.ReverseProxy,'
```

While this uses the `--parse-as` flag, these are _not_ considered standard
parsers so the detector will not automatically use them when checking
directories for lockfiles.

### Auxiliary output commands

The detector supports a few auxiliary commands that have it output information
which can be useful for debugging issues and general exploring.

#### `--list-ecosystems`

Lists all the ecosystems that the detector knows about (aka there is a parser
that results in packages from that ecosystem):

```
$ osv-detector --list-ecosystems
The detector supports parsing for the following ecosystems:
  npm
  crates.io
  RubyGems
  Packagist
  Go
  PyPI
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

## Public packages (`pkg/`)

A couple of the core packages that power the detector have been made public to
allow others to use in their own projects, and to help encourage improvements.

These packages will not receive their own versions and so should be constrained
based on their commit hash. They should also be considered as "hopefully
stable" - while they're not expected to change, as the detector evolves it might
become necessary to change their api in a manner that could be breaking to
downstream consumers.

It is hoped that someday these packages will be published independently, once
their apis have proven to be truly stable.

Improvements, feature requests, and bug reports for these packages are welcome!
