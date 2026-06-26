#!/usr/bin/env bash
#
# mayhem/build.sh — build this repo's go-fuzz harnesses as sanitized libFuzzer
# binaries (OSS-Fuzz Go path: go-fuzz-build -libfuzzer + clang link). Runs inside the
# commit image (GO mayhem/Dockerfile) as `mayhem` in /mayhem.
#
# GOROOT/GOPATH/GOMODCACHE are pinned by the Dockerfile ENV (under /opt/toolchains —
# absolute, $HOME-independent), so the module cache survives the PATCH re-run identity.
#
# AIR-GAPPED CONTRACT (SPEC §6.5): the PATCH tier re-runs THIS script OFFLINE.
#   - This FIRST build (online, in CI) populates the module cache under $GOMODCACHE,
#     which doubles as a FILE PROXY at $GOMODCACHE/cache/download.
#   - GOPROXY points at that file proxy FIRST, network LAST: the offline re-run
#     resolves entirely from the cache; the network entries only fill cache-misses on
#     this first online build. -mod=mod lets go-fuzz-build's `go get` of go-fuzz-dep
#     (and the AdaLogics fuzz-headers) update go.mod from the cache.
set -euo pipefail

[ -n "${SOURCE_DATE_EPOCH:-}" ] || unset SOURCE_DATE_EPOCH

: "${SRC:=/mayhem}"
: "${CC:=clang}" ; : "${CXX:=clang++}" ; : "${LIB_FUZZING_ENGINE:=-fsanitize=fuzzer}"
# OSS-Fuzz Go path links with ASan (keep it even if the base default is empty).
: "${SANITIZER_FLAGS=-fsanitize=address}"
# DWARF < 4 (SPEC §6.2 item 10): Go's gc compiler always emits DWARF4 with no
# downgrade flag, but go-fuzz links the harness via clang++, whose C (cgo) shims land
# FIRST in the binary. Forcing those C compilation units to DWARF3 (via CGO_*FLAGS +
# the final clang link) makes the first .debug_info CU DWARF3 — what verify-repo's
# `readelf` check reads.
: "${GO_DEBUG_FLAGS:=-g -gdwarf-3}"
: "${MAYHEM_JOBS:=$(nproc)}"
export CC CXX LIB_FUZZING_ENGINE SANITIZER_FLAGS GO_DEBUG_FLAGS MAYHEM_JOBS

export CGO_CFLAGS="${CGO_CFLAGS:+$CGO_CFLAGS }$GO_DEBUG_FLAGS"
export CGO_CXXFLAGS="${CGO_CXXFLAGS:+$CGO_CXXFLAGS }$GO_DEBUG_FLAGS"

# Resolve modules offline-first from the in-image cache; network only as a fallback.
# $(go env GOMODCACHE) reads the pinned ENV, so it is correct under ANY $HOME.
export GOFLAGS="${GOFLAGS:--mod=mod}"
export GOPROXY="${GOPROXY:-file://$(go env GOMODCACHE)/cache/download,https://proxy.golang.org,direct}"

cd "$SRC"
go version

# go-fuzz-build needs go-fuzz-dep on the module graph; the harnesses need the AdaLogics
# fuzz-headers. With -mod=mod + the file-proxy GOPROXY these resolve from the cache
# offline (no-op once present). Done here (not via a committed go.mod edit) to keep the
# mayhem layer purely additive vs upstream.
go get github.com/dvyukov/go-fuzz/go-fuzz-dep
go get github.com/AdaLogics/go-fuzz-headers

build_target() {
  local dir="$1" out="$2"
  echo "=== building $out (go-fuzz-build -libfuzzer) ==="
  ( cd "$SRC/$dir" && go-fuzz-build -libfuzzer -o "$SRC/mayhem-build/$(basename "$out").a" )
  $CXX $SANITIZER_FLAGS $GO_DEBUG_FLAGS $LIB_FUZZING_ENGINE "$SRC/mayhem-build/$(basename "$out").a" -o "$out"
  echo "built $out"
}

mkdir -p "$SRC/mayhem-build"
build_target "mayhem/fuzz_virtualpaper_process" "/mayhem/fuzz_virtualpaper_process"
build_target "mayhem/fuzz_virtualpaper_delete"  "/mayhem/fuzz_virtualpaper_delete"

echo "build.sh complete"
