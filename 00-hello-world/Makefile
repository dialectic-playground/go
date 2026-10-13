# 00-hello-world: Standard Terminal Hello World
#
# Project Makefile template. It extends the universal build contract in
# shared/templates/BASE-MAKEFILE, which defines the standard targets:
#
#   make all | install | uninstall | test | check | lint | format | clean | distclean
#   make test-hidden | grade | help
#
# A language implementation repository copies this file next to its sources, sets BIN_DIR,
# fills in each DP_* hook with its toolchain's commands, and points CURRICULUM_ROOT at its
# checkout of sdlc-curriculum (for example a git submodule). Hooks that are left empty make
# their target fail with an explanation. `make help` shows which hooks are defined.
#
# Container-first: goals run in the toolchain container built from the repository's
# Containerfile (podman, then docker, then local; override with DP_RUNTIME). The image needs
# the language toolchain, make, and python3 >= 3.11. Use DP_RUNTIME=local for native builds.

PROJECT := 00-hello-world
TOOLS   := dp-hello

# Path to the sdlc-curriculum checkout (the specs, test suites, and runner live there).
CURRICULUM_ROOT ?= ../sdlc-curriculum

# Go (modules): go.mod lives here (no third-party dependencies, so there is no go.sum);
# the command is src/cmd/dp-hello and its units are packages under src/internal/. Build and
# module caches live outside the project (GOCACHE, GOPATH).
GO      ?= go
GOFMT   ?= gofmt
PACKAGE := ./src/cmd/$(TOOLS)

# Directory that contains the built executables after `make all`.
BIN_DIR := bin

# Language-specific commands. Use $(TOOLS), $(EXE), $(PREFIX), and $(BINDIR) as needed.
DP_BUILD     = $(GO) build -trimpath -o "$(BIN_DIR)/$(TOOLS)$(EXE)" $(PACKAGE)
define DP_INSTALL
mkdir -p "$(BINDIR)"
cp "$(BIN_DIR)/$(TOOLS)$(EXE)" "$(BINDIR)/$(TOOLS)$(EXE)"
endef
DP_UNINSTALL = rm -f "$(BINDIR)/$(TOOLS)$(EXE)"
DP_UNIT_TEST = $(GO) test ./...
# go vet, then fail if gofmt would change any file.
define DP_LINT
$(GO) vet ./...
@unformatted="$$($(GOFMT) -l src)"; [ -z "$$unformatted" ] || { printf 'gofmt: run make format; unformatted files:\n%s\n' "$$unformatted" >&2; exit 1; }
endef
DP_FORMAT    = $(GOFMT) -w src
DP_CLEAN     = rm -rf "$(BIN_DIR)"
# Nothing else lives in the project directory: Go keeps its caches in GOCACHE and GOPATH.
DP_DISTCLEAN = rm -rf "$(BIN_DIR)"

include $(CURRICULUM_ROOT)/shared/templates/BASE-MAKEFILE
