# How to contribute

Thank you for deciding to contribute to this project! We're open to most changes
so long as they're maintainable and within the scope and spirit of this project.

## Reporting bugs, suggesting features, and asking questions

The best place to start is with creating an issue on the GitHub repository - try
to include as much detail as possible.

If you're reporting a bug, make sure to include information like:

- what environment you're using
- how you're calling the tool (including the contents of lockfiles being parsed)
- the _full_ output of the tool
- depending on the issue, outputs of the auxiliary commands can be helpful too

If you're suggesting a feature, make sure to include as much information on your
use-case and why you think the feature belongs in this tool. Also remember that
while we'd love to empower everyone by providing useful features, we do have to
maintain everything in the tool, so if your feature is very specific for your
use-case and complex we might have to say no (but please don't let that stop you
from opening an issue - we want to work _with_ you, which includes trying to
find middle grounds and alternatives if we think something doesn't belong in the
tool itself for some reason).

If you're asking a question, make sure to be as concise as possible, and
remember that we only have limited time so might not always be able to answer
every question.

Always make sure your issue is well formatted - use codeblocks to wrap terminal
output and code and wrap large content dumps in `<details>`. We know no one's
perfect (including ourselves!), but be aware we will edit issues descriptions if
we think they need a touch up.

## Testing

We try and include tests for as much of the codebase as possible, though some of
the more complicated parts (such as the CLI flags) currently don't have tests.

Tests are run as part of CI and are required to pass before a change can be
landed. You can run tests locally with:

```shell
make test
```

## Linting & formatting

We use [`golangci-lint`](https://github.com/golangci/golangci-lint) & `gofmt` to
keep the codebase healthy, which can be run with:

```shell
make lint
```

Currently, there are _6_ unresolved linting errors which are yet to be handled -
because of this, we're currently not running linting as part of CI. Changes
should not include any additional linting errors.

Here are the unresolved errors:

```
detector/database/cache.go:32:15: err113: do not define dynamic errors, use wrapped static errors instead: "errors.New(\"--offline can only be used when a local version of the OSV database is available\")" (goerr113)
                return nil, errors.New("--offline can only be used when a local version of the OSV database is available")
                            ^
detector/database/cache.go:36:30: should rewrite http.NewRequestWithContext or add (*Request).WithContext (noctx)
                req, err := http.NewRequest("GET", db.ArchiveURL, nil)
                                           ^
detector/parsers/parse-composer-lock.go:22:2: Consider preallocating `packages` (prealloc)
        var packages []PackageDetails
        ^
detector/parsers/parse-npm-lock.go:32:2: Consider preallocating `details` (prealloc)
        var details []PackageDetails
        ^
detector/parsers/parse-yarn-lock.go:99:2: Consider preallocating `packages` (prealloc)
        var packages []PackageDetails
        ^
detector/parsers/parsers.go:33:30: err113: do not define dynamic errors, use wrapped static errors instead: "fmt.Errorf(\"cannot parse %s\", path.Base(pathToLockfile))" (goerr113)
                return []PackageDetails{}, fmt.Errorf("cannot parse %s", path.Base(pathToLockfile))
```

Markdown documents and yaml files should ideally be formatted with
[`prettier`](https://prettier.io/) with
[`--prose-wrap always`](https://prettier.io/).

You can run this with:

```shell
npx prettier --prose-wrap always --write .
```

## Submitting changes

Make a pull request on this repository with a clear description of the change
you're made and why. Ideally include a test or two if possible, and keep changes
atomic (one feature per PR, since we squash when merging).

Commit messages should be
[conventional](https://www.conventionalcommits.org/en/v1.0.0/), to make it
easier to write changelogs and determine version numbers when releasing.
