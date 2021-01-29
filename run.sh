#!/bin/bash -e

command -v gopls &> /dev/null || {
    echo "gopls not found: try 'GO111MODULE=on go get golang.org/x/tools/gopls@latest' or check PATH is set"
    exit
}

command -v go-fuzz-build &> /dev/null || {
    echo "go-fuzz-build not installed: try 'go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build' or check PATH is set"
    exit
}

function fuzz() {
    echo "Running fuzzers in $1"
    pushd "$1" &> /dev/null
    go-fuzz-build
    for f in ./*.go; do
        [ -f "$f" ] && {
            gopls symbols "$f"|grep "Fuzz"|cut -f 1 -d' '|while read -r func; do
                echo "Running $func"
                start=$SECONDS
                workdir=$(mktemp -d)
                go-fuzz -bin "$1"-fuzz.zip -workdir "$workdir" -func "$func" 2>&1 | while read -r line; do
                    echo "$line" | grep -q "crashers: [1-9][0-9]*" && {
                        break
                    }
                done
                end=$SECONDS
                duration=$(( end - start ))
                echo -n "Found crasher for $func after $duration seconds: "
                find "$workdir/crashers" -name "*.quoted" -exec cat {} \; | while read -r line; do
                    echo "$line"
                done
                rm -rf "$workdir"
            done
        }
    done
    popd &> /dev/null
}

dirs=("sonar" "literals")


for dir in "${dirs[@]}"
do
fuzz "$dir"
done





