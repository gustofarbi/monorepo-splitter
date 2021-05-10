# Repository splitting
## Usage
1. install task: https://github.com/go-task/task
2. run `task install` and open a new shell
3. create config file `<some name>.yaml`
4. run `splitter -c <path to your config>.yaml`
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
