# Simple Websocket Server provide with TCP server

## How to run
```
go run TcpServer.go 8000
```

## Client Test

run this on your javascript browser

```
var ws = new Websocket("ws://localhost:8000");
```

Voila!!
Now check on your browser network.
You have made one websocket connection

## Description

A WebSocket connection begins with a HTTP GET request from the client to the server, called the ‘handshake’. This request carries with it a Connection: upgrade and Upgrade: websocket header to tell the server that it wants to begin a WebSocket connection, and a Sec-WebSocket-Version header that indicates the kind of response it wants. In this guide we’ll only focus on version 13 of the protocol.

The request headers also include a Sec-WebSocket-Key field. With this, the server creates a Sec-WebSocket-Accept header that forms part of its response.

If there is a value to Sec-WebSocket-Key, according to the regular expression above, we can take that value and create the Sec-WebSocket-Accept header in our response. It does so by taking the value of the Sec-WebSocket-Key and concatenating it with "258EAFA5-E914-47DA-95CA-C5AB0DC85B11", a 'magic string', defined in the protocol specification. It takes this concatenation, creates a SHA1 digest of it, then encodes this digest in Base64. We can do this using the built-in sha1 and base64 libraries.