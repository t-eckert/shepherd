# Shepherd CLI

The Shepherd CLI is used to move and modify open GitHub issues across repositories.

## Prerequisites

- [Go](https://golang.org)
  ```bash
  brew install go
  ```
- [GitHub Command Line](https://cli.github.com/)
  ```bash
  brew install gh
  ```
  
## Installation

```bash
go get github.com/t-eckert/shepherd 
```

## Usage

Shepherd takes two required parameters, an origin repository and a destination repository. All open issues will be migrated from the origin repository to the destination repository by the tool. Both repositories passed in should be formatted as `OWNER/REPO`. For example, `hashicorp/consul-k8s`.

An optional parameter can be passed in with the flag `-p` or `--modify-prepend`. This tells Shepherd to modify the title of each open issue in the origin repository. If the issue title has text prepended and terminated by a `:` -- e.g. `consul:Issue about something` -- the text of the prepend will be replaced with the value passed for the flag. If the issue title does not have existing prepended text, the value passed in will be added. There is no need to pass in the `:` when setting the flag.

### Examples

Migrate all open issues from `hashicorp/consul-helm` to `hashicorp/consul-k8s` and prepend `helm:` to each issue on migration.

```bash
shepherd hashicorp/consul-helm hashicorp/consul-k8s -p helm
```

