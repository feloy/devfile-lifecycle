```mermaid
graph TB
my-container["container: my-container<br/>image: my-image"]
my-build["command: my-build"]
my-container-my-run-expose["Expose ports<br/>http: 8080"]
my-container-stop["Stop containers"]
my-run["command: my-run"]
post-start-1["Post Start<br/>command: post-start-1"]
post-start-2["Post Start<br/>command: post-start-2"]
start["start"]
sync-all-my-container["Sync All Sources"]
sync-modified-my-container["Sync Modified Sources"]
start -->|"dev"| my-container
my-container -->|"container running"| sync-all-my-container
sync-all-my-container -->|"sources synced"| post-start-1
post-start-1 -->|"post-start-1 done"| post-start-2
post-start-2 -->|"post-start-2 done"| my-build
my-build -->|"build done, with run"| my-run
my-run -->|"command running"| my-container-my-run-expose
my-container-my-run-expose -->|"User quits"| my-container-stop
my-container-my-run-expose -->|"source changed"| sync-modified-my-container
sync-modified-my-container -->|"source synced"| my-build
my-container-my-run-expose -->|"devfile changed"| my-container
```
