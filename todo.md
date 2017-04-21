### important performance improvements

- compare speed against java implementation
    - seems generally faster in Go, but...
    - ... there seems to be a bug. Takes extremely long
     to process ~/Pictures, overheating CPU (interrupted)
- test run on home pc, entire media disk

### polishing

- add size scale shortcut for -minSize
- add -stdin flag to read path list from stdin
- add -0 flag to read null-delimited path list from stdin
- clean up interface, hide implementation details
- document public types and functions
- cleanup all TODOs
- 100% test coverage
- publish documentation
- setup build badges like portping
- post on code review
- post on blog
