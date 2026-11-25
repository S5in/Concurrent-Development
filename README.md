# Concurrent Development – Labs

**Author:** Maksym Redchenko  
**Module:** Concurrent Development, SETU Carlow, 2026/27  
**Licence:** GPL-3.0 (see LICENSE)  
**Language:** Go 1.22 or newer  
**GitHub:** https://github.com/S5in/Concurrent-Development

Each lab is in its own folder and contains:
- the Go code (one program per folder, each with a licence header and comments),
- `README.md` – how to run it, what it shows, results,
- `DOCUMENTATION.md` – design, correctness argument, testing and a code reference.

Each lab is a separate Go module: `cd` into the lab folder and use `go run .` (or `go run ./<subfolder>`).

| Lab | Folder | Topic |
|---|---|---|
| 1 | [lab01-mutex](lab01-mutex) | Mutual exclusion (mutex, atomic) and signalling |
| 2 | [lab02-barrier](lab02-barrier) | Rendezvous and single-use barrier |
| 3 | [lab03-reusable-barrier](lab03-reusable-barrier) | Reusable two-phase barrier (mutex and atomic) |
| 4 | [lab04-semaphore](lab04-semaphore) | Semaphore type (buffered channel), signal and multiplex |
| 5 | [lab05-dining-philosophers](lab05-dining-philosophers) | Dining philosophers: deadlock demo, footman, asymmetric |
| 6 | [lab06-producer-consumer](lab06-producer-consumer) | Producer–consumer: bounded buffer (mutex + two semaphores), wrong-order deadlock |

## Requirements
- Go 1.22 or newer. No external libraries.
