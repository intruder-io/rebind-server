This is a web server that will shut off in response to requests to `/block`. This is useful for exploiting DNS rebinding against Chrome and Edge, as described in the accompanying [blog post](https://intruder.io/research/tricks-for-split-second-dns-rebinding).

# Installation
Ensure that you have golang properly installed, and then install this tool with:
```
go get github.com/intruder-io/rebind-server@latest
```

# Usage
The port to listen on can be specified with `-p`, and the directory to serve files from can be specified with `-a` (default `./assets`). So, to listen on port 9000, serving files from `./my-exploit`, you can run:
```
rebind-server -p 9000 -a ./my-exploit
```

The server will shut off after a request is made to `/block`. While testing, it can often be helpful to run the server in a loop:
```
while :; do rebind-server -p 8080; sleep 2; done
```

You will likely have to run this server directly on the host - TCP forwarding won't work.
