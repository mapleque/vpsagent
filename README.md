# vpsagent
An agent service for VPS with http api

You can run a set of commands on your VPS by the http request.

## Security

1. Using TLS on network transportation.
1. Making a signature for reuqest params.
1. Checking the timestamp to avoid request replay.
1. Ip white list.

## Deployment

1. Download [release version](https://github.com/mapleque/vpsagent/releases).
1. General cert and token for config.
1. Run with the command:
```
vspagent \
  -p <port> \
  -sign-token 'token_for_signature' \
  -tls-key-file '/path/to/tls.key' \
  -tls-cert-file '/path/to/tls.cert' \
  -ip-allow '127.0.0.1' \
  -ip-allow '192.168.0.1'
```

## Api

All request must be a POST request with HTTP Header:
- `Signature: <signature with token>`. See [Signature](#signature) for more information.
- `Timestamp: <unix timestamp>`. See [Timestamp](#timestamp) for more information.

The request body should be a runable script (bat for windows or shell for linux).

For an example:

```
# This is a http request, which can
# create a file on your linux vps.

# Request begin
POST / HTTP/1.1
Content-Type: text/plain
Signature: <sign>
Timestamp: <timestamp>

pwd
echo 'Hello vpsagent!' > hello-vpsagent.txt
ls -l
cat hello-vpsagent.txt
# Request end
```

### Signature

Signature is used to protect your request data from being modified.

The calculate method is: `md5(timestamp@md5(body))`.

In which,
- The timestamp is comes from your request time. See [Timestamp](#timestamp) for more information.
- The body is just the whole request body.

### Timestamp

The timestamp is an Unix timestamp, which comes from the request time.

With the timestamp, following rules will be check:
1. The request will be deny if there are more then 5 seconds between the timestamp and the time of request recieving.
1. If there are multiple reuqest with same timestamp and signature, only one will be processed.

