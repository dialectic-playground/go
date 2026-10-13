# Low-Level Design (LLD) — Standard Terminal Hello World

| Field | Value |
| :-- | :-- |
| Project | `00-hello-world` |
| Tools | dp-hello |
| Implementation language | go (Go 1.27, standard library only) |
| Author | Brandon Perkins |
| Version / date | 1.0 — 2026-10-13 |
| Status | Approved |
| Inputs | `hld.md`, `specs/SRS.md`, `shared/CLI-CONVENTIONS.md` |

## 1. Source layout

Module `github.com/dialectic-playground/go/00-hello-world` (`go.mod` in the project
directory). Paths below are relative to `00-hello-world/`.

```text
go.mod                                  module path, go 1.27, no dependencies
src/
├── cmd/dp-hello/main.go                process boundary: os.Args, os.Stdout, os.Stderr, os.Exit
└── internal/
    ├── app/app.go, app_test.go         Run: orchestration and exit statuses
    ├── cli/cli.go, cli_test.go         Parse: arguments → Invocation (CONV-07 … CONV-12)
    ├── validate/validate.go, …_test.go Validate, ParseCount: Invocation → Config
    ├── greeting/greeting.go, …_test.go Line, HelpText, VersionLine (pure)
    ├── diagnostics/diagnostics.go, …   UsageError and the CONV-13 message constructors
    └── output/output.go, output_test.go Write, WriteRepeated
bin/dp-hello[.exe]                      build output (BIN_DIR), not committed
```

## 2. Data types

### 2.1 `cli.Invocation`
| Field | Type | Default | Constraint / invariant | SRS |
| :-- | :-- | :-- | :-- | :-- |
| `Action` | `cli.Action` (`Greet`, `Help`, `Version`) | `Greet` | `Help`/`Version` only when that option was the first one reached | FR-07, FR-08, CONV-12 |
| `Name` | `*cli.Value` | `nil` (not given) | Last occurrence of `-n`/`--name` | FR-02, FR-06 |
| `Count` | `*cli.Value` | `nil` (not given) | Last occurrence of `-c`/`--count` | FR-04, FR-06 |
| `Operands` | `[]string` | empty | In command-line order, including every token after `--` | FR-09, CONV-09, CONV-10 |

`cli.Value` holds `Text` (the raw value) and `Option` (the spelling, `-n` or `--name`).

### 2.2 `validate.Config`
| Field | Type | Default | Constraint / invariant | SRS |
| :-- | :-- | :-- | :-- | :-- |
| `Name` | `string` | `"World"` | Non-empty | FR-01 … FR-03 |
| `Count` | `int` | `1` | `1 ≤ Count ≤ 1000` | FR-04, FR-05 |

### 2.3 Error type
`*diagnostics.UsageError{Message}`; `Report()` returns `dp-hello: <Message>\n` +
`Try 'dp-hello --help' for more information.\n`. Exit status 2 for every variant.

| Constructor | Message (after `dp-hello: `) |
| :-- | :-- |
| `InvalidShortOption(x)` | `invalid option -- 'x'` |
| `UnrecognizedOption(token)` | `unrecognized option 'token'` |
| `ShortOptionRequiresArgument(x)` | `option requires an argument -- 'x'` |
| `LongOptionRequiresArgument(name)` | `option '--name' requires an argument` |
| `LongOptionDisallowsArgument(name)` | `option '--name' doesn't allow an argument` |
| `InvalidArgument(value, option)` | `invalid argument 'value' for 'option'` |
| `ExtraOperand(operand)` | `extra operand 'operand'` |

A failed write to stdout is not an error type: `app.Run` returns status 1 without a message.

## 3. Module and function specifications

### 3.1 `cli.Parse`
- **Signature:** `func Parse(args []string) (Invocation, *diagnostics.UsageError)`
- **Responsibility:** syntactic scan of the arguments.
- **Preconditions:** `args` excludes the program name.
- **Postconditions:** on success, an `Invocation` as in §2.1; on error, a zero `Invocation`.
- **Errors:** every CONV-13 option-syntax row (unknown short/long, missing argument,
  `--help=`/`--version=`).
- **Algorithm:** for each token, in order:
  1. `--`: append all remaining tokens to `Operands`; stop.
  2. `--NAME[=VALUE]`: split at the first `=`; look NAME up exactly (no abbreviations). Unknown
     → `UnrecognizedOption(token)`. `help`/`version`: `=` present → `LongOptionDisallowsArgument`;
     otherwise set `Action` and stop. `name`/`count`: value is the inline text, else the next
     token, else `LongOptionRequiresArgument`. Store with spelling `--NAME`.
  3. `-X…` (not `-`): `parseShortBundle`.
  4. Otherwise: append to `Operands`.
- **Complexity:** O(total argument length) time, O(number of operands) memory.
- **Satisfies:** FR-02, FR-06, FR-07, FR-08, FR-10, CONV-07 … CONV-12.
- **Unit tests:** `TestParseAccepted`, `TestParseRejected` (`cli_test.go`).

### 3.2 `(*Invocation).parseShortBundle` (private)
- **Signature:** `func (inv *Invocation) parseShortBundle(token string, rest []string) (consumedNext bool, err *diagnostics.UsageError)`
- **Algorithm:** decode the first letter after `-` as one UTF-8 rune (so `-ë` is reported as
  `'ë'`). Unknown → `InvalidShortOption`. Both short options take a value, so the value is the
  rest of the token (`-c2`, `-nc3` → NAME `c3`) or, if that is empty, `rest[0]`
  (`consumedNext = true`); with no next token → `ShortOptionRequiresArgument`.
- **Satisfies:** CONV-07. **Unit tests:** `TestParseAccepted/bundled_values`,
  `TestParseAccepted/bundle_consumes_rest_of_token`, `TestParseRejected/unknown_short_option`.

### 3.3 `validate.Validate`
- **Signature:** `func Validate(inv cli.Invocation) (Config, *diagnostics.UsageError)`
- **Algorithm:** (1) any operand → `ExtraOperand(Operands[0])`; (2) NAME given and empty →
  `InvalidArgument("", spelling)`; (3) COUNT given and `ParseCount` fails →
  `InvalidArgument(text, spelling)`; (4) return the Config with defaults filled in.
- **Satisfies:** FR-03, FR-05, FR-06, FR-09. **Unit tests:** `TestValidateDefaults`,
  `TestValidateValues`, `TestValidateRejected`.

### 3.4 `validate.ParseCount`
- **Signature:** `func ParseCount(text string) (count int, ok bool)`
- **Algorithm:** reject empty text; for each byte, reject anything outside `'0'…'9'`; accumulate
  `count = count*10 + digit` only while `count ≤ 1000`, so long inputs cannot overflow; accept
  when `1 ≤ count ≤ 1000`. Leading zeros are allowed (`007` → 7).
- **Complexity:** O(len(text)) time, O(1) memory.
- **Satisfies:** FR-04, FR-05. **Unit tests:** `TestParseCountAccepted`, `TestParseCountRejected`.

### 3.5 `greeting.Line`, `greeting.HelpText`, `greeting.VersionLine`
- **Signature:** `func Line(name string) string`; two string constants.
- **Responsibility:** pure text construction: `"Hello, " + name + "!\n"`; the help text exactly
  as in the user manual; `dp-hello (dialectic-playground) 1.0.0\n`.
- **Satisfies:** FR-01, FR-02, FR-07, FR-08, NFR-05. **Unit tests:** `TestLine`,
  `TestHelpTextLayout`, `TestVersionLine`.

### 3.6 `output.Write`, `output.WriteRepeated`
- **Signatures:** `func Write(w io.Writer, text string) error`;
  `func WriteRepeated(w io.Writer, line string, count int) error`
- **Algorithm:** `Write` uses `io.WriteString`. `WriteRepeated` writes `line` `count` times into a
  `bufio.Writer` (4 KiB) and flushes; the first error is returned.
- **Satisfies:** FR-04, NFR-02, CONV-03, CONV-25. **Unit tests:** `TestWrite`,
  `TestWriteRepeated`, `TestWriteReportsFailure`.

### 3.7 `app.Run`
- **Signature:** `func Run(args []string, stdout, stderr io.Writer) int`
- **Algorithm:** `cli.Parse` → on error `reportUsage` (write `Report()` to stderr, return 2).
  `Help`/`Version` → write the constant, return 0. Otherwise `validate.Validate` → on error
  `reportUsage`; else `output.WriteRepeated(stdout, greeting.Line(name), count)`. A write error
  returns `ExitOutputFailure` (1) with no message; success returns 0.
- **Satisfies:** FR-01 … FR-11, SRS §3.7, CONV-14, CONV-24, CONV-25.
- **Unit tests:** `TestRunSuccess`, `TestRunUsageErrors`, `TestRunOutputFailureIsSilent`.

### 3.8 `main`
`os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))`. Covered by P01–P10.

## 4. Algorithms

The only non-trivial algorithm is the option scan (§3.1–§3.2), a single left-to-right pass in
the style of GNU `getopt_long` with permutation: operands may be interleaved with options, the
last value wins, and the first `--help`/`--version` ends the scan. Value checks are deferred to
`Validate` (HLD D2), so `dp-hello -c 0 -c 5` succeeds and `dp-hello -c 0 --help` prints help, as
the curriculum's reference model does. A table of four options with linear lookup is the
simplest adequate structure (Principle 12).

## 5. Platform-specific details

| Concern | Linux | macOS | Windows | Abstraction |
| :-- | :-- | :-- | :-- | :-- |
| Argument encoding | Raw bytes (UTF-8) in `os.Args` | Same | Go runtime converts the UTF-16 command line to UTF-8 `os.Args` | Go runtime; `main` passes `os.Args[1:]` |
| Newline translation | None | None | None: `os.File` has no text mode | Not needed |
| Console output | Bytes | Bytes | UTF-8 converted to UTF-16 by Go for consoles; raw bytes for pipes/files | Go runtime |
| Broken pipe | Write on fd 1 raises SIGPIPE; the runtime exits silently | Same | Write error → `Run` returns 1 silently | `app.Run` |
| Executable name | `bin/dp-hello` | `bin/dp-hello` | `bin/dp-hello.exe` (`$(EXE)`) | Makefile |

## 6. Resource management

No files, sockets, or child processes are opened. Memory: the argument slice, one greeting
line, and a 4 KiB `bufio.Writer`; output size never affects memory. The buffer is flushed
before `Run` returns, and `os.Exit` is called only after that.

## 7. Sequence of one invocation

Public test P06, `dp-hello -c2 -n Ken`:
1. `main` → `app.Run(["-c2", "-n", "Ken"], os.Stdout, os.Stderr)`.
2. `cli.Parse`: `-c2` → `parseShortBundle` stores Count `{"2", "-c"}`; `-n` has no attached text,
   so it takes `Ken` (`consumedNext`), storing Name `{"Ken", "-n"}`. Result: `Action = Greet`.
3. `validate.Validate`: no operands; NAME non-empty; `ParseCount("2")` → 2. Config `{Ken, 2}`.
4. `output.WriteRepeated(stdout, "Hello, Ken!\n", 2)` writes 24 bytes and flushes.
5. `Run` returns 0; `main` exits 0. stderr is empty.

## 8. Design-to-code checklist

| HLD component | LLD section | Source file(s) | Unit-test file(s) |
| :-- | :-- | :-- | :-- |
| `cli` | §3.1, §3.2 | `src/internal/cli/cli.go` | `src/internal/cli/cli_test.go` |
| `validate` | §3.3, §3.4 | `src/internal/validate/validate.go` | `src/internal/validate/validate_test.go` |
| `greeting` | §3.5 | `src/internal/greeting/greeting.go` | `src/internal/greeting/greeting_test.go` |
| `output` | §3.6 | `src/internal/output/output.go` | `src/internal/output/output_test.go` |
| `diagnostics` | §2.3 | `src/internal/diagnostics/diagnostics.go` | `src/internal/diagnostics/diagnostics_test.go` |
| `app` | §3.7 | `src/internal/app/app.go` | `src/internal/app/app_test.go` |
| `cmd/dp-hello` | §3.8 | `src/cmd/dp-hello/main.go` | public black-box P01–P10 |
