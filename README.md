# fastuniq

Accelerate uniq with multi process and large amount memory.

## Installation

```
$ go get github.com/syucream/fastuniq
```

## Usage

```
  -c    print with count
  -max-workers int
        number of max workers (default 1)
```

* for example:

```
$ cat file | fastuniq -c -workers 2
```
