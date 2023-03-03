# Devfile lifecycle

The purpose of this tool is to analyze a [Devfile](https://devfile.io/) and determine its lifecycle.

The first iteration of the tool is focusing on building a graphical representation of the lifecycle of a Devfile.

This is  a work in progress and all Devfile features are not supported yet. Please explore examples below to check which features are supported.

For example, for a Devfile [container-build-run.yaml](./tests/devfiles/container-build-run.yaml) defining a single component and two commands, *build* and *run*, the lifecycle is represented as follows:

```mermaid
graph TB
my-container["container: my-container<br/>image: my-image"]
my-build["command: my-build"]
my-container-my-run-expose["Expose ports<br/>http: 8080"]
my-container-stop["Stop container<br/>container: my-container"]
my-run["command: my-run"]
start["start"]
sync-all-my-container["Sync All Sources"]
sync-modified-my-container["Sync Modified Sources"]
start -->|"dev"| my-container
my-container -->|"container running"| sync-all-my-container
sync-all-my-container -->|"sources synced"| my-build
my-build -->|"build done, with run"| my-run
my-run -->|"command running"| my-container-my-run-expose
my-container-my-run-expose -->|"User quits"| my-container-stop
my-container-my-run-expose -->|"source changed"| sync-modified-my-container
sync-modified-my-container -->|"source synced"| my-build
my-container-my-run-expose -->|"devfile changed"| my-container
```

The output of the command is in text format, using the Mermaid standard. This output can be exploited by Mermaid-compatible tools (GitHub pages, IDEs, inline vizualizers) to obtain a graphical representation.

## Usage

```bash
# From a file on disk
$ go run main.go path/to/my/devfile.yaml > devfile.md

# From a Devfile Registry (don't forget trailing /)
$ go run main.go https://registry.devfile.io/devfiles/nodejs-basic/ > nodejs-basic.md
```

## Examples

Examples are visible in the [tests/graphs directory](./tests/graphs), with corresponding Devfile sources in [tests/devfiles directory](./tests/devfiles).


## Next steps

The next iteration of the tool would be to develop a library (at least in Go) providing a state-machine based on the Devfile, usable by tools working with Devfiles.
