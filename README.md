# PoC WebSocket Proxy

## Suggestions

- websockify (python solution)
- ngingx

### Using `websockify`

`websockify` is a tool that can proxy WebSocket connections. It can be used to multiplex multiple WebSocket connections through a single connection to the target server.

#### Installation

You can install `websockify` using `pip`:

```bash
pip install websockify
```

#### Usage

To start a WebSocket proxy with `websockify`, you can use the following command:

```bash
websockify --web /path/to/web 8080 target_server:target_port
```

- `8080` is the port on which `websockify` will listen for incoming WebSocket connections.
- `target_server:target_port` is the address of the target WebSocket server.

### Using Nginx

Nginx can also be configured to proxy WebSocket connections. This approach allows you to use Nginx's powerful load balancing and connection management features.

#### Installation

If you don't have Nginx installed, you can install it using your package manager. For example, on Debian-based systems:

```bash
sudo apt-get update
sudo apt-get install nginx
```

#### Configuration

Create or modify your Nginx configuration file to include a WebSocket proxy configuration. Here is an example configuration:

```nginx
http {
    upstream websocket_backend {
        server target_server:target_port;
    }

    server {
        listen 8080;

        location / {
            proxy_pass http://websocket_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "Upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

- `target_server:target_port` is the address of the target WebSocket server.
- `8080` is the port on which Nginx will listen for incoming WebSocket connections.

#### Start Nginx

After configuring Nginx, start or reload the Nginx service:

```bash
sudo systemctl restart nginx
```

### Summary

Both `websockify` and Nginx can be used to proxy WebSocket connections. `websockify` is a lightweight tool specifically designed for WebSocket proxying, while Nginx offers more advanced features like load balancing and connection management. Choose the tool that best fits your needs and infrastructure.