This is a Work in progress.

Example usage:

```
$ go run .\cmd\main.go -- .\data\example.go
.\data\example.go:9:2: function call to 'slog.Info': last key has no associated value.
.\data\example.go:15:2: function call to 'slog.Debug': argument #1 should be of type string, but is int.
```