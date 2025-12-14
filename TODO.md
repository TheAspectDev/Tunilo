## TODO

This document tracks planned improvements and features for the project,
grouped by expected implementation timeframe.

### Short-Term (Improvements)

- **Custom HTML Error pages**
  - Add "no tunnel connected" html page
  - "Internal Server Error" html page

- **Refactor client architecture** 
  - ~~✅ Simplify client responsibilities~~
  - ~~✅ Improve separation of concerns~~ 

- **Add TLS**
  - Add TLS encryption
  - Authenticate the server

- **Build a TUI (Terminal User Interface)** 
  - ~~✅ Add TUI to server~~
  - Add TUI to client
  - Improve usability over raw CLI output 
  - Configuration using CLI

### Long-Term ( Major Changes )

- **Chunking requests**
  - Split large requests into smaller chunks
  - Reduce memory pressure for large payloads