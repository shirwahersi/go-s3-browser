# Go S3 Browser

A lightweight, standalone Go web application for browsing S3-compatible storage buckets. Provides a clean web UI with folder navigation, file listing, and secure downloads via pre-signed URLs.

## Features

- Browse S3 bucket contents with folder navigation
- File type icons for visual recognition
- Secure file downloads via pre-signed URLs (1-hour expiry)
- Responsive design (mobile-friendly)
- Single binary deployment (static files embedded)
- Customizable page title and description
- YAML configuration with environment variable overrides
- Works with any S3-compatible storage (AWS S3, Contabo, MinIO, etc.)

## Quick Start

### Prerequisites

- Go 1.21 or later

### Build

```bash
git clone https://github.com/shirwahersi/go-s3-browser.git
cd go-s3-browser
go build -o go-s3-browser
```

### Configure

Copy the example config and edit with your S3 credentials:

```bash
cp config.example.yaml config.yaml
```

Edit `config.yaml`:

```yaml
port: 3000

site:
  title: "My File Storage"
  description: "Browse and download files"

s3:
  endpoint: https://eu2.contabostorage.com
  region: eu2
  bucket: your-bucket-name
  access_key_id: your-access-key-id
  secret_access_key: your-secret-access-key
```

### Run

```bash
./go-s3-browser
```

Open http://localhost:3000 in your browser.

## Configuration

The application supports two configuration methods. Environment variables take precedence over the config file.

### Config File (`config.yaml`)

```yaml
port: 3000

site:
  title: "My File Storage"
  description: "Browse and download files from our storage"

s3:
  endpoint: https://eu2.contabostorage.com
  region: eu2
  bucket: my-bucket
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `3000` |
| `SITE_TITLE` | Page title/header | `S3 Bucket Browser` |
| `SITE_DESCRIPTION` | Description below title | (empty) |
| `S3_ENDPOINT` | S3-compatible endpoint URL | - |
| `S3_REGION` | Storage region | - |
| `S3_BUCKET` | Bucket name | - |
| `S3_ACCESS_KEY_ID` | Access key credentials | - |
| `S3_SECRET_ACCESS_KEY` | Secret key credentials | - |

Example using environment variables only:

```bash
export S3_ENDPOINT=https://eu2.contabostorage.com
export S3_REGION=eu2
export S3_BUCKET=my-bucket
export S3_ACCESS_KEY_ID=your-access-key
export S3_SECRET_ACCESS_KEY=your-secret-key
./go-s3-browser
```

## Deployment with Systemd

### Installation

1. Build and copy the binary:

```bash
go build -o go-s3-browser
sudo mkdir -p /opt/go-s3-browser
sudo cp go-s3-browser /opt/go-s3-browser/
```

2. Create the systemd service file:

```bash
sudo nano /etc/systemd/system/go-s3-browser.service
```

### Option A: Using Config File

```ini
[Unit]
Description=Go S3 Browser
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/go-s3-browser
ExecStart=/opt/go-s3-browser/go-s3-browser
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Copy your config file:

```bash
sudo cp config.yaml /opt/go-s3-browser/
sudo chown www-data:www-data /opt/go-s3-browser/config.yaml
sudo chmod 600 /opt/go-s3-browser/config.yaml
```

### Option B: Using Environment Variables

```ini
[Unit]
Description=Go S3 Browser
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/go-s3-browser
Environment="PORT=3000"
Environment="S3_ENDPOINT=https://eu2.contabostorage.com"
Environment="S3_REGION=eu2"
Environment="S3_BUCKET=your-bucket-name"
Environment="S3_ACCESS_KEY_ID=your-access-key-id"
Environment="S3_SECRET_ACCESS_KEY=your-secret-access-key"
ExecStart=/opt/go-s3-browser/go-s3-browser
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Alternatively, use an environment file for better security:

```bash
sudo nano /opt/go-s3-browser/.env
```

```
PORT=3000
S3_ENDPOINT=https://eu2.contabostorage.com
S3_REGION=eu2
S3_BUCKET=your-bucket-name
S3_ACCESS_KEY_ID=your-access-key-id
S3_SECRET_ACCESS_KEY=your-secret-access-key
```

```bash
sudo chown www-data:www-data /opt/go-s3-browser/.env
sudo chmod 600 /opt/go-s3-browser/.env
```

Then reference it in the service file:

```ini
[Service]
EnvironmentFile=/opt/go-s3-browser/.env
```

### Enable and Start the Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable go-s3-browser
sudo systemctl start go-s3-browser
```

### Check Status

```bash
sudo systemctl status go-s3-browser
sudo journalctl -u go-s3-browser -f  # View logs
```

## Nginx Reverse Proxy

### Basic HTTP Configuration

Create a new site configuration:

```bash
sudo nano /etc/nginx/sites-available/s3-browser
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:

```bash
sudo ln -s /etc/nginx/sites-available/s3-browser /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### HTTPS with Let's Encrypt

Install certbot and obtain a certificate:

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

Certbot will automatically update your Nginx configuration to use HTTPS.

### HTTPS Configuration (Manual)

If you have your own SSL certificates:

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/fullchain.pem;
    ssl_certificate_key /path/to/privkey.pem;

    # SSL settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    location / {
        proxy_pass http://127.0.0.1:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## License

MIT License
