# Requirements Traceability Matrix (RTM) — Standard Terminal Hello World

| Field | Value |
| :-- | :-- |
| Project | `00-hello-world` |
| Tools | dp-hello |
| Author | Brandon Perkins |
| Version / date | 1.0 — 2026-10-13 |

Code paths are relative to `00-hello-world/src/`. Unit tests run with `go test ./...`
(`make test`) and are cited as `TestFunction` or `TestFunction/subtest`; public black-box tests
are `P01` … `P10`.

## 1. User requirements to system requirements

| User requirement (URD) | Derived system requirements (SRS) |
| :-- | :-- |
| UR-01 | FR-01 |
| UR-02 | FR-02 |
| UR-03 | FR-04 |
| UR-04 | FR-07 |
| UR-05 | FR-08 |
| UR-06 | FR-03, FR-05, FR-09, FR-10, FR-11 |
| UR-07 | FR-01, SRS §3.7 |
| UR-08 | NFR-01, CONV-03 |
| UR-09 | FR-02, CONV-05 |
| UR-10 | FR-06, FR-10, CONV-07 … CONV-15 |
| UR-11 | NFR-02 |
| UR-12 | NFR-06 |

## 2. Forward traceability

| Req. ID | Summary | Design (HLD/LLD) | Code (file::function) | Verification (unit / black-box) | Status |
| :-- | :-- | :-- | :-- | :-- | :-- |
| UR-01 … UR-12 | User needs | HLD §2 | via the FR/NFR/CONV rows below | via the FR/NFR/CONV rows below | Verified |
| FR-01 | Default greeting | HLD §3.3, LLD §3.3, §3.5 | `internal/validate/validate.go::Validate`, `internal/greeting/greeting.go::Line` | `TestValidateDefaults`, `TestLine`, `TestRunSuccess/default_greeting`, P01 | Verified |
| FR-02 | Named greeting, verbatim | HLD D1, LLD §3.1, §3.2, §3.5 | `internal/cli/cli.go::Parse`, `::parseShortBundle`, `internal/greeting/greeting.go::Line` | `TestParseAccepted`, `TestLine`, `TestRunSuccess/named`, `TestRunSuccess/non-ASCII_name`, P02–P04, P06 | Verified |
| FR-03 | Empty name rejected | LLD §3.3 | `internal/validate/validate.go::Validate` | `TestValidateRejected/empty_name_(long)`, `TestValidateRejected/empty_name_(short)`, `TestRunUsageErrors/empty_name` | Verified |
| FR-04 | Repetition 1–1000, leading zeros | LLD §3.4, §3.6 | `internal/validate/validate.go::ParseCount`, `internal/output/output.go::WriteRepeated` | `TestParseCountAccepted`, `TestWriteRepeated`, `TestRunSuccess/maximum_count`, `TestRunSuccess/leading_zeros`, P05, P06 | Verified |
| FR-05 | Invalid count rejected | LLD §3.3, §3.4 | `internal/validate/validate.go::ParseCount`, `::Validate` | `TestParseCountRejected`, `TestValidateRejected`, `TestRunUsageErrors/count_zero`, `…/count_non-digit`, `…/count_above_maximum` | Verified |
| FR-06 | Last occurrence wins | HLD D2, LLD §3.1 | `internal/cli/cli.go::store` | `TestParseAccepted/last_occurrence_wins`, `TestParseAccepted/count_last_wins_with_spelling`, `TestRunSuccess/last_wins;_earlier_invalid_values_ignored` | Verified |
| FR-07 | Help text, scan order | HLD §3.3, LLD §3.1, §3.5 | `internal/cli/cli.go::Parse`, `internal/greeting/greeting.go::HelpText`, `internal/app/app.go::Run` | `TestParseAccepted/help_stops_the_scan`, `TestHelpTextLayout`, `TestRunSuccess/help_ends_the_scan`, `TestRunSuccess/help_beats_earlier_operands`, P07 | Verified |
| FR-08 | Version line | LLD §3.5 | `internal/greeting/greeting.go::VersionLine`, `internal/app/app.go::Run` | `TestVersionLine`, `TestRunSuccess/version`, P08 | Verified |
| FR-09 | Operands rejected | LLD §3.3 | `internal/validate/validate.go::Validate`, `internal/diagnostics/diagnostics.go::ExtraOperand` | `TestValidateRejected/operand_first`, `TestRunUsageErrors/operand`, `TestRunUsageErrors/operand_after_terminator`, P10 | Verified |
| FR-10 | Unknown and malformed options | LLD §2.3, §3.1, §3.2 | `internal/cli/cli.go::Parse`, `::parseShortBundle`, `internal/diagnostics/diagnostics.go` | `TestParseRejected`, `TestMessages`, `TestRunUsageErrors`, P09 | Verified |
| FR-11 | No stdout on error | HLD §3.3, LLD §3.7 | `internal/app/app.go::Run` | `TestRunUsageErrors` (asserts empty stdout), P09, P10 | Verified |
| NFR-01 | Byte-identical, LF-only output on 3 OSes | HLD §5, LLD §5 | `internal/greeting/greeting.go`, `internal/output/output.go` | `TestHelpTextLayout`, P01–P10 in CI on Linux, macOS, Windows | Verified |
| NFR-02 | < 100 ms, < 32 MiB | HLD D5, LLD §6 | `internal/output/output.go::WriteRepeated` | MTP §4.4 measurement; `TestRunSuccess/maximum_count` | Verified |
| NFR-03 | No crash, panic, or stack trace | HLD D4, §7 | all packages (errors as values) | `TestParseRejected`, `TestParseCountRejected` (40-digit input), `TestRunOutputFailureIsSilent`, P09, P10 | Verified |
| NFR-04 | Separate parse, validate, greet, output units | HLD §3.2 | `internal/cli`, `internal/validate`, `internal/greeting`, `internal/output` | Design review (LLD §8) | Verified |
| NFR-05 | Pure, unit-tested greeting construction | LLD §3.5 | `internal/greeting/greeting.go::Line` | `TestLine` | Verified |
| NFR-06 | Every BASE-MAKEFILE target works | HLD §2 | `00-hello-world/Makefile` (DP_* hooks) | CI: `make check roundtrip` (container and native 3 OSes); `make format`, `make clean distclean` locally | Verified |
| NFR-07 | No files, environment, or subprocesses | HLD §1.1, §7 | `cmd/dp-hello/main.go::main` (only `os.Args`, `os.Stdout`, `os.Stderr`, `os.Exit`) | Inspection | Verified |
| CONV-01 | Canonical name in messages | LLD §2.3 | `internal/diagnostics/diagnostics.go::Tool`, `::Report` | `TestReportFormatsBothLines`, P07–P10 | Verified |
| CONV-02 | Version line format | LLD §3.5 | `internal/greeting/greeting.go::VersionLine` | `TestVersionLine`, P08 | Verified |
| CONV-03 | LF output, no text mode | HLD §5 | `internal/output/output.go`, `cmd/dp-hello/main.go` | P01–P10 (Windows CI) | Verified |
| CONV-05 | UTF-8 arguments, echoed exactly | HLD §5, LLD §5 | `cmd/dp-hello/main.go::main` (`os.Args`) | `TestRunSuccess/non-ASCII_name`, `TestRunUsageErrors/non-ASCII_operand`, `TestParseRejected/unknown_short_option_is_a_whole_letter`; CI on Windows | Verified |
| CONV-06 | Locale-independent, no colors | HLD §1.1 | all packages (no locale or terminal APIs) | Inspection; P01–P10 run with a minimal environment | Verified |
| CONV-07 | Short options, attached values | LLD §3.2 | `internal/cli/cli.go::parseShortBundle` | `TestParseAccepted/short_attached`, `…/bundled_values`, `…/bundle_consumes_rest_of_token`, P06 | Verified |
| CONV-08 | Long options, no abbreviations | LLD §3.1 | `internal/cli/cli.go::Parse`, `::findLong` | `TestParseAccepted/long_inline`, `TestParseRejected/no_abbreviations`, P04 | Verified |
| CONV-09 | `--` terminator | LLD §3.1 | `internal/cli/cli.go::Parse` | `TestParseAccepted/terminator_makes_operands`, `TestRunUsageErrors/operand_after_terminator` | Verified |
| CONV-10 | Operands, `-` is an operand | LLD §3.1 | `internal/cli/cli.go::Parse` | `TestParseAccepted/single_dash_is_an_operand`, `…/operands_interleaved`, P10 | Verified |
| CONV-11 | Last occurrence wins | LLD §3.1 | `internal/cli/cli.go::store` | `TestParseAccepted/last_occurrence_wins` | Verified |
| CONV-12 | Scan order: only syntax errors precede `--help`/`--version`; operand and value errors after the scan | HLD D2, LLD §4 | `internal/cli/cli.go::Parse`, `internal/validate/validate.go::Validate`, `internal/app/app.go::Run` | `TestParseRejected/error_before_help_wins`, `TestParseAccepted/help_after_operand`, `TestRunSuccess/help_beats_earlier_operands`, `TestRunSuccess/help_beats_an_earlier_invalid_long_value`, `TestRunSuccess/version_beats_earlier_operands_and_empty_name`, `TestRunUsageErrors/error_before_help` | Verified |
| CONV-13 | Exact usage messages | LLD §2.3 | `internal/diagnostics/diagnostics.go` | `TestMessages`, `TestParseRejected`, `TestValidateRejected`, P09, P10 | Verified |
| CONV-14 | Hint line, exit 2 | LLD §2.3, §3.7 | `internal/diagnostics/diagnostics.go::Report`, `internal/app/app.go::reportUsage` | `TestReportFormatsBothLines`, `TestRunUsageErrors`, P09, P10 | Verified |
| CONV-15 | Values verbatim, empty as `''` | LLD §2.3 | `internal/diagnostics/diagnostics.go::InvalidArgument` | `TestMessages/empty_invalid_argument`, `TestValidateRejected/count_empty` | Verified |
| CONV-18 | Exit statuses 0 / 2 | HLD §6, LLD §3.7 | `internal/app/app.go::Run` (`ExitSuccess`, `ExitUsage`) | `TestRunSuccess`, `TestRunUsageErrors`, P01–P10 | Verified |
| CONV-22 | `--help`, `--version` take no argument | LLD §3.1 | `internal/cli/cli.go::Parse` | `TestParseRejected/help_with_argument`, `…/version_with_argument` | Verified |
| CONV-23 | Help text verbatim, trailing LF | LLD §3.5 | `internal/greeting/greeting.go::HelpText` | `TestHelpTextLayout`, P07 | Verified |
| CONV-24 | No crash or stack trace | HLD D4 | all packages | as NFR-03 | Verified |
| CONV-25 | Broken pipe ends silently | HLD §6, LLD §5 | `internal/app/app.go::Run`, `internal/output/output.go` | `TestRunOutputFailureIsSilent`, `TestWriteReportsFailure` | Verified |
| CONV-31 | Memory limits for native builds | HLD §7, MTP §4.4 | `internal/output/output.go::WriteRepeated` (bounded buffer) | MTP §4.4 measurement (whole-process peak RSS, no baseline subtracted) | Verified |
| CONV-26 | Stdin not read | HLD §3.1 | `cmd/dp-hello/main.go::main` (no stdin access) | Inspection | Verified |

## 3. Backward traceability (test → requirement)

| Test | Type | Requirement(s) verified |
| :-- | :-- | :-- |
| P01 | black-box, public | FR-01, CONV-03 |
| P02, P03 | black-box, public | FR-02 |
| P04 | black-box, public | FR-02, CONV-08 |
| P05 | black-box, public | FR-04 |
| P06 | black-box, public | FR-04, CONV-07 |
| P07 | black-box, public | FR-07, CONV-23 |
| P08 | black-box, public | FR-08, CONV-02 |
| P09 | black-box, public | FR-10, FR-11, CONV-13, CONV-14 |
| P10 | black-box, public | FR-09, FR-11, CONV-13, CONV-14 |
| `TestParseAccepted` | unit | FR-02, FR-06, FR-07, FR-08, CONV-07 … CONV-12 |
| `TestParseRejected` | unit | FR-10, CONV-12, CONV-13, CONV-22 |
| `TestParseCountAccepted`, `TestParseCountRejected` | unit | FR-04, FR-05, NFR-03 |
| `TestValidateDefaults`, `TestValidateValues` | unit | FR-01, FR-02, FR-04 |
| `TestValidateRejected` | unit | FR-03, FR-05, FR-09, CONV-15 |
| `TestLine` | unit | FR-01, FR-02, NFR-05 |
| `TestHelpTextLayout`, `TestVersionLine` | unit | FR-07, FR-08, CONV-02, CONV-23 |
| `TestReportFormatsBothLines`, `TestMessages` | unit | CONV-01, CONV-13 … CONV-15 |
| `TestWrite`, `TestWriteRepeated`, `TestWriteReportsFailure` | unit | FR-04, NFR-02, CONV-25 |
| `TestRunSuccess` | unit (end to end) | FR-01 … FR-09, CONV-05, CONV-12, CONV-18 |
| `TestRunUsageErrors` | unit (end to end) | FR-03, FR-05, FR-09 … FR-11, CONV-13, CONV-14 |
| `TestRunOutputFailureIsSilent` | unit (end to end) | CONV-25, NFR-03 |

## 4. Coverage summary

| Category | Total | Designed | Implemented | Verified |
| :-- | :-- | :-- | :-- | :-- |
| User requirements (UR) | 12 | 12 | 12 | 12 |
| Functional (FR) | 11 | 11 | 11 | 11 |
| Non-functional (NFR) | 7 | 7 | 7 | 7 |
| Conventions (CONV, applicable) | 21 | 21 | 21 | 21 |

## 5. Gaps and justifications

No applicable requirement is unverified. Not applicable to `dp-hello`: CONV-04 (it copies no
input), CONV-16 and CONV-17 (no files or operands), and CONV-19 … CONV-21 (no paths).
CONV-27 … CONV-30 describe the black-box harness and are satisfied by running P01–P10 through
`shared/tools/run-tests.py`. NFR-02 and NFR-07 are verified by measurement and inspection, as
SRS §7 prescribes.
