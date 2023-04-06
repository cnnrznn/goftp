# goftp

GoFTP is a simple file transfer library and cli tool.
Simply specify a target machine and a filename, and watch it go.

Keep this library simple so it can be a flexible tool in other projects.

## Timeline

### Update 2023/04/06

In working on this project, I realize I need a couple of generic libraries that
can be used for other projects outside of this one. One of which already exists as part
of another repo.

- github.com/cnnrznn/gonet:  A networking lib:
- github.com/cnnrznn/goargs: A command line parser

I will be pausing work on this tool until both of those are complete.

### Stage 1

- [ ] Server - listen on a port and perform receive protocol
- [ ] Client - Send file to destination via protocol
- [ ] No retries, fail hard