gowatch
---

A simple program that executes a command on file changes in the selected directory (current directory by default).

---
To install run:
---
```
go install github.com/dele256/gowatch@latest
```


---
When rerunning the given command, all sub-processes spawned by the command are also killed.

```
gowatch -c "go run ."
```
go run launches your program after it has built it, meaning we are launching the `go` executable with our command. But since we are assigning a PGID before launching the command (the `go` executable) we can kill all sub-processes aswell before rerunning the command again.

Usage of gowatch:
---
  * -c string
        
        command to run e.g. "go run ."
  * -d string
        
        run directory, the directory to run the command in (default "./")
  * -f string
        
        comma separated extensions to watch e.g. "go,cc,c,h" (default="")
  * -w string
        
        watch directory, directory to recursively watch (default "./")