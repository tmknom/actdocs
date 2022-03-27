# General Template

A template repository for any languages that keep clean code.

## Description

N/A

## Usage

N/A

## Developer Guide

<!-- markdownlint-disable no-inline-html -->
<details>
<summary>Click to see details</summary>

### Requirements

- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/)
- [GitHub CLI](https://cli.github.com/)

### Development

N/A

### Test

Run the following command:

```shell
make test
```

### CI

When create a pull request, the following workflows are executed automatically at GitHub Actions.

- [Lint Markdown](/.github/workflows/lint-markdown.yml)
- [Lint YAML](/.github/workflows/lint-yaml.yml)
- [Lint Action](/.github/workflows/lint-action.yml)
- [Lint Shell](/.github/workflows/lint-shell.yml)

### Dependency management

Use Dependabot version updates.
For more information, see [dependabot.yml](/.github/dependabot.yml).

### Release management

#### 1. Bump up to a new version

Run the following command to bump up.

```shell
make bump
```

This command will execute the following steps:

1. Update [VERSION](/VERSION)
2. Commit, push, and create a pull request
3. Open the web browser automatically for reviewing pull request

Then review and merge, so the release is ready to go.

#### 2. Publish the new version

Run the following command to publish a new tag at GitHub.

```shell
make release
```

Finally, we can use the new version! :tada:

</details>
<!-- markdownlint-enable no-inline-html -->

## Changelog

See [CHANGELOG.md](/CHANGELOG.md).

## License

Apache 2 Licensed. See [LICENSE](/LICENSE) for full details.
