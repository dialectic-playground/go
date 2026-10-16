# Toolchain image for dialectic-playground/go: Go (container-first; see tools/dp-run).
#
# Go 1.27; build and package tool: Go modules.
# Official base image pinned by multi-architecture index digest, plus make and python3 (for the
# curriculum's black-box runner). Dependabot (docker) refreshes the digest; moving to a new
# release line is a deliberate change (keep `make toolchain` and the native CI in step).
FROM docker.io/library/golang:1.27-trixie@sha256:47667934874913a449b96faa7ef4777fa285d9ec29126ee80ed61c45b7cccf0e

# Debian keeps only the current version of each package, so apt versions cannot be pinned
# reproducibly; the base image is pinned by digest instead.
# hadolint ignore=DL3008
RUN apt-get update \
    && apt-get install -y --no-install-recommends make python3 \
    && rm -rf /var/lib/apt/lists/*

# GOPATH in the mounted cache; never download a different Go toolchain.
ENV GOPATH=/dp-cache/go \
    GOTOOLCHAIN=local \
    DP_IN_CONTAINER=1
