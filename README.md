# Tough Fuzzer

My experience of introducing
go-fuzz to developers has taught me that the project README and preconceptions
about *fuzzing* can lead to limiting mental models of go-fuzz. 
In particular, the idea that *randomness* is the most important property
of the fuzzer is common. It's hard to imagine a simple harness
for injecting random inputs into programs being so [unreasonably effective](https://github.com/dvyukov/go-fuzz#trophies). Go-fuzz has a lot of clever
strategies for covering as much code as possible to root out bugs.

I think it's important for developers using go-fuzz to understand these
clever strategies so they can maximize the effectiveness of their fuzzing
resources. Fuzzing can be resource intensive and often runs asynchronously to
the normal development process. The quicker a fuzzer finds bugs, the more
chance that fixes land before code is shipped. Without knowing how the
fuzzer works it's easy (and, in my experience, common) to reduce the
effectiveness of fuzzing substantially.

Tough Fuzzer is an obstacle course for [go-fuzz](https://github.com/dvyukov/go-fuzz) composed of a series of small code samples
which encapsulate the most common *obstacles* to code-coverage the fuzzer
will encounter. In each case, the obstacle is insurmountable in a reasonable
period of time using random inputs or even coverage-guided mutation.

In each case, running the fuzzer should result in near instantaneous discovery
of a crasher. You'll have to use your imagination and replace the ```panic``` calls
with plausible code paths. The point of the examples is showing how certain
paths through the code can be reached.


## What does it mean to overcome "obstacles" in a program?

Coverage-guided fuzzers, like go-fuzz, generate input using a kind of genetic algorithm.
Given a function to test and an initial *corpus* of test cases, the algorithm runs the function for each test case
and measures the resulting [*code-coverage*](https://en.wikipedia.org/wiki/Code_coverage). Each element of the corpus
is *mutated* (bit-flipping etc) to expand the corpus. That's the random
part of the fuzzer. These new test cases are tried with the function under test
and the code-coverage is measured. The clever part is that inputs which result in
greater code-coverage are deemed to be *fit* and prioritized for further
mutation. By this evolutionary process, new paths through the 
code are discovered which hopefully trigger hidden bugs.

![Genetics](https://media.giphy.com/media/G8k4UcUNIhFSM/giphy.gif)

For example, it's plausible for a coverage-guided fuzzer to navigate the following and cause a panic:

```golang
func buggy(data []byte) {
    if len(data) < 4 {
        return
    }
    if data[0] == 1 {
        if data[1] == 2 {
            if data[2] == 3 {
                if data[3] == 4 {
                    panic("bug found") // try to reach here
                }
            }
        }
    }
}
```
The fuzzer can do this without having to guess all four bytes of ```data``` (taking on average 2^31 tries with uniformly distributed guesses). It can,
once on the right track, guess each byte one after the other with on the order of 4*2^7 mutations. 

An obstacle would be some code which causes an evolutionary dead-end. For example, consider the following function:

```golang
func buggy(i int) {
    if i == 1337 {
        panic("bug found")
    }
}
```

In this case, mutating an initial corpus will only result in new code-coverage once the magic ```1337``` value is discovered.
This is equivalent to guessing an int (64 bits on my machine) at random. That's not tractable. Happily, go-fuzz has some
strategies for coping with this which will be discussed later.

Another example arises when CRC checksums are used:

```golang
func buggyCRC(data []byte) {
	if len(data) <= 4 {
		return
	}
	sum := crc32.ChecksumIEEE(data[4:])
	if binary.BigEndian.Uint32(data[:4]) == sum {
        panic("bug found")
    }
}
```

How can mutation of inputs hope to find inputs which have a valid CRC checksum as the first four bytes? Mutation will be no better than
random guessing. Luckily, go-fuzz has tricks to get around these obstacles and that's what reading the code samples (and running them!) 
will hopefully showcase.

In each sub-directory, there are example fuzzers which demonstrate how go-fuzz can leap over some common fuzzing obstacles. Reaching
a ```panic``` is considered to be "finding a bug". This is obviously quite artificial but it shouldn't take too much imagination to
see how each applies to real-world code.

> Try ```./run.sh``` to run each of the fuzzers


## Int and String Literals

Go-fuzz [extracts int and string literals](https://github.com/dvyukov/go-fuzz/blob/6a8e9d1f2415cf672ddbe864c2d4092287b33a21/go-fuzz-build/main.go#L570) from the program under test and uses these in the [mutator](https://github.com/dvyukov/go-fuzz/blob/90825f39c90b713570ea0cc748b0987937ae6288/go-fuzz/mutator.go#L346) to construct new inputs.

In [literals/sha256.go](./literals/sha256.go) an input will cause a crash if its sha256 sum matches that
of the string ```really too long to be guessed```. Go-fuzz finds this immediately. Try it!

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

When ```go-fuzz-build``` instruments code it adds functions which
report when certain types of comparison are made. For example the [code](https://github.com/stevenjohnstone/toughfuzzer/blob/master/sonar/integers.go#L17)
```golang
func findInt(i int) {
	if 1337 == i {
		panic("Found int")
	}
}
```
is instrumented as
```golang
func findInt(i int) {
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:11
	_go_fuzz_dep_.CoverTab[23576]++
											if func() _go_fuzz_dep_.Bool {
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:12
		__gofuzz_v2 := i
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:12
		_go_fuzz_dep_.Sonar(1337, __gofuzz_v2, 793408)
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:12
		return 1337 == __gofuzz_v2
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:12
	}() == true {
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:12
		_go_fuzz_dep_.CoverTab[46109]++
												panic("Found int")
	} else {
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:14
		_go_fuzz_dep_.CoverTab[26146]++
//line /home/stevie/go/src/github.com/stevenjohnstone/toughfuzzer/sonar/integers.go:14
	}
}
```

When go-fuzz runs, ```_go_fuzz_dep_.Sonar(1337, __gofuzz_v2, 793408)``` reports that a comparison has been made
between a variable and ```1337```. The value ```1337``` is used
to generate new inputs. Not only that, various encodings of the
integer ```1337``` are added to the corpus. In particular, as
demonstrated in [sonar/integers.go](./sonar/integers.go), ```1337```
is added to the corpus as a

* big-endian array of bytes
* little-endian array of bytes
* hex string
* decimal string

This helps go-fuzz deftly manoeuvre past comparison with magic
numbers. Helpfully, it anticipates that the function under test
may use a variety of representations for an integer. 

> To explore the go-fuzz instrumentation, use ```go-fuzz-build -work``` to preserve the work directory: the working directory path will be displayed when the build completes

Note that values reported by sonar don't need to be constants: values can
be calculated at runtime. For example, CRC checksums could be considered to
be a barrier to simple mutation fuzzing. So much so that [AFL](https://github.com/google/AFL) requires libpng
to be [patched](https://github.com/google/AFL/tree/master/experimental/libpng_no_checksum) to remove checksums to make progress. Not so with go-fuzz. See [this](./sonar/checksums.go) example for sonar ducking under a CRC check without
missing a beat.

![Jesse Magic](https://media.giphy.com/media/NmerZ36iBkmKk/giphy.gif)

> Convince yourself that sonar is doing the work here by adding ```-sonar=false``` to the invocation of ```go-fuzz```

## Big Event

In [race.go](./bigevent/race.go) I've strung together some obstacles covered above to make this a bit more challenging. Go-fuzz
should make short work of it.

On my mid-range laptop, I get the following results on running ```run.sh```

```shell
Running fuzzers in sonar
Running FuzzCheckSum
Found crasher for FuzzCheckSum after 6 seconds: "\xa3\xe7t\x01\x03\x80\xedӱ\x92Hգ\f*w"
Running FuzzIntegerBigEndian
Found crasher for FuzzIntegerBigEndian after 6 seconds: "\x00\x00\x00\x00\x00\x00\x059"
Running FuzzIntegerLittleEndian
Found crasher for FuzzIntegerLittleEndian after 6 seconds: "9\x05\x00\x00\x00\x00\x00\x00"
Running FuzzIntegerDecimalString
Found crasher for FuzzIntegerDecimalString after 6 seconds: "1337"
Running FuzzIntegerHexString
Found crasher for FuzzIntegerHexString after 48 seconds: "539"
Running FuzzString
Found crasher for FuzzString after 6 seconds: "gfedcba gnirts 7331"
Running fuzzers in literals
Running FuzzLiteral
Found crasher for FuzzLiteral after 6 seconds: "really too long to b" +
"e guessed"
Running fuzzers in bigevent
Running FuzzBigEvent
Found crasher for FuzzBigEvent after 27 seconds: "race2021 is on! F$ck" +
" COVID\x05\xee?\xa6"
```



