# redirector

Redirector is a Go webserver that responds to requests with either a 404 or a 302. It is configured with a
simple text configuration.

## Usage

```
$ ./redirector -h
Usage of ./redirector:
  -addr="localhost:9310": Listen addr
  -conf="sample.conf": Config file
```

## Configuration

The configuration file is extremely simple:

```
from to
...  ...
```

`from` and `to` must both be present on each line and must be separated by any amount of whitespace.

If the configuration is

```
foo http://google.com
```

then requests to `/foo`, `/foo/`, etc will get a 302 redirect to http://google.com, while everything else will
get a 404.

The configuration is reloaded on each request in order to pick up changes (under load it should come out of
disk cache pretty much all the time).
