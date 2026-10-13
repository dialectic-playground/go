# Master Test Plan (MTP) — Standard Terminal Hello World

| Field | Value |
| :-- | :-- |
| Project | `00-hello-world` |
| Tools | dp-hello |
| Implementation language | go (Go 1.27, standard library only) |
| Author | Brandon Perkins |
| Version / date | 1.0 — 2026-10-13 |
| Standard | Structure adapted from IEEE 829 / ISO/IEC/IEEE 29119-3 |

## 1. Introduction and scope

### 1.1 Objectives
"Tested enough" means: every SRS acceptance criterion (FR-01 … FR-11) is exercised by a unit
or black-box test, every CONV-13 diagnostic is checked byte for byte, the public and hidden
suites pass on Linux, macOS, and Windows, `make lint` is clean, and the RTM shows every
applicable requirement as Verified.

### 1.2 In scope / out of scope
| In scope | Out of scope |
| :-- | :-- |
| Every FR/NFR in the SRS and applicable CONV rule, on Linux, macOS, and Windows | Interactive terminal rendering of non-ASCII names (console fonts) |
| The BASE-MAKEFILE targets, including install/uninstall (`make roundtrip`) | Locales and translations (not specified) |
| Robustness against arbitrary argument lists | Files, stdin, environment (the tool uses none) |

## 2. Test items

| Item | Version / location |
| :-- | :-- |
| dp-hello | `00-hello-world/src/` (module `github.com/dialectic-playground/go/00-hello-world`) on branch `feat/00-hello-world` |
| Unit tests | `src/internal/*/*_test.go` (18 test functions, 97 test results with subtests) |
| Public black-box suite | `sdlc-curriculum/projects/00-hello-world/tests/public/` (10 tests) |
| Hidden black-box suite | `sdlc-curriculum/projects/00-hello-world/tests/private/` (10 tests, run by CI) |

## 3. Test levels and approach

| Level | Purpose | Technique(s) | Tooling | Owner |
| :-- | :-- | :-- | :-- | :-- |
| Unit | Each package in isolation: parser, validator, greeting, diagnostics, output | Equivalence partitioning, boundary-value analysis, table-driven tests | `go test` (`testing` package), `make test` | Implementer |
| Integration | `app.Run` end to end with in-memory stdout/stderr (`bytes.Buffer`) and a failing writer | Specification-based (SRS AC), error guessing | `go test` (`internal/app`) | Implementer |
| System (black-box) | The specified CLI behavior, byte-exact | Specification-based | `make test` → `run-tests.py` | Implementer |
| Acceptance | Hidden suite passes on all OSes | — | CI `make check` (`test-hidden`) | CI |
| Static analysis | Vet checks and formatting | — | `make lint` (`go vet`, `gofmt -l`) | Implementer / CI |

## 4. Test design

### 4.1 Equivalence partitions and boundaries

| Input | Valid partitions | Invalid partitions | Boundaries |
| :-- | :-- | :-- | :-- |
| NAME | ASCII, spaces, punctuation, non-ASCII (`Zoë`), leading `-` | Empty string | Length 0 (invalid), 1 |
| COUNT | Digits only, value 1–1000, leading zeros | Empty, non-digits (`abc`, `-1`, `+1`, ` 1`, `1.0`, Unicode digits), 0, > 1000, very long | 0, 1, 999, 1000, 1001, 40 digits |
| Option syntax | `-nV`, `-n V`, `--name=V`, `--name V`, bundles, repetition | Unknown short/long, abbreviation, missing value, `--help=`/`--version=` | Value is the last token; value is `--` |
| Operands | None | Any operand, before or after `--`, including `-` | First of several operands |
| Scan order | `--help`/`--version` first among actions | Syntax error before `--help` | `--help` after an operand or an invalid value |

### 4.2 Additional test cases (beyond the supplied suite)

| ID | Requirement(s) | Preconditions | Steps / command | Expected result | Level |
| :-- | :-- | :-- | :-- | :-- | :-- |
| T-01 | FR-04 | — | `ParseCount` of `1`, `1000`, `007` | 1, 1000, 7 accepted | unit |
| T-02 | FR-05 | — | `ParseCount` of `0`, `1001`, `""`, `abc`, `+1`, 40 nines | all rejected, no overflow | unit |
| T-03 | FR-03 | — | `--name ''` and `-n ''` | `invalid argument '' for '--name'` / `'-n'`, exit 2 | unit |
| T-04 | FR-06 | — | `-n '' -n B -c 0 -c 1` | `Hello, B!`, exit 0 | integration |
| T-05 | CONV-07 | — | `-nc3` | NAME `c3` (value consumes rest of token) | unit |
| T-06 | FR-02 | — | `-n-x`, `--name=-x`, `--name -x` | NAME `-x` | unit |
| T-07 | CONV-09, FR-09 | — | `-- -n` | `extra operand '-n'`, exit 2 | integration |
| T-08 | CONV-12 | — | `-x --help` | `invalid option -- 'x'`, exit 2 | integration |
| T-09 | CONV-12, FR-07 | — | `-c 0 --help`, `extra --help` | help text, exit 0 | integration |
| T-10 | CONV-05 | — | `-c2 -nZoë`; operand `Zoë`; `-ë` | UTF-8 bytes echoed exactly | integration |
| T-11 | CONV-08 | — | `--nam=Ada`, `--=x` | `unrecognized option '<token>'` | unit |
| T-12 | CONV-25 | Writer that always fails | `Run(nil, failing, stderr)` | status 1, stderr empty | integration |
| T-13 | NFR-02 | — | `-c 1000` | 1000 identical lines | integration |

### 4.3 Negative and robustness testing
Every CONV-13 row is covered by `TestParseRejected`, `TestValidateRejected`, and
`TestRunUsageErrors`, which also assert that stdout stays empty (FR-11). Overflow is covered by
a 40-digit COUNT; a failing output stream by `TestRunOutputFailureIsSilent`. Files, stdin, and
permissions do not apply (the tool uses none). No code path panics: slicing is guarded by
length checks, and every error is a returned value.

### 4.4 Performance and resource testing
NFR-02 is measured in the container with a Python harness that runs
`bin/dp-hello -c 1000 -n <100 000-character name>` 20 times with output to the null device:
average 1.7 ms per run and 11 MiB peak RSS (limits: 100 ms, 32 MiB). Memory is independent of
output size because output is streamed through a 4 KiB buffer.

## 5. Test environments

| Environment | OS / version | Toolchain | How provisioned |
| :-- | :-- | :-- | :-- |
| Local (container) | Fedora 44 host, Debian trixie container | Go 1.27 (`golang:1.27-trixie`), make, python3 | `make …` via `tools/dp-run` and `./Containerfile` |
| CI | ubuntu-latest, macos-latest, windows-latest (Git Bash) | Go 1.27.x via `actions/setup-go`, Python 3.12 | `.github/workflows/ci.yml` (from `shared/ci/cross-platform-ci.yml`) |

## 6. Entry and exit criteria

- **Entry:** HLD and LLD written; `make all` succeeds; `make lint` is clean.
- **Exit:** all unit tests and the public and hidden black-box suites pass in the container and
  natively on all three OSes; `make roundtrip` passes; package statement coverage is at least
  80 % (`go test -cover ./...`); no open severity-1/2 defects; the RTM shows every applicable
  requirement as Verified.
- **Suspension / resumption:** suspend when the toolchain or CI infrastructure fails (not the
  code); resume once `make toolchain` passes again.

## 7. Defect management

Defects are recorded as GitHub issues in `dialectic-playground/go` with the failing command,
expected and actual bytes, and the OS. Each fix adds a unit test that fails before the fix, and
the full CI matrix re-runs on the pull request.

| Severity | Definition |
| :-- | :-- |
| 1 — Critical | Crash, panic, or wrong output/exit status on a mainline path |
| 2 — Major | Wrong diagnostic text or status on an error path, or a platform-specific difference |
| 3 — Minor | Documentation, lint, or build-hook issue without a behavior change |

## 8. Schedule and deliverables

| Milestone | Deliverable | Target date |
| :-- | :-- | :-- |
| Design complete | `hld.md`, `lld.md` | 2026-10-13 |
| Implementation complete | `src/`, unit tests, Makefile hooks | 2026-10-13 |
| Verification complete | RTM all Verified, CI green | 2026-10-13 |

## 9. Risks and contingencies

| Risk | Mitigation |
| :-- | :-- |
| Windows-only failure discovered late | CI runs all three OSes on every push |
| Go version drift between container and native CI | `make toolchain` requires go1.27.x; gofmt and vet ship with Go |
| Help text drift | Single constant, checked byte for byte by P07 |

## 10. Test summary (complete at the end)

| Metric | Value |
| :-- | :-- |
| Unit tests (pass / total) | 97 / 97 (18 test functions with subtests) |
| Coverage (line / branch) | 100 % of statements in `cli`, `validate`, `greeting`, `diagnostics`, `app`; 85.7 % in `output` (the in-loop write-error branch; failures surface at the final flush); `cmd/dp-hello` covered by black-box tests |
| Public black-box (pass / 10) | 10 / 10 |
| Hidden black-box (pass / 10, from CI) | 10 / 10 |
| Open defects by severity | 0 / 0 / 0 |
