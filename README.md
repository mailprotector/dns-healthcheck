# dns-healthcheck
A light-weight utility for testing a DNS server via an HTTP request

## How it works
dns-healthcheck sets up a web server that sends a DNS query to the specified DNS server and returns a 2xx HTTP code for successful queries, and a 5xx for failed queries.
## Environment Variables

| Name     | Description                        |  Type  |       Default       | Required |
| -------- | ---------------------------------- | :----: | :-----------------: | :------: |
| RESOLVER | IP address of the resolver to test | string |      127.0.0.1      |    no    |
| LOOKUP   | DNS name to look up                | string | 2.0.0.127.my.domain |    no    |
| HTTPPORT | HTTP listening port                | string |        8080         |    no    |
| LOGLEVEL | Log verbosity                      | string |        info         |    no    |

## Build

### Go
`env GOOS=linux GOARCH=amd64 go build -o bin/dns-healthcheck_0.0.1_linux_amd64`

### Docker
```
docker build -t dns-healthcheck
```

## Usage
The most simplistic way to run the test

```
curl -s http://127.0.0.1/health

{"result": alive}
```

You can also send a custom HTTP header to lookup an address other than the defined env var
```
curl -H 'lookup:some.address.com' -s http://127.0.0.1/health

{"result": alive}
```