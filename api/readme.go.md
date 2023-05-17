## Go

These settings apply only when `--go` is specified on the command line.

``` yaml $(go)
go:
    license-header: MICROSOFT_MIT_NO_VERSION
    clear-output-folder: false
    export-clients: true
    output-folder: $(go-sdk-folder)
    use: "@autorest.go@4.0.0-preview.46"
    file-prefix: zz_generated_
    verbose: true
```

### Tag: preview-2023-03-01 and go

These settings apply only when `--tag=2023-03-01-preview --go` is specified on the command line.
Please also specify `--go-sdk-folder=../sdk`.

``` yaml $(tag) == '2023-03-01-preview' && $(go)
output-folder: $(go-sdk-folder)
```
