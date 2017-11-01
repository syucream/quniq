# quniq

Accelerate uniq with using multi goroutines and large amount of memory.

## Installation

```
$ go get github.com/syucream/quniq
```

## Usage

```
Usage of ./quniq:
  -c    print with count
  -d    output only duplicated lines
  -max-workers int
        number of max workers (default 1)
  -u    output only uniuqe lines
```

* for example:

```
$ cat file | quniq -c -max-workers 2
```

## Note

* It doesn't require input data is sorted.
* Its order of output lines is random.

