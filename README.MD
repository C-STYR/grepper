## Multi-line Grepping Tool written in Golang

&nbsp;
#### ---this tool is heavily based on one written by [Jayson Lennon](https://github.com/jayson-lennon) and is mainly for me to play with---
&nbsp;

### To use this tool, once available locally, cd into ./grepper and on the command line run:

```bash
go run ./grepper [search term] [filepath]
```

Using this tool will return results in the form of "filename[line number]: [line of code/text]", 
for example: 

```go
worklist/worklist.go[8] jobs chan Entry
```