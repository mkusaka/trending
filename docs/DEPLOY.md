# Deploying trending (Go only)

- Prerequisites: Go (matching `go.mod`), GitHub Pages enabled for “GitHub Actions”.

## Local build
- Generate site: `go run main.go`
- Output: `out/` (e.g. `out/javascript/daily/index.xml`)

## Publish to GitHub Pages
- CI workflow: `.github/workflows/daily.yml`
  - Builds with Go and uploads `out/` as Pages artifact
  - Deploys via `actions/deploy-pages`
- How to trigger:
  - Scheduled: daily at 13:00 UTC
  - Manual: GitHub → Actions → “daily trending check” → Run workflow

Node/yarn are no longer required; `gh-pages` CLI has been removed.

