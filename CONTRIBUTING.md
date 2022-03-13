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
keep the codebase healthy.

These are run as part of CI and are required to pass before a change can be
landed. You can run linting locally with:

```shell
make lint
```

Markdown documents, json, and yaml files should be formatted with
[`prettier`](https://prettier.io/). You can run this with:

```shell
make format-with-prettier
```

This is also run as part of CI.

## Submitting changes

Make a pull request on this repository with a clear description of the change
you're made and why. Ideally include a test or two if possible, and keep changes
atomic (one feature per PR, since we squash when merging).

Commit messages should be
[conventional](https://www.conventionalcommits.org/en/v1.0.0/), to make it
easier to write changelogs and determine version numbers when releasing.

## Releasing a new version

> This section is primarily for maintainers

Releases are done by tagging commits with a semantic version prefixed with a
`v`, which trigger a CI workflow that uses
[`goreleaser`](https://goreleaser.com/) - once it's built a new release, it
creates a new draft release on GitHub and pushes the freshly built artifacts.

A maintainer should fill out the changelog for the release, and then publish it.
Note that the changelog can be edited after a release has been published, so it
doesn't have to be perfect.

Version numbers should be based on what commits are in the new release - `fix:`
commits represent patch versions, and `feat:` commits represent minor versions.
So if a new version is made up of commits that are all `fix:`s, then it should
be patch bump. If there's at least one `feat:` commit, it should be a minor bump
(which resets the patch number to 0).

See the
[conventional commit spec](https://www.conventionalcommits.org/en/v1.0.0/) for
more.

The final version tag should be prefixed with a `v`:

```
v0.1.0
v1.0.0
v1.1.0
```
