# GitHub

## Requirements

- [gh](https://cli.github.com/)

> **Note** You should be authenticated with GitHub using the `gh auth login` command before running the scripts.

## Install

```bash
sunbeam extension install github https://gist.github.com/pomdtr/dc1363d8f641928893ca8d3e670c9c3d
```

## Usage

```bash
sunbeam github # List all repositories
sunbeam github list-prs <repo> # List all pull requests for a repository
```

## Code

```bash
{{#include ./sunbeam-extension}}
```