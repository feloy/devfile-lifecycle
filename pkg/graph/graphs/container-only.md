```mermaid
graph TB
my-container["container: my-container<br/>image: my-image"]
sync-all-my-container["Sync All Sources"]
my-container -->|"container running"| sync-all-my-container
```
