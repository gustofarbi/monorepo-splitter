# Repository splitting
## Download
You can download binary build for MacOS: https://github.com/myposter-de/monorepo-splitter/releases
## Installation (this requires golang)
1. run `make install` and open a new shell
2. create config file `<some name>.yaml`
3. run `splitter -c <path to your config>.yaml`
## Config file syntax
```
version: 0.0.3
root:
  branch: feature/canvas-renderer/split
  path: ~/mp/myposter
  remote: origin
packages:
  prefix: packages
  branch: master
  items:
    - remote: foo
      url: https://github.com/foo-de/canvas.git
      path: some/path/to/repo
    - remote: bar
      url: https://github.com/bar-de/enum.git
actions:
  - validate
  - set-packages-dependencies
  - update-configs
  - write-changes
  - push-changes
  - split-packages
  - reset
```
