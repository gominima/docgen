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

| Name | Position | Description | Default |
|---|---|---|---|
| path | 1 | The directory which contains the docgen files | `.` The current directory |
| outputfile | 2 | The name of the output JSON file | `output.json` |

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
    steps:
      - name: Checkout repository
        uses: actions/checkout@v2.4.0
        
      - name: Setup Go environment
        uses: actions/setup-go@v2.1.5
        with:
          go-version: 1.17

      - name: Get Docgen
        run: go get github.com/gominima/docgen 
        
      - name: Run Docgen
        run: go run github.com/gominima/docgen . $(git rev-parse --abbrev-ref HEAD).json
        
      - name: Extract branch name
        shell: bash
        run: echo "##[set-output name=branch;]$(echo ${GITHUB_REF#refs/heads/}.json)"
        id: extract_branch

      - name: Push docs file
        uses: dmnemec/copy_file_to_another_repo_action@main
        env:
          API_TOKEN_GITHUB: ${{ secrets.API_TOKEN_GITHUB }}
        with:
          source_file: ${{ steps.extract_branch.outputs.branch }}
          destination_repo: 'gominima/docs'
          user_email: ''
          user_name: ''
          commit_message: 'chore: docs'
```

Add your email and name in the `user_email: ''` `user_name: ''` sections.

## Format

The docgen format is heavily inspired from the JSDoc format, it however does not completely support the JSDoc specification and is NOT a replacement for JS documentation Generators

### Functions
- `@info` tag to add `Description` in JSON
- `@param {type} [name] optional description` to add a value to the `Parameters` array section
- `@returns {type} optional description` to add `Returns` value in JSON
```go
/**
  @info function description
  @param {string} [example] param description
  @returns {int}
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
  @info structure description
  @property {string} [example] property description
*/
type example struct {
  example string
}
```
