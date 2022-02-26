# Gominima/Docgen

Generates documentation from code comments automatically

## Installation

```bash
go get github.com/gominima/docgen
```

## Usage

```bash
go run github.com/gominima/docgen
```

## CLI Arguments

| Name       | Position | Description                                   | Default                   |
| ---------- | -------- | --------------------------------------------- | ------------------------- |
| path       | 1        | The directory which contains the docgen files | `.` The current directory |
| outputfile | 2        | The name of the output JSON file              | `output.json`             |

Gominima/Docgen does not utilises flags because there are only 2 arguments and if you skip the first one `.` and write the second one, even as, `-O main.json` thats still a character more than just writing `. main.json`

## Workflow

You can also automate the documentation generation by using a workflow which pushes to a seprate repository

Here is an example of the workflow we use in minima repository:

```yml
name: Generate Documentation
on:
  push:
    branches:
      - "main"
      - "stable"

jobs:
  docs:
    name: Generate Docs
    runs-on: ubuntu-latest
    if: github.repository_owner == 'gominima'
    outputs:
      REPO_NAME: ${{ steps.env.outputs.REPO_NAME }}
      BRANCH_NAME: ${{ steps.env.outputs.BRANCH_NAME }}
      BRANCH_OR_TAG: ${{ steps.env.outputs.BRANCH_OR_TAG }}
      SHA: ${{ steps.env.outputs.SHA }}
    steps:
      - name: Checkout repository
        uses: actions/checkout@v2.4.0
        
      - name: Setup Go environment
        uses: actions/setup-go@v2.1.5
        with:
          go-version: 1.17

      - name: Get docgen
        run: go get github.com/gominima/docgen 
        
      - name: Run docgen
        run: go run github.com/gominima/docgen . $(basename `git rev-parse --show-toplevel`)-$(git rev-parse --abbrev-ref HEAD).json
        
      - name: Checkout docs repository
        uses: actions/checkout@v2
        with:
          repository: 'gominima/docs'
          token: ${{ secrets.API_TOKEN_GITHUB }}
          path: 'out'

      - name: Commit and push
        run: |
          mv $(basename `git rev-parse --show-toplevel`)-$(git rev-parse --abbrev-ref HEAD).json out
          cd out
          git config user.name github-actions[bot]
          git config user.email 41898282+github-actions[bot]@users.noreply.github.com
          git add .
          git commit -m "Docs build for ${GITHUB_REF_TYPE} ${GITHUB_REF_NAME}: ${GITHUB_SHA}" || true
          git push
```

## Format

The docgen format is heavily inspired from the JSDoc format, it however does not completely support the JSDoc specification and is NOT a replacement for JavaScript documentation generators

### Functions

- `@info` tag to add `Description` in JSON
- `@param {type} [name] optional description` to add a value to the `Parameters` array section
- `@returns {type} optional description` to add `Returns` value in JSON

```go
/**
 * @info function description
 * @param {string} [example] param description
 * @returns {int}
*/
func someFunc(example string) int {
	return 1;
}
```

### Structures

- `@info` tag to add `Description` in JSON
- `@property {type} [name] optional description` to add a value to the `Properties` array section
- Functions on Structures are automatically added to `Functions` array section inside the Structure if they have comments

```go
/**
 * @info structure description
 * @property {string} [example] property description
*/
type example struct {
  example string
}
```
