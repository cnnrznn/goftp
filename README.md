# goftp

GoFTP is a simple file transfer library and cli tool.
Simply specify a target machine and a filename, and watch it go.

Keep this library simple so it can be a flexible tool in other projects.

This library is push based.
The reason for this is that the destination machine does not have the metadata necessary to request the file.
It is up to the program using this library to determine when file transfer is necessary and expected, and invoke the receiving machine to listen for an incoming file when necessary.


## Timeline

### Update 2023/07/18

- Integration tests working
- Simplified API
- Removed unnecessary packages. All under `ftp` now.

### Update 2023/07/16

A couple of days ago I got this module into "working order," meaning untested
in any dependent projects but the integration tests are working.
In doing so, I realize that exposing the `Run()` method for the API is a big
anti-pattern.
The functions on client/server side are inherrently synchronous and it should
be left to the caller to run them in the background or not.
I will be refactoring this package to be simpler.

Further, the "Client" and "Server" naming convention is obtuse.
I will rename these to something clearer, or better combine them into a single
package.

### Update 2023/04/06

In working on this project, I realize I need a couple of generic libraries that
can be used for other projects outside of this one. One of which already exists as part
of another repo.

- github.com/cnnrznn/gonet:  A networking lib:
- github.com/cnnrznn/goargs: A command line parser

I will be pausing work on this tool until both of those are complete.

### Stage 1

- [x] Server - listen on a port and perform receive protocol
- [x] Client - Send file to destination via protocol
- [x] No retries, fail hard