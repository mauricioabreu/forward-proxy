# Forward Proxy

> A simple forward proxy to test my understanding of servers and the HTTP protocol.

## What is it?

[Challenge reference](https://codingchallenges.fyi/challenges/challenge-forward-proxy)

> This challenge is to build your own HTTP Proxy Server. A proxy is a server that sits between a client that wants to get a resource and a server that provides the resource.

## Why go?

* Easy to deal with HTTP request and response cycle;
* Concurrent server out of the box;
* Testing HTTP servers is a breeze.

## How to run?

If you have [just](https://github.com/casey/just) installed, here is the list of commands you can run:

```
Available recipes:
    lint # Check code quality
    run  # Run proxy server
    test # Execute test suite
```

If you don't have it, you can run the commands described inside `justfile` manually.

## How to use

`just run` starts the server. Then you can use `curl` to make a request and check the response and the headers.

```
curl -v --proxy "http://localhost:8989" "http://httpbin.org/ip"

*   Trying [::1]:8989...
* Connected to localhost (::1) port 8989
> GET http://httpbin.org/ip HTTP/1.1
> Host: httpbin.org
> User-Agent: curl/8.4.0
> Accept: */*
> Proxy-Connection: Keep-Alive
>
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 45
< Content-Type: application/json
< Date: Tue, 05 Mar 2024 23:30:45 GMT
< Server: gunicorn/19.9.0
<
{
  "origin": "[::1]:50021, 127.0.0.1"
}
```
