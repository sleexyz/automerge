# Automerge Tool - Design Plan

## Overview
Create a tool that polls GitHub API for status checks on the current branch and exits with appropriate codes based on check results.

## Design Options

### Option 1: Pure `gh` wrapper (shell script)
**Pros:**
- Minimal implementation
- Leverages existing `gh` authentication
- No additional dependencies

**Cons:**
- Limited error handling and parsing
- Shell scripting complexity for JSON parsing
- Less maintainable for complex logic

### Option 2: Go tool using `gh` for authentication
**Pros:**
- Better JSON parsing and error handling
- More robust implementation
- Consistent with existing `j` tool patterns
- Can shell out to `gh` for authentication/API calls

**Cons:**
- Slightly more complex than pure shell
- Still depends on `gh` binary

### Option 3: Go tool with direct GitHub API integration
**Pros:**
- Full control over API interactions
- No external dependencies
- Most flexible for future enhancements

**Cons:**
- Need to handle authentication ourselves
- More complex initial implementation
- Reinventing wheel for GitHub auth

## Recommended Approach: Option 2 (Go tool using `gh`)

**Rationale:**
- Balances simplicity with maintainability
- Follows the pattern established by the `j` tool
- `gh` handles authentication complexity
- Go provides excellent JSON parsing and error handling
- Easy to extend in the future

## Implementation Plan

1. **Core functionality:**
   - Use `gh api` to fetch status checks for current branch
   - Parse JSON response to check status of all checks
   - Exit 0 if all pass, exit 1 if any fail
   - Print clear failure messages

2. **Error handling:**
   - Handle cases where branch has no checks
   - Handle API rate limiting
   - Handle network errors
   - Validate we're in a git repo

3. **Structure:**
   - Single main.go file initially (can expand later)
   - Follow `j` tool's Nix flake pattern
   - Include proper version management

4. **Dependencies:**
   - Require `gh` binary in PATH
   - Use Go's standard library for JSON parsing
   - Use `os/exec` to shell out to `gh`