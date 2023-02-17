```mermaid
graph TB
my-container["container: my-container<br/>image: my-image"]
my-build-1["command: my-build-1"]
my-build-2a["command: my-build-2a"]
my-build-2b["command: my-build-2b"]
my-build-3["command: my-build-3"]
my-container-my-run-expose["Expose ports<br/>http: 8080"]
my-run["command: my-run"]
sync-all-my-container["Sync All Sources"]
sync-modified-my-container["Sync Modified Sources"]
my-container -->|"container running"| sync-all-my-container
sync-all-my-container -->|"sources synced"| my-build-1
my-build-1 -->|"my-build-1 done"| my-build-2a
my-build-2a -->|"my-build-2a done"| my-build-2b
my-build-2b -->|"my-build-2 done"| my-build-3
my-build-3 -->|"build done, with run"| my-run
my-run -->|"command running"| my-container-my-run-expose
my-container-my-run-expose -->|"source changed"| sync-modified-my-container
sync-modified-my-container -->|"source synced"| my-build-1
my-container-my-run-expose -->|"devfile changed"| my-container
```
