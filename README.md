## 📡 Tunilo — Simple, Extremely Fast Reverse Tunnel for Exposing Local Services

Tunio is a lightweight **Reverse Tunnel** that lets you expose a local HTTP service to the public internet through a TCP control channel.

It works similarly to ngrok, cloudflare tunnel, and localtunnel, but is intentionally minimal and easy to extend.

> The difference this project makes besides minimalism is that the public facing server is yours, basically making it a **self-hosted reverse-tunnel**.

## How it works
Device #1:
- HTTP Handlers ( Web Server )
- Tunilo/Client installed and running
- Can be under **NAT**

Device #2:
- Tunilo/Server
- - Control Server :9090
- - Public Server :4311
- Is the server that the public is going to request to

> Running
**Device** #2 starts running **Tunilo/Server**, after it is turned on, **Device #1** should turn on and connect to **Tunilo/Server** on address `device_2_ip:4311` by forming a TCP Connection.

>HTTP requests are forwarded using the **control server connection** and returned back on the same connection using tunilos framed protocol

## Tunilo Protocol
You can read more about the protocol at [Protocol](/internal/protocol/protocol.md)

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

## License
All the code is licensed under [MIT](https://choosealicense.com/licenses/mit/)