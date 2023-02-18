```mermaid
graph TB
my-container["container: my-container<br/>image: my-image"]
my-build["command: my-build"]
my-container-my-debug-expose["Expose ports<br/>http: 8080<br/>debug: 5858"]
my-container-my-run-expose["Expose ports<br/>http: 8080"]
my-container-stop["Stop container<br/>container: my-container"]
my-debug["command: my-debug"]
my-run["command: my-run"]
sync-all-my-container["Sync All Sources"]
sync-modified-my-container["Sync Modified Sources"]
my-container -->|"container running"| sync-all-my-container
sync-all-my-container -->|"sources synced"| my-build
my-build -->|"build done, with run"| my-run
my-run -->|"command running"| my-container-my-run-expose
my-container-my-run-expose -->|"User quits"| my-container-stop
my-container-my-run-expose -->|"source changed"| sync-modified-my-container
sync-modified-my-container -->|"source synced"| my-build
my-container-my-run-expose -->|"devfile changed"| my-container
my-build -->|"build done, with debug"| my-debug
my-debug -->|"command running"| my-container-my-debug-expose
my-container-my-debug-expose -->|"User quits"| my-container-stop
my-container-my-debug-expose -->|"source changed"| sync-modified-my-container
my-container-my-debug-expose -->|"devfile changed"| my-container
```
