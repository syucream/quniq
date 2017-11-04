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
  -i    enable case insentive comparison
  -inbuf-weight int
        number of input buffer items(used specified value * 1024 * 1024) (default 1)
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

## Benchmarks

### Environment 

* MacBook Air Early 2014
* Core i7 4650U 1.6GHz
* Mem 8GB

### Target file

```
$ cat /dev/urandom | tr -dc '0-9' | fold -w 4 | head -n 100000000 > randlog_0
$ cp randlog_0 randlog_1
...
```

### Execution time comparison

* sort | uniq

```
bash-3.2$ time cat randlog_* | LANG=C gsort | guniq > /dev/null

real    15m39.761s
user    12m59.955s
sys     2m10.857s
```

* sort | uniq with --parallel

```
bash-3.2$ time cat randlog_* | LANG=C gsort --parallel 4 | guniq > /dev/null

real    14m24.231s
user    12m32.867s
sys     2m7.031s
```

* awk

```
bash-3.2$ time cat randlog_* | awk '!_[$0]++' > /dev/null

real    11m13.350s
user    10m59.538s
sys     0m5.868s
```

* sort -u

```
bash-3.2$ time cat randlog_* | LANG=C gsort -u > /dev/null

real    6m4.100s
user    5m46.810s
sys     0m10.659s
```

* sort -u with --parallel

```
bash-3.2$ time cat randlog_* | LANG=C gsort -u --parallel 4 > /dev/null

real    5m56.870s
user    5m40.977s
sys     0m10.251s
```

* quniq

```
bash-3.2$ time cat randlog_* | ./quniq --max-workers 4 > /dev/null

real    1m45.362s
user    4m29.294s
sys     0m10.177s
```
