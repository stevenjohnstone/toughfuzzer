# Tough Fuzzer

Obstacle course for go-fuzz to help explain its inner workings. Using a series of fuzzers which have to follow certain code
paths to reach a ```panic```, the different features of go-fuzz are showcased.

## Int and String Literals

Go-fuzz [extracts int and string literals](https://github.com/dvyukov/go-fuzz/blob/6a8e9d1f2415cf672ddbe864c2d4092287b33a21/go-fuzz-build/main.go#L570) from the program under test and uses these in the [mutator](https://github.com/dvyukov/go-fuzz/blob/90825f39c90b713570ea0cc748b0987937ae6288/go-fuzz/mutator.go#L346) to construct new inputs.

In [literals/sha256.go](./literals/sha256.go) an input will cause a crash if its sha256 sum matches that
of the string ```really to long to be guessed```. Go-fuzz finds this immediately. Try it!

```bash
$ cd literals
$ go-fuzz-build
$ go-fuzz -func FuzzLiteral -bin literals-fuzz.zip
go-fuzz -bin literals-fuzz.zip -func FuzzLiteral
2021/01/29 20:09:46 workers: 8, corpus: 4 (2s ago), crashers: 1, restarts: 1/0, execs: 0 (0/sec), cover: 0, uptime: 3s
^C2021/01/29 20:09:49 shutting down...
```

A fuzzer using randomly generated strings would take a very long time to find a suitable input while go-fuzz will
find the string immediately as it has collected all the string (and int) literals.

## Sonar    


## Versify


## Mutations


## Fuzz Return Values