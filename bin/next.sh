#!/usr/bin/env bash
set -euo pipefail

if [[ "${1:-}" != "" && "${1:-}" != "--run" ]] || [[ $# -gt 1 ]]; then
  printf 'Usage: just next [--run]\n' >&2
  exit 2
fi

printf 'REPOSITORY: r14r/agentstrap\n'
printf 'BRANCH: %s\n' "$(git branch --show-current)"
if [[ -n "$(git status --porcelain)" ]]; then
  printf 'WORKTREE: has changes\n'
else
  printf 'WORKTREE: clean\n'
fi

printf '\nOPEN PULL REQUESTS:\n'
if command -v gh >/dev/null 2>&1; then
  gh pr list --repo r14r/agentstrap --state open --limit 30 || true
else
  printf '  gh CLI not installed; see GitHub for open PRs.\n'
fi

printf '\nNEXT STEPS:\n'
printf '  git diff --check\n'
printf '  just build   # includes formatting, vet, tests, template checks, and race tests\n'
printf '  gh pr checks # when current branch has an open PR\n'
printf '  Review and merge only after explicit approval.\n'

if [[ "${1:-}" == "--run" ]]; then
  printf '\nEXECUTING LOCAL CHECKS:\n'
  git diff --check
  just build
  if command -v gh >/dev/null 2>&1 && gh pr view --repo r14r/agentstrap >/dev/null 2>&1; then
    gh pr checks --repo r14r/agentstrap
  fi
fi
