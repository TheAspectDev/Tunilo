# Tunio Protocol
Tunio uses a basic and well-known protocol for transfering states/data between server/client.

## Message Types

```go
type Type byte // 0-255

const (
	MsgReady    Type = 1
	MsgPing     Type = 2
	MsgPong     Type = 3
	MsgRequest  Type = 4
	MsgError    Type = 5
	MsgDataAddr Type = 6
)
```

## Structure

```js
[MESSAGE_TYPE][PAYLOAD_LENGTH][        PAYLOAD       ]
```

```js
[   1 Byte   ][    4 Bytes   ][ PAYLOAD_LENGTH Bytes ]
```

## Examples

```js
// Data
[5][12][syntax error]
// Structure
[MsgError][Payload length is 12 bytes][12 Bytes of data]
```

