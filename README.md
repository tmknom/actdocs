# actdocs

Generate documentation from Actions and Reusable Workflows.

## Description

The actdocs is a utility to generate documentation from GitHub Actions in Markdown format.
It's identified Actions or Reusable Workflows automatically, then formats appropriately.

## Getting Started

Documentation is generated by the following command:

```shell
docker run --rm -v "$(pwd):/work" -w "/work" \
ghcr.io/tmknom/actdocs generate action.yml
```

### Actions

For example, write the following Actions and save it as `action.yml`.

```yaml:action.yml
name: Example Action
description: A example for Actions.
inputs:
  hello:
    default: "Hello, world."
    required: false
    description: "A input value."
  answer:
    default: 42
    required: true
    description: "Answer to the Ultimate Question of Life, the Universe, and Everything."
outputs:
  result:
    value: ${{ steps.main.outputs.result }}
    description: "A output value."
runs:
  using: composite
  steps:
    - id: main
      shell: bash
      run: echo "result=example" >> "${GITHUB_OUTPUT}"
```

Run `actdocs generate action.yml`, the following Markdown is output.

<!-- prettier-ignore-start -->
```markdown:README.md
## Description

A example for Actions.

## Inputs

| Name | Description | Default | Required |
| :--- | :---------- | :------ | :------: |
| hello | A input value. | `Hello, world.` | no |
| answer | Answer to the Ultimate Question of Life, the Universe, and Everything. | `42` | yes |

## Outputs

| Name | Description |
| :--- | :---------- |
| result | A output value. |
```
<!-- prettier-ignore-end -->

These outputs can be sorted or injected into a specified file.
For more information, see [Usage](#usage).

### Reusable Workflows

Simply change the file you specify.

```shell
docker run --rm -v "$(pwd):/work" -w "/work" \
ghcr.io/tmknom/actdocs generate .github/workflows/lint.yml
```

The actdocs automatically switches its behavior for Reusable Workflows.

## Installation

### Pull Docker image

You can pull from Docker Hub or GitHub Packages, whichever you prefer.

**Docker Hub:**

```shell
docker pull tmknom/actdocs
```

**GitHub Packages:**

```shell
docker pull ghcr.io/tmknom/actdocs
```

### Download binary

Download the latest compiled binaries and put it anywhere in your executable path.

- [GitHub Releases](https://github.com/tmknom/actdocs/releases/latest)

### Build from source code

If you have Go 1.18+ development environment:

```shell
git clone https://github.com/tmknom/actdocs
cd actdocs/
make install
actdocs --help
```

## Usage

### Injection

You can inject to existing file.
Write the injection comments to Markdown.

```markdown
<!-- actdocs start -->
<!-- actdocs end -->
```

Use `inject` command with `--file` or `-f` option.

```shell
docker run --rm -v "$(pwd):/work" -w "/work" \
ghcr.io/tmknom/actdocs inject --file README.md action.yml
```

Then, output is injected to the specified file.

> **Note**
>
> `inject` command can be used with `--dry-run` option to check the behavior without overwriting the file.

### Sort

You can sort items by name and required.
Run actdocs with `--sort` or `-s` option.

```shell
docker run --rm -v "$(pwd):/work" -w "/work" \
ghcr.io/tmknom/actdocs generate --sort action.yml
```

Of course, it can be used in combination with `inject` command.

If you prefer to sort in another way, try the following options:

- `--sort-by-name`: sort by name only
- `--sort-by-required`: sort by required only

### Format

You can format to json.
Run actdocs with `--format` option.

```shell
docker run --rm -v "$(pwd):/work" -w "/work" \
ghcr.io/tmknom/actdocs generate --format=json action.yml
```

Supported format is `markdown` and `json`.

### Show help

For full details, run `docker run --rm ghcr.io/tmknom/actdocs --help`.

```shell
Generate documentation from Actions and Reusable Workflows

Usage:
  actdocs [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate documentation
  help        Help about any command
  inject      Inject generated documentation to existing file

Flags:
      --debug              show debugging output
      --format string      output format [markdown json] (default "markdown")
  -h, --help               help for actdocs
      --omit               omit for markdown if item not exists
  -s, --sort               sort items by name and required
      --sort-by-name       sort items by name
      --sort-by-required   sort items by required
  -v, --version            version for actdocs

Use "actdocs [command] --help" for more information about a command.
```

## Maintainer Guide

<!-- markdownlint-disable no-inline-html -->
<details>
<summary>Click to see details</summary>

### Requirements

- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/)
- [GitHub CLI](https://cli.github.com/)

### Development

You can use the `make` command.

**Build**:

```shell
make build
```

**Test**:

```shell
make test
```

**Lint**:

```shell
make lint
```

For more information, run `make help`.

### CI

When create a pull request, the following workflows are executed automatically at GitHub Actions.

- [Test](/.github/workflows/test.yml)
- [Lint Go](/.github/workflows/lint-go.yml)
- [Lint Markdown](/.github/workflows/lint-markdown.yml)
- [Lint YAML](/.github/workflows/lint-yaml.yml)
- [Lint Action](/.github/workflows/lint-action.yml)
- [Lint Shell](/.github/workflows/lint-shell.yml)

### Release

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

Then [releasing workflow with GoReleaser](/.github/workflows/release.yml) is run automatically at GitHub Actions
that executes the following steps.

1. Build executable binaries for Linux, Windows and Darwin
2. Create a new GitHub Release, and publish binaries
3. Push Docker images to Docker Hub and GitHub Packages

Finally, we can use the new version! :tada:

### Administration

#### Package management

- Binaries
  - [GitHub Releases](https://github.com/tmknom/actdocs/releases/latest)
- Docker images
  - [Docker Hub](https://hub.docker.com/repository/docker/tmknom/actdocs)
  - [GitHub Packages](https://github.com/tmknom/actdocs/pkgs/container/actdocs)

#### Dependency management

Use Dependabot version updates.
For more information, see [dependabot.yml](/.github/dependabot.yml).

#### Secrets management

Stored environment secrets for the following environments in this repository.

- **release**
  - `DOCKERHUB_TOKEN`: Personal access token used to log against Docker Hub, and it's used by the [releasing workflow](/.github/workflows/release.yml).

</details>
<!-- markdownlint-enable no-inline-html -->

## Changelog

See [Releases](https://github.com/tmknom/actdocs/releases).

## Source Repository

See [tmknom/actdocs](https://github.com/tmknom/actdocs/).

## License

Apache 2 Licensed. See [LICENSE](/LICENSE) for full details.
