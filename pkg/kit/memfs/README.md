# The `memfs` package

The `memfs` package provides in-memory filesystem.

## The `File`

The `File` structure implements multiple I/O interfaces plus `fs.File` 
interface. Its purpose is to be a versatile buffer. 

```go
type file interface {
    io.Seeker
    io.Reader
    io.ReaderAt
    io.Closer
    io.ReaderFrom
    io.Writer
    io.WriterAt
    io.StringWriter
    fs.File
    
    Truncate(size int64) error
}
```
