# AgentStrap — Go CLI. Run recipes from repository root.
set shell := ["bash", "-euo", "pipefail", "-c"]
set working-directory := "."

default:
    @just --list

# Display all available commands.
help:
    @just --list

version:
    @cat VERSION

# Update code with fast-forward only, then run the full CI/build chain once.
update:
    git pull --ff-only
    just build

# Prepare Go modules without changing system-wide configuration.
bootstrap:
    @command -v go >/dev/null 2>&1 || { echo "Go toolchain required" >&2; exit 1; }
    @go version
    go mod download

setup: bootstrap

# Format Go source files.
fmt:
    gofmt -w cmd internal

# Check formatting without modifying the working tree.
fmt-check:
    @files="$(gofmt -l cmd internal)"; if [[ -n "$files" ]]; then printf 'Unformatted Go files:\n%s\n' "$files" >&2; exit 1; fi

vet:
    go vet ./...

lint: fmt-check vet

test:
    go test -count=1 ./...

# Preserve validation of the inherited AI-project skills and examples.
validate:
    python3 scripts/validate_template.py

# Format check + static analysis + unit tests + project validation.
precommit: bootstrap lint test validate

# Includes precommit; no need to run precommit and test separately.
ci: precommit
    go test -race -count=1 ./...

# Includes CI, then builds and smoke-tests the versioned binary.
build: ci
    mkdir -p dist
    go build -trimpath -ldflags "-X main.version=$(tr -d '\r\n' < VERSION)" -o dist/agentstrap ./cmd/agentstrap
    @test "$(./dist/agentstrap version)" = "agentstrap $(tr -d '\r\n' < VERSION)"
    ./dist/agentstrap --help

# Install to /usr/local/bin (override with INSTALL_DIR for custom locations/CI).
# Automatically use sudo only when the destination cannot be written directly.
install: build
    @destination="${INSTALL_DIR:-/usr/local/bin}"; \
      if [[ -d "$destination" && -w "$destination" ]]; then \
        install -m 0755 dist/agentstrap "$destination/agentstrap"; \
      elif [[ ! -e "$destination" && -w "$(dirname "$destination")" ]]; then \
        mkdir -p "$destination"; \
        install -m 0755 dist/agentstrap "$destination/agentstrap"; \
      else \
        command -v sudo >/dev/null 2>&1 || { echo "Cannot write to $destination and sudo is unavailable" >&2; exit 1; }; \
        sudo install -d -m 0755 "$destination"; \
        sudo install -m 0755 dist/agentstrap "$destination/agentstrap"; \
      fi; \
      actual="$("$destination/agentstrap" version)"; \
      expected="agentstrap $(tr -d '\r\n' < VERSION)"; \
      if [[ "$actual" != "$expected" ]]; then \
        printf 'Version mismatch: installed %s, expected %s\n' "$actual" "$expected" >&2; exit 1; \
      fi; \
      printf 'Installed %s (%s)\n' "$destination/agentstrap" "$actual"

run *args:
    go run ./cmd/agentstrap {{args}}

# Inspect branch/worktree, all open PRs, and suggested acceptance checks.
# --run executes the checks, but never merges or deploys automatically.
next *args:
    bash scripts/next.sh {{args}}

clean:
    rm -rf dist
