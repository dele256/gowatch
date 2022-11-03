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
  * -g
      
        include .git folder (default=false)
  * -w string
        
        watch directory, directory to recursively watch (default "./")

---
Makefile example to run a make command and to only watch go files:
---
```
gowatch -c "make run" -f go
```

Then inside the Makefile we can have something like this:
```makefile
run:
      go build -o build/example.out main.go
      (cd build && ./example.out)
```
