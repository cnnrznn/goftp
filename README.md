# goftp

GoFTP is a simple file transfer library and cli tool.
Simply specify a target machine and a filename, and watch it go.

Keep this library simple so it can be a flexible tool in other projects.

## Stage 1 goals

[ ] Server - listen on a port and perform receive protocol
[ ] Client - Send file to destination via protocol
[ ] No retries, fail hard