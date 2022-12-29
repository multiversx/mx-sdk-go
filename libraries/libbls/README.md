# libbls

BLS signature utilities, compiled as a library (shared object).

## Build

On Linux:

```
go build -buildmode=c-shared -o libbls.so .
```

On MacOS:

```
go build -buildmode=c-shared -o libbls.dylib .
```
