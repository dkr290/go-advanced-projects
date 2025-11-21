# Deployment Guide

Production deployment guide for Wan2.1 Video Server.

## Table of Contents

- [Deployment Options](#deployment-options)
- [Docker Deployment](#docker-deployment)
- [Bare Metal Deployment](#bare-metal-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Monitoring](#monitoring)
- [Security](#security)
- [Scaling](#scaling)

---

## Deployment Options

1. **Docker** - Recommended for isolated, reproducible deployments
2. **Bare Metal** - Direct installation on server
3. **Cloud** - AWS, GCP, Azure with GPU instances
4. **Kubernetes** - For large-scale deployments

---

## Docker Deployment

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- NVIDIA Docker runtime (for GPU support)

### Install NVIDIA Docker

```bash
# Ubuntu/Debian
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list

sudo apt-get update
sudo apt-get install -y nvidia-docker2
sudo systemctl restart docker
```

### Build and Run

```bash
# Build image
docker build -t wan2-video-server:latest .

# Run with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop
docker-compose down
```

### Docker Compose Configuration

Edit `docker-compose.yml` for production:

```yaml
version: '3.8'

services:
  wan2-video-server:
    image: wan2-video-server:latest
    restart: always
    ports:
      - "8080:8080"
    environment:
      - ENABLE_GPU=true
      - LOG_LEVEL=info
    volumes:
      - ./models:/app/models
      - ./outputs:/app/outputs
      - ./logs:/app/logs
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
```

---

## Bare Metal Deployment

### System Requirements

- Ubuntu 20.04/22.04 or similar Linux distribution
- NVIDIA GPU with CUDA 11.8+
- 16GB+ RAM (32GB recommended)
- 50GB+ SSD storage

### Installation Steps

1. **Install System Dependencies**

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install build essentials
sudo apt install -y build-essential git curl wget

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Python
sudo apt install -y python3.10 python3.10-venv python3-pip

# Install CUDA (if not already installed)
wget https://developer.download.nvidia.com/compute/cuda/11.8.0/local_installers/cuda_11.8.0_520.61.05_linux.run
sudo sh cuda_11.8.0_520.61.05_linux.run
```

2. **Setup Application**

```bash
# Clone repository
git clone <repo-url>
cd wan2-video-server

# Run setup
chmod +x setup.sh
./setup.sh

# Download model
./wan2-video-server download
```

3. **Create Systemd Services**

**Python Backend Service:** `/etc/systemd/system/wan2-python-backend.service`

```ini
[Unit]
Description=Wan2.1 Python Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/wan2-video-server/python_backend
Environment="PATH=/opt/wan2-video-server/python_backend/venv/bin"
ExecStart=/opt/wan2-video-server/python_backend/venv/bin/python server.py
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Go Server Service:** `/etc/systemd/system/wan2-video-server.service`

```ini
[Unit]
Description=Wan2.1 Video Server
After=network.target wan2-python-backend.service
Requires=wan2-python-backend.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/wan2-video-server
ExecStart=/opt/wan2-video-server/wan2-video-server
Restart=always
RestartSec=10
Environment="PYTHON_BACKEND_URL=http://localhost:5000"

[Install]
WantedBy=multi-user.target
```

4. **Enable and Start Services**

```bash
sudo systemctl daemon-reload
sudo systemctl enable wan2-python-backend
sudo systemctl enable wan2-video-server
sudo systemctl start wan2-python-backend
sudo systemctl start wan2-video-server

# Check status
sudo systemctl status wan2-python-backend
sudo systemctl status wan2-video-server
```

---

## Cloud Deployment

### AWS EC2

**Recommended Instance:** `g4dn.xlarge` or `g4dn.2xlarge`

```bash
# Launch instance
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type g4dn.xlarge \
  --key-name my-key \
  --security-group-ids sg-xxxxxx \
  --subnet-id subnet-xxxxxx

# SSH to instance
ssh -i my-key.pem ubuntu@<instance-ip>

# Follow bare metal installation steps
```

**Security Group Rules:**
- Port 22 (SSH) - Your IP only
- Port 8080 (API) - Your application network
- Port 443 (HTTPS) - If using SSL

### GCP Compute Engine

**Recommended Instance:** `n1-standard-4` with NVIDIA T4

```bash
# Create instance
gcloud compute instances create wan2-video-server \
  --zone=us-central1-a \
  --machine-type=n1-standard-4 \
  --accelerator=type=nvidia-tesla-t4,count=1 \
  --image-family=ubuntu-2004-lts \
  --image-project=ubuntu-os-cloud \
  --maintenance-policy=TERMINATE \
  --boot-disk-size=100GB

# SSH and install
gcloud compute ssh wan2-video-server
```

### Azure VM

**Recommended Instance:** `Standard_NC6s_v3`

```bash
# Create VM
az vm create \
  --resource-group myResourceGroup \
  --name wan2-video-server \
  --image UbuntuLTS \
  --size Standard_NC6s_v3 \
  --admin-username azureuser \
  --generate-ssh-keys
```

---

## Reverse Proxy Setup

### Nginx Configuration

```nginx
# /etc/nginx/sites-available/wan2-video-server

upstream wan2_backend {
    server localhost:8080;
}

server {
    listen 80;
    server_name your-domain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # File upload size
    client_max_body_size 100M;

    location / {
        proxy_pass http://wan2_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts for long-running requests
        proxy_connect_timeout 600s;
        proxy_send_timeout 600s;
        proxy_read_timeout 600s;
    }

    location /outputs/ {
        alias /opt/wan2-video-server/outputs/;
        expires 1h;
        add_header Cache-Control "public, immutable";
    }
}
```

Enable site:

```bash
sudo ln -s /etc/nginx/sites-available/wan2-video-server /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## Monitoring

### Prometheus Metrics

Add to your Go server (future enhancement):

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

http.Handle("/metrics", promhttp.Handler())
```

### Logging

Configure centralized logging:

```bash
# Install Loki for log aggregation
docker run -d --name=loki -p 3100:3100 grafana/loki:latest

# Install Promtail for log shipping
docker run -d --name=promtail \
  -v /var/log:/var/log \
  -v /opt/wan2-video-server/logs:/app/logs \
  grafana/promtail:latest
```

### Grafana Dashboard

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

---

## Security

### Best Practices

1. **Use HTTPS only**
   - Get SSL certificate with Let's Encrypt
   - Force HTTPS redirects

2. **Authentication**
   - Add API key authentication
   - Use JWT tokens

3. **Rate Limiting**
   - Already implemented in application
   - Add nginx rate limiting as backup

4. **Firewall**
   ```bash
   sudo ufw allow 22/tcp
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```

5. **Regular Updates**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

---

## Scaling

### Horizontal Scaling

Use load balancer with multiple instances:

```nginx
upstream wan2_cluster {
    least_conn;
    server 10.0.1.10:8080;
    server 10.0.1.11:8080;
    server 10.0.1.12:8080;
}
```

### Vertical Scaling

- Increase `MAX_CONCURRENT_REQUESTS`
- Add more GPU devices
- Increase memory allocation

### Queue System

For production, consider adding Redis queue:

```bash
# Install Redis
docker run -d -p 6379:6379 redis:alpine

# Modify application to use Redis for job queue
```

---

## Backup and Recovery

### Backup Strategy

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR=/backups/wan2-video-server
DATE=$(date +%Y%m%d)

# Backup models
tar -czf $BACKUP_DIR/models-$DATE.tar.gz /opt/wan2-video-server/models/

# Backup configuration
cp /opt/wan2-video-server/.env $BACKUP_DIR/env-$DATE

# Backup outputs (optional)
# tar -czf $BACKUP_DIR/outputs-$DATE.tar.gz /opt/wan2-video-server/outputs/

# Keep last 7 days
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

Add to crontab:
```bash
0 2 * * * /opt/wan2-video-server/backup.sh
```

---

## Health Checks

### Uptime Monitoring

Use services like:
- UptimeRobot
- Pingdom
- StatusCake

Monitor endpoint: `https://your-domain.com/health`

### Automated Restart

```bash
#!/bin/bash
# healthcheck.sh

if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "Service unhealthy, restarting..."
    sudo systemctl restart wan2-video-server
fi
```

Crontab:
```bash
*/5 * * * * /opt/wan2-video-server/healthcheck.sh
```

---

## Troubleshooting

See logs:
```bash
# Systemd logs
sudo journalctl -u wan2-video-server -f
sudo journalctl -u wan2-python-backend -f

# Application logs
tail -f /opt/wan2-video-server/logs/app.log

# Nginx logs
tail -f /var/log/nginx/error.log
```

Common issues documented in `README.md` troubleshooting section.
