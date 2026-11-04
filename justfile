# AgentStrap — Go CLI. All recipes run from the repository root.
set shell := ["bash", "-euo", "pipefail", "-c"]
set working-directory := "."

default:
    @just --list

# Fast-forward and run the complete check once (build already includes ci).
update:
    git pull --ff-only
    just build

bootstrap:
    @go version
    go mod download

fmt:
    gofmt -w cmd internal

fmt-check:
    @files="$(gofmt -l cmd internal)"; if [[ -n "$files" ]]; then printf 'Unformatted Go files:\n%s\n' "$files" >&2; exit 1; fi

vet:
    go vet ./...

test:
    go test -count=1 ./...

# Keep the inherited AI-template structure valid.
validate:
    python3 scripts/validate_template.py

precommit: fmt-check vet test validate

ci: precommit
    go test -race ./...

build: ci
    mkdir -p dist
    go build -trimpath -ldflags "-X main.version=$(tr -d '\r\n' < VERSION)" -o dist/agentstrap ./cmd/agentstrap
    @test "$(./dist/agentstrap version)" = "agentstrap $(tr -d '\r\n' < VERSION)"
    ./dist/agentstrap --help

run *args:
    go run ./cmd/agentstrap {{args}}

install: build
    @destination="${INSTALL_DIR:-$HOME/.local/bin}"; mkdir -p "$destination"; install -m 0755 dist/agentstrap "$destination/agentstrap"; "$destination/agentstrap" version

# Review local validation and open GitHub pull requests. --run executes just build.
next *args:
    bash scripts/next.sh {{args}}

clean:
    rm -rf dist
