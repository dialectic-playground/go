<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents** _generated with [DocToc](https://github.com/thlorenz/doctoc)_

- [go](#go)
  - [Go at a glance](#go-at-a-glance)
  - [Toolchain](#toolchain)
  - [Layout](#layout)
  - [Development](#development)
  - [Implementing a project](#implementing-a-project)
  - [Contributing](#contributing)
  - [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/dialectic-playground/.github/main/assets/go/dist/logo-dark.svg">
    <img alt="dialectic playground / go" src="https://raw.githubusercontent.com/dialectic-playground/.github/main/assets/go/dist/logo.svg" height="96">
  </picture>
</p>

# go

Go implementation of the Dialectic Playground SDLC curriculum.

📖 **Documentation:** <https://dialectic-playground.github.io/docs/reference/languages/go/>

This repository implements the eleven projects of the
[SDLC curriculum](https://github.com/dialectic-playground/sdlc-curriculum) in Go. Each project is
specified to the byte by the curriculum and verified by its black-box tests on Linux, macOS,
and Windows.

**Status:** Implemented projects: 00-hello-world.

## Go at a glance

| Aspect                  | Go                                                            |
| :---------------------- | :------------------------------------------------------------ |
| File extension          | `.go`                                                         |
| CLI classification      | Included                                                      |
| Image category          | First-Class                                                   |
| Access to OS primitives | High (wrapped for async goroutine scheduler)                  |
| Entry-point mechanism   | Compiled-in runtime sets up engine and slices args to os.Args |

## Toolchain

|                      |                                                                           |
| :------------------- | :------------------------------------------------------------------------ |
| Release line         | Go 1.27                                                                   |
| Build / package tool | Go modules                                                                |
| Container image      | `golang:1.27-trixie`, pinned by digest ([`Containerfile`](Containerfile)) |
| Native CI            | actions/setup-go (1.27.x)                                                 |

`make toolchain` verifies that the container and native toolchains are on this release line.

| Task       | Tool                                                        | Project hook                                     |
| :--------- | :---------------------------------------------------------- | :----------------------------------------------- |
| Build      | `go build -trimpath -o bin/<tool>`                          | `make all` (executables in the project's `bin/`) |
| Unit tests | `go test ./...` (standard `testing` package)                | `make test` (then the public black-box suite)    |
| Lint       | `go vet ./...` and `gofmt -l` (fails if any file is listed) | `make lint`                                      |
| Format     | `gofmt -w`                                                  | `make format`                                    |

All of these ship with the Go distribution, so the container image and the native CI need
nothing beyond Go itself, and the formatter is always the toolchain's own version. Each
project is a Go module (`go.mod` in the project directory) with its command in
`src/cmd/<tool>/` and its units in `src/internal/`. The projects use only the standard
library, so there is no `go.sum`; Dependabot (`gomod`) watches each module.

## Layout

| Path                                  | Contents                                                                                                          |
| :------------------------------------ | :---------------------------------------------------------------------------------------------------------------- |
| `00-hello-world/` … `10-sys-monitor/` | One workspace per curriculum project: `Makefile` (with `DP_*` hooks), HLD, LLD, RTM, master test plan, and `src/` |
| `sdlc-curriculum/`                    | The curriculum, as a git submodule pinned to a release tag                                                        |
| `Makefile`                            | `make check` (repository, toolchain, implemented projects), `toolchain`, `projects`, `roundtrip`                  |
| `Containerfile`, `.devcontainer/`     | The toolchain image and the matching dev container                                                                |
| `tools/dp-run`, `tools/container.mk`  | Shared container tooling (copies of the canonical files in `dialectic-playground/.github`)                        |
| `.github/workflows/ci.yml`            | CI: container and native (Linux, macOS, Windows) `make` targets, shared-tooling check, DCO                        |

## Development

Container-first: `make` targets and `tools/dp-run <command>` run in the toolchain container
via Podman, then Docker, then the local system. `DP_RUNTIME=podman|docker|local` overrides
the choice.

```sh
git clone --recurse-submodules https://github.com/dialectic-playground/go.git
cd go
make check                    # repository, toolchain, and every implemented project
DP_RUNTIME=local make check   # the same with your native toolchain
make -C 00-hello-world help   # a project's configuration and hooks
```

## Implementing a project

1. Complete the project's `hld.md` and `lld.md` before writing code.
2. Write the implementation and its unit tests under the project's `src/`, add the
   repository's Go modules configuration, and fill in `BIN_DIR` and the `DP_*` hooks in
   the project's `Makefile`.
3. Pass `make -C <project> test` (container) and `DP_RUNTIME=local make -C <project> test`.
4. Add the project directory to `IMPLEMENTED` in the root `Makefile`. CI then builds, lints,
   and tests it in the container and natively on all three operating systems.
5. Keep `rtm.md` and `master-test-plan.md` current.

Specifications: [curriculum reference](https://dialectic-playground.github.io/docs/reference/curriculum/). Language setup:
[Set up a language repository](https://dialectic-playground.github.io/docs/how-to/set-up-a-language-repository/).

## Contributing

See the [organization contributing guide](https://github.com/dialectic-playground/.github/blob/main/CONTRIBUTING.md).
Pull requests are squash-merged, titles follow Conventional Commits, and every commit needs a
DCO sign-off.

## License

[MIT](LICENSE)
