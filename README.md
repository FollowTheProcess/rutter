# Rutter

[![License](https://img.shields.io/github/license/FollowTheProcess/rutter)](https://github.com/FollowTheProcess/rutter)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/rutter?logo=github&sort=semver)](https://github.com/FollowTheProcess/rutter)
[![CI](https://github.com/FollowTheProcess/rutter/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/rutter/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/rutter/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/rutter)

Sail through your shell history ⚓

> [!WARNING]
> **Rutter is in early development and is not yet ready for use**

![caution](https://assets.followtheprocess.codes/shared/caution-gopher.png)

## Project Description

A [rutter] was a mariner's private notebook of previously sailed routes, kept so the voyage could be repeated before the
advent of nautical charts.

In the same vein, `rutter` is your private notebook of previously executed shell commands. If you're familiar with things like [atuin], [mcfly] and
friends, this is basically the same concept.

Which begs the question... why reinvent the ~~wheel~~ helm?

- I used [atuin] for a long time, it's great, but I use less than 1% of it's functionality
- This seems like a fun thing to implement, and I can mould it to _exactly_ my requirements
- Never really worked with SQLite, sqlc etc. before outside of throwaway projects

## Installation

Compiled binaries for all supported platforms can be found in the [GitHub release]. There is also a [homebrew] tap:

```bash
brew install --cask FollowTheProcess/tap/rutter
```

## Quickstart

### Credits

This package was created with [copier] and the [FollowTheProcess/go-template] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go-template]: https://github.com/FollowTheProcess/go-template
[GitHub release]: https://github.com/FollowTheProcess/rutter/releases
[homebrew]: https://brew.sh
[rutter]: https://en.wikipedia.org/wiki/Rutter_(nautical)
[atuin]: https://atuin.sh
[mcfly]: https://github.com/cantino/mcfly
