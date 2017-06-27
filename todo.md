### stress test

- compare speed against java implementation
- test run on home pc, entire media disk

### polishing

- normalize paths in finder, otherwise dupfinder tmp tmp/ will have not real dups
- get the pools more directly in dupTracker.getDupGroups,
  without the ugly delete
- cleanup all TODOs
- 100% test coverage
- add size scale shortcut for -minSize
- document public types and functions
- publish documentation
- setup build badges like portping
- post on code review
- post on blog
