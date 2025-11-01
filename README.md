# File Utility CLI Tool (Go)

Small CLI tool written in Go to:
- separate a pipe-delimited list of files into per-year files, and
- compare two pipe-delimited lists (internal vs external) producing difference reports.

Files in this workspace:
- [main.go](main.go) — main program and implementation
- [README.md](README.md)

Key symbols in the code:
- [`main.main`](main.go) — program entry and interactive menu
- [`main.showBanner`](main.go) — prints header/banner
- [`main.separateByYear`](main.go) — separates input file into files named <year>.txt
- [`main.compare`](main.go) — compares internal and external lists and writes reports

Build & run
1. Run without building:

   ```sh
   go run main.go
   ```