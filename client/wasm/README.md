# Web Demo 

This demo uses the devfile-lifecyle library (thanks to WebAssembly) to build a graph from a Devfile content, and display it on a webpage.

Usage:

```shell
(cd ../.. && odo dev --platform podman)
```

Then visit localhost:20001, enter a Devfile content into the textarea. The graph should be displayed on the right.

![Screenshot](./webclient-devfile-lifecycle.png)
