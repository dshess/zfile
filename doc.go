/*
Package zfile implements readers and writers which inject de/compression
based on file extension.

# Usage

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

# What is the motivating use case?

I often find myself writing utility commands which run like this:

	go run ./ generate-data ... output.zst

Adding flags to determine the kind of output is error-prone, because you
can specify to, say, xz compress data into a file with .zst extension.  It
is convenient to be able to do things like up-arrow and change the
extension and run, without having to make two or three parallel
command-line changes.
*/
package zfile
