remember to `git add` before running any nix build / install commands.

Bump the patch number up every update so we can keep track of what version is installed at the current time.

## Build & Test

- `go build -o automerge .` - Build binary
- `./automerge` - Test locally built binary  
- `nix build` - Build via Nix
- `nix profile install .` - Install via nix profile

## Usage

- `automerge` - Poll status checks on current branch
- `automerge --help` - Show help
- `automerge --version` - Show version

## Requirements

- Must be in a git repository
- `gh` CLI must be installed and authenticated
- Repository must be on GitHub

## How it works

1. Checks if in git repo and `gh` is available
2. Gets current branch and repo info via `gh`
3. Polls GitHub API for status checks and check runs
4. Exits 0 if all pass, exits 1 if any fail
5. Prints clear failure messages for failed checks

## Testing

The tool includes a GitHub Actions workflow that runs:
- Go build and test
- Nix build verification  
- Linting with golangci-lint

### Manual Testing

To test the automerge tool:

1. **Success case**: Create a PR with passing CI
   ```bash
   # Make a simple change, commit, push
   git checkout -b test-branch
   echo "test" >> CLAUDE.md
   git add . && git commit -m "test change" && git push -u origin test-branch
   gh pr create --title "Test PR" --body "Testing automerge"
   ./automerge  # Should exit 0 when all checks pass
   ```

2. **Failure case**: Create a PR with failing CI
   ```bash
   # Add unused import to trigger linting error
   # Edit main.go to add: import "net/http" (unused)
   git add . && git commit -m "trigger CI failure" && git push
   ./automerge  # Should exit 1 and show failure details
   ```

Expected behavior:
- Success: Polls with dots, then "✅ All status checks passed!" and exit 0
- Failure: Shows "❌ [check]: [reason]" and exit 1