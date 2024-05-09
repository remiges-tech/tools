This program parses an sqlc (https://sqlc.dev) file to extract query names and then generates Go code with those names as constants.

```
go run main.go help
```

OR

```
go install github.com/remiges-tech/tools/sqlc-query-names-const@latest
sqlc-query-names-const help
```

It is usefule for cases like instrumenting Go code to collect metrics on query execution.
