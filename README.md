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
