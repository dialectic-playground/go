<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents** _generated with [DocToc](https://github.com/thlorenz/doctoc)_

- [High-Level Design (HLD) — Standard Terminal Hello World](#high-level-design-hld--standard-terminal-hello-world)
  - [1. Introduction](#1-introduction)
    - [1.1 Scope](#11-scope)
    - [1.2 Definitions and acronyms](#12-definitions-and-acronyms)
    - [1.3 References](#13-references)
  - [2. Requirements summary](#2-requirements-summary)
  - [3. Architectural overview](#3-architectural-overview)
    - [3.1 Context diagram](#31-context-diagram)
    - [3.2 Component (module) decomposition](#32-component-module-decomposition)
    - [3.3 Data flow](#33-data-flow)
  - [4. Key design decisions](#4-key-design-decisions)
  - [5. Cross-platform strategy](#5-cross-platform-strategy)
  - [6. Error handling and diagnostics](#6-error-handling-and-diagnostics)
  - [7. Non-functional design](#7-non-functional-design)
  - [8. Risks and open issues](#8-risks-and-open-issues)
  - [9. Approval](#9-approval)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# High-Level Design (HLD) — Standard Terminal Hello World

| Field                   | Value                                                                               |
| :---------------------- | :---------------------------------------------------------------------------------- |
| Project                 | `00-hello-world`                                                                    |
| Tools                   | dp-hello                                                                            |
| Implementation language | go (Go 1.27, standard library only)                                                 |
| Author                  | Brandon Perkins                                                                     |
| Version / date          | 1.0 — 2026-10-13                                                                    |
| Status                  | Approved                                                                            |
| Inputs                  | `specs/URD.md`, `specs/SRS.md`, `specs/user-manual.md`, `shared/CLI-CONVENTIONS.md` |

## 1. Introduction

### 1.1 Scope

`dp-hello` writes `Hello, NAME!` + LF to standard output COUNT times (NAME defaults to
`World`, COUNT to 1, range 1–1000), prints help and version text, and reports usage errors
with the exact CONV-13/CONV-14 diagnostics and exit status 2. As SRS §1.2 states, it reads no
input, takes no operands, and touches no files; it also reads no environment variables and
starts no subprocesses (NFR-07).

### 1.2 Definitions and acronyms

| Term        | Meaning                                                                                                                                    |
| :---------- | :----------------------------------------------------------------------------------------------------------------------------------------- |
| Invocation  | The syntactic result of option parsing: the requested action, the last NAME and COUNT values with their option spellings, and the operands |
| Config      | A validated greeting request: a non-empty NAME and a COUNT in 1–1000                                                                       |
| Usage error | A CONV-13 condition; reported on stderr with the CONV-14 hint, exit status 2                                                               |
| Spelling    | The option as the user wrote it (`-n` or `--name`), needed by the invalid-argument diagnostic                                              |

### 1.3 References

- SRS §3 (interfaces and exit status), §4 (FR-01 … FR-11), §5 (NFR-01 … NFR-07), §6 (diagnostics)
- CLI-CONV CONV-03, CONV-05, CONV-07 … CONV-15, CONV-24, CONV-25
- User manual, "Help output" (normative help text)

## 2. Requirements summary

| Area                  | SRS requirements                                     | Architectural impact                                                                                        |
| :-------------------- | :--------------------------------------------------- | :---------------------------------------------------------------------------------------------------------- |
| Option parsing        | FR-02, FR-06, FR-07, FR-08, FR-10, CONV-07 … CONV-12 | A hand-written scanner (`cli`) that records the last value of each option and stops at `--help`/`--version` |
| Validation            | FR-03, FR-04, FR-05, FR-09                           | A separate `validate` unit that checks only the winning values, after the scan                              |
| Greeting construction | FR-01, FR-02, NFR-05                                 | A pure `greeting` unit, plus the help and version texts as constants                                        |
| Diagnostics           | FR-10, FR-11, §6, CONV-13 … CONV-15                  | One `diagnostics` unit owns every message; errors are values, never panics                                  |
| Output                | NFR-01, NFR-02, CONV-03, CONV-25                     | An `output` unit that writes raw bytes through a bounded buffer                                             |
| Build and install     | NFR-06                                               | Project `Makefile` hooks around `go build`, `go test`, `go vet`, and `gofmt`                                |

## 3. Architectural overview

### 3.1 Context diagram

```text
   argv (UTF-8) ──▶ ┌───────────┐ ──▶ stdout (greeting, help, or version; LF only)
                    │  dp-hello │ ──▶ stderr (usage diagnostic + hint)
                    └───────────┘ ──▶ exit status 0 / 2 (1 only if stdout cannot be written)
```

No stdin, files, environment variables, network, or child processes.

### 3.2 Component (module) decomposition

All units are Go packages under `src/internal/`, so nothing outside this module can import
them; `src/cmd/dp-hello` is the executable.

| Component      | Responsibility                                                                                          | Depends on                                             |
| :------------- | :------------------------------------------------------------------------------------------------------ | :----------------------------------------------------- |
| `cmd/dp-hello` | Process boundary: passes `os.Args[1:]`, `os.Stdout`, `os.Stderr` to `app.Run` and exits with its status | `app`                                                  |
| `app`          | Orchestrates one run; the only place that maps outcomes to exit statuses                                | `cli`, `validate`, `greeting`, `output`, `diagnostics` |
| `cli`          | Scans the arguments into an `Invocation` or a syntax usage error                                        | `diagnostics`                                          |
| `validate`     | Checks operands, NAME, and COUNT; produces a `Config` or a usage error                                  | `cli`, `diagnostics`                                   |
| `greeting`     | Pure text: greeting line, help text, version line                                                       | —                                                      |
| `diagnostics`  | The `UsageError` type and every CONV-13/14 message                                                      | —                                                      |
| `output`       | Writes text to a stream (`io.Writer`), repeated output through a buffer                                 | —                                                      |

### 3.3 Data flow

1. `main` calls `app.Run(os.Args[1:], os.Stdout, os.Stderr)`.
2. `cli.Parse` scans left to right, only collecting operands and recording values. An option
   syntax error returns a `UsageError`; `--help` or `--version` returns immediately with that
   action (CONV-12).
3. For `Help`/`Version`, `app` writes the constant text and exits 0.
4. Otherwise `validate.Validate` reports the first operand, then an empty NAME, then an invalid
   COUNT, or returns a `Config`.
5. `app` writes `greeting.Line(name)` COUNT times via `output.WriteRepeated` and exits 0.
6. Any `UsageError` is written to stderr as `UsageError.Report()` and the status is 2. No
   stdout output has happened at that point (FR-11).

## 4. Key design decisions

| #   | Decision                                                    | Alternatives considered              | Rationale                                                                                                                                                                                                                                               |
| :-- | :---------------------------------------------------------- | :----------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| D1  | Hand-written option scanner over a fixed option table       | `flag` package; third-party parsers  | `flag` does not implement CONV-07 bundling, GNU diagnostics, or the CONV-12 scan order; third-party parsers are not allowed                                                                                                                             |
| D2  | Parse (syntax) and validate (values) are separate passes    | Validate each value as it is scanned | Last-wins (FR-06) means only the final value may be validated, and CONV-12 (curriculum v1.4.0) detects operand and value errors only after the scan: `extra --help` and `--count 0 --help` print the help, while `-x --help` reports the invalid option |
| D3  | `Run(args, stdout, stderr) int` core with a one-line `main` | Logic in `main` using `os` directly  | End-to-end behavior is unit-testable with `bytes.Buffer`; no globals (Principle 16)                                                                                                                                                                     |
| D4  | Errors are values (`*diagnostics.UsageError`), never panics | `panic`/`recover`, `log.Fatal`       | CONV-24 forbids stack traces; explicit returns keep control flow visible                                                                                                                                                                                |
| D5  | Stream repeated greetings through `bufio.Writer`            | Build the whole output in memory     | Memory stays bounded for any NAME length × 1000 (NFR-02)                                                                                                                                                                                                |
| D6  | Standard library only; no `go.sum`                          | —                                    | Nothing to lock or audit; the module is self-contained                                                                                                                                                                                                  |

## 5. Cross-platform strategy

- **Line endings (CONV-03):** every generated line ends in `\n` in Go source constants. Go's
  `os.File` has no text mode, so bytes are written unchanged on Linux, macOS, and Windows;
  there is nothing to switch to binary mode. Console output on Windows is valid UTF-8, which
  Go converts for the console itself.
- **Arguments (CONV-05):** on Windows the Go runtime reads the UTF-16 command line
  (`GetCommandLineW`) and converts each argument to UTF-8 before `main` runs; on POSIX
  `os.Args` holds the raw bytes, interpreted as UTF-8. The code treats arguments as opaque
  strings and echoes them unchanged, so `Zoë` round-trips on every OS.
- **Paths (CONV-19 … CONV-21):** not applicable; the tool uses no paths.
- **Isolation:** the only platform-facing code is `cmd/dp-hello/main.go`; every other unit
  works on strings and `io.Writer` values.

## 6. Error handling and diagnostics

`cli.Parse` and `validate.Validate` return `*diagnostics.UsageError`, built by one constructor
per CONV-13 row (`InvalidShortOption`, `UnrecognizedOption`, …). `UsageError.Report()` adds the
`dp-hello:` prefix and the CONV-14 hint. `app.Run` is the single place that turns a usage
error into stderr output and status 2. A failure to write standard output (for example a
broken pipe) ends the run silently with status 1; CONV-25 leaves that status unspecified, and
on POSIX the Go runtime may already have ended the process with SIGPIPE, which is equally
silent.

## 7. Non-functional design

| Quality attribute             | SRS reference   | How the design achieves it                                                                                                                                                            |
| :---------------------------- | :-------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Performance / memory          | NFR-02, CONV-31 | Native executable, so the 32 MiB limit applies to the whole process (no runtime baseline is subtracted); O(arguments) parsing; buffered streaming output (4 KiB buffer plus one line) |
| Portability                   | NFR-01          | Raw-byte output with LF constants; runtime-provided UTF-8 arguments; identical code on every OS                                                                                       |
| Robustness                    | NFR-03          | No panics: slicing guarded by length checks, overflow-free COUNT parsing, errors returned as values                                                                                   |
| Security / input validation   | NFR-07          | No file, environment, or process APIs are used; every value is validated before use                                                                                                   |
| Maintainability / testability | NFR-04, NFR-05  | One package per responsibility; `greeting.Line` is pure; `app.Run` is tested end to end                                                                                               |
| Build                         | NFR-06          | Every BASE-MAKEFILE hook is defined in the project `Makefile`                                                                                                                         |

## 8. Risks and open issues

| Risk / issue                                                 | Likelihood | Impact | Mitigation                                                                                       |
| :----------------------------------------------------------- | :--------- | :----- | :----------------------------------------------------------------------------------------------- |
| Help text drifts from the user manual                        | L          | H      | Text is a single constant; public test P07 compares it byte for byte                             |
| Windows argument quoting differs from POSIX shells           | L          | M      | The black-box runner starts the tool directly; the Go runtime follows the standard Windows rules |
| Go toolchain versions differ between container and native CI | L          | M      | `make toolchain` requires go1.27.x everywhere; gofmt ships with that toolchain                   |

## 9. Approval

| Role     | Name                                            | Date       | Decision                      |
| :------- | :---------------------------------------------- | :--------- | :---------------------------- |
| Author   | Brandon Perkins                                 | 2026-10-13 | Submitted                     |
| Reviewer | CI (container and native Linux, macOS, Windows) | 2026-10-13 | Approved when all checks pass |
