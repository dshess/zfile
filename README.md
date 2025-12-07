# zfile

[![Go Reference](https://pkg.go.dev/badge/github.com/dshess/zfile.svg)](https://pkg.go.dev/github.com/dshess/zfile)
[![Go Report Card](https://goreportcard.com/badge/github.com/dshess/zfile)](https://goreportcard.com/report/github.com/dshess/zfile)

Package zfile implements readers and writers which inject de/compression
based on file extension.

## Install

`go get github.com/dshess/zfile`

Or just import it and let go guide you.

## Requirements

This was developed under go `1.25.5`.  TBD: Tests work under 1.20.14.

## Documentation

https://pkg.go.dev/github.com/dshess/zfile

## Usage

This code writes a plain-text file:

	path := "file"
	contents := "This is a test"
	err := os.WriteFile(path, contents, 0644)
	if err != nil {
	    return err
	}

This writes a gzip-compressed file:

	path := "file.gz"
	contents := "This is a test"
	err := zfile.WriteFile(path, contents, 0644)
	if err != nil {
	    return err
	}

Similar wrappers for os.ReadFile(), os.Open(), and os.Create().  This
module handles gzip, zstd and xz.

## License

Licensed under the MIT License, see `LICENSE` file.
