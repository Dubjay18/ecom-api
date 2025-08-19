# 🚀 Deployment Guide

This guide covers various deployment strategies for the E-commerce API, from local development to production environments.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Configuration](#environment-configuration)
- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Railway Deployment](#railway-deployment)
- [Heroku Deployment](#heroku-deployment)
- [AWS Deployment](#aws-deployment)
- [Production Considerations](#production-considerations)
- [Monitoring and Logging](#monitoring-and-logging)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **Go**: 1.23.3 or higher
- **PostgreSQL**: 12.0 or higher
- **Docker**: 20.10 or higher (for containerized deployment)
- **Make**: For build automation

### External Services

- **Cloudinary Account**: For image storage (optional)
- **Stripe Account**: For payment processing (optional)

## Environment Configuration

### Required Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
SERVER_HOST=0.0.0.0
PORT=8080
SERVER_MODE=release
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=5s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=ecom_api_prod
DB_SSL_MODE=require
DB_MAX_IDLE_CONNS=25
DB_MAX_OPEN_CONNS=100
DB_MAX_LIFETIME=1h

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-min-32-chars
JWT_ACCESS_TOKEN_EXPIRY=24h
JWT_REFRESH_TOKEN_EXPIRY=168h

# External Services
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_KEY=your_api_key
CLOUDINARY_SECRET=your_api_secret
STRIPE_KEY=your_stripe_key
```

### Security Considerations

- **JWT_SECRET_KEY**: Use a strong, random key (minimum 32 characters)
- **Database Password**: Use a complex password
- **SSL Mode**: Set to `require` or `verify-full` for production
- **Server Host**: Use `0.0.0.0` for containerized deployments

## Local Development

### Quick Setup

1. **Clone and setup:**
   ```bash
   git clone https://github.com/Dubjay18/ecom-api.git
   cd ecom-api
   go mod tidy
   ```

2. **Database setup:**
   ```bash
   # Start PostgreSQL
   brew services start postgresql  # macOS
   sudo service postgresql start  # Linux
   
   # Create database
   createdb ecom_api_dev
   ```

3. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your local configuration
   ```

4. **Run migrations:**
   ```bash
   make migrate-up
   ```

5. **Start development server:**
   ```bash
   make dev  # With live reload
   # or
   make run  # Standard run
   ```

### Development with Docker

```bash
# Start database
make docker-up

# Run application
make run
```

## Docker Deployment

### Dockerfile

Create a `Dockerfile`:

```dockerfile
# Build stage
FROM golang:1.23.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN make build

# Production stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/bin/api .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./api"]
```

### Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=ecom_api
      - DB_SSL_MODE=disable
      - JWT_SECRET_KEY=your-super-secret-jwt-key-here
      - SERVER_MODE=release
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=ecom_api
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  postgres_data:
```

### Build and Deploy

```bash
# Build and start services
docker-compose up --build -d

# Run migrations
docker-compose exec api ./api migrate-up

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

## Railway Deployment

Railway provides automatic deployments from Git repositories.

### Setup Steps

1. **Connect Repository:**
   - Visit [Railway](https://railway.app)
   - Connect your GitHub repository
   - Select the `ecom-api` repository

2. **Configure Environment Variables:**
   ```env
   PORT=8080
   SERVER_MODE=release
   DB_HOST=${{Postgres.PGHOST}}
   DB_PORT=${{Postgres.PGPORT}}
   DB_USER=${{Postgres.PGUSER}}
   DB_PASSWORD=${{Postgres.PGPASSWORD}}
   DB_NAME=${{Postgres.PGDATABASE}}
   DB_SSL_MODE=require
   JWT_SECRET_KEY=your-production-jwt-secret
   CLOUDINARY_CLOUD_NAME=your_cloud_name
   CLOUDINARY_KEY=your_api_key
   CLOUDINARY_SECRET=your_api_secret
   ```

3. **Add PostgreSQL Service:**
   - Click "New" → "Database" → "PostgreSQL"
   - Railway will automatically configure connection variables

4. **Deploy:**
   - Push to main branch triggers automatic deployment
   - Monitor deployment in Railway dashboard

### Railway Configuration

Create a `railway.toml`:

```toml
[build]
builder = "NIXPACKS"

[deploy]
healthcheckPath = "/health"
healthcheckTimeout = 300
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10
```

## Heroku Deployment

### Heroku Setup

1. **Install Heroku CLI:**
   ```bash
   # macOS
   brew tap heroku/brew && brew install heroku
   
   # Linux
   curl https://cli-assets.heroku.com/install.sh | sh
   ```

2. **Login and Create App:**
   ```bash
   heroku login
   heroku create your-ecom-api
   ```

3. **Add PostgreSQL:**
   ```bash
   heroku addons:create heroku-postgresql:mini
   ```

4. **Configure Environment Variables:**
   ```bash
   heroku config:set SERVER_MODE=release
   heroku config:set JWT_SECRET_KEY=your-production-jwt-secret
   heroku config:set CLOUDINARY_CLOUD_NAME=your_cloud_name
   heroku config:set CLOUDINARY_KEY=your_api_key
   heroku config:set CLOUDINARY_SECRET=your_api_secret
   ```

5. **Deploy:**
   ```bash
   git push heroku main
   ```

### Procfile

Create a `Procfile`:

```
web: ./bin/api
release: ./bin/api migrate-up
```

## AWS Deployment

### Using AWS EC2

1. **Launch EC2 Instance:**
   - Choose Ubuntu 22.04 LTS
   - Configure security groups (ports 22, 80, 443, 8080)
   - Create or select key pair

2. **Connect and Setup:**
   ```bash
   ssh -i your-key.pem ubuntu@your-ec2-ip
   
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install Go
   wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   
   # Install PostgreSQL
   sudo apt install postgresql postgresql-contrib -y
   
   # Install Nginx
   sudo apt install nginx -y
   ```

3. **Deploy Application:**
   ```bash
   # Clone repository
   git clone https://github.com/Dubjay18/ecom-api.git
   cd ecom-api
   
   # Build application
   go mod tidy
   make build
   
   # Setup systemd service
   sudo cp ecom-api.service /etc/systemd/system/
   sudo systemctl enable ecom-api
   sudo systemctl start ecom-api
   ```

### Systemd Service

Create `/etc/systemd/system/ecom-api.service`:

```ini
[Unit]
Description=E-commerce API
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/ecom-api
ExecStart=/home/ubuntu/ecom-api/bin/api
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=SERVER_MODE=release
EnvironmentFile=/home/ubuntu/ecom-api/.env

[Install]
WantedBy=multi-user.target
```

### Nginx Configuration

Create `/etc/nginx/sites-available/ecom-api`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/ecom-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## Production Considerations

### Security

1. **HTTPS/TLS:**
   ```bash
   # Install Certbot
   sudo apt install certbot python3-certbot-nginx
   
   # Get SSL certificate
   sudo certbot --nginx -d your-domain.com
   ```

2. **Firewall:**
   ```bash
   sudo ufw allow 22/tcp
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```

3. **Database Security:**
   - Use strong passwords
   - Enable SSL/TLS
   - Restrict network access
   - Regular backups

### Performance

1. **Database Optimization:**
   ```sql
   -- Add indexes for frequently queried fields
   CREATE INDEX idx_users_email ON users(email);
   CREATE INDEX idx_products_sku ON products(sku);
   CREATE INDEX idx_orders_user_id ON orders(user_id);
   CREATE INDEX idx_orders_status ON orders(status);
   ```

2. **Connection Pooling:**
   ```env
   DB_MAX_IDLE_CONNS=25
   DB_MAX_OPEN_CONNS=100
   DB_MAX_LIFETIME=1h
   ```

3. **Caching (Optional):**
   - Implement Redis for session storage
   - Add caching middleware for product listings

### Monitoring

1. **Health Checks:**
   ```bash
   # Setup monitoring script
   #!/bin/bash
   response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
   if [ $response != "200" ]; then
       echo "API is down! Response code: $response"
       # Send alert or restart service
   fi
   ```

2. **Log Management:**
   ```bash
   # Setup log rotation
   sudo cat > /etc/logrotate.d/ecom-api << EOF
   /var/log/ecom-api/*.log {
       daily
       missingok
       rotate 52
       compress
       delaycompress
       notifempty
       sharedscripts
       postrotate
           systemctl reload ecom-api
       endscript
   }
   EOF
   ```

### Backup Strategy

1. **Database Backups:**
   ```bash
   # Daily backup script
   #!/bin/bash
   BACKUP_DIR="/home/ubuntu/backups"
   DATE=$(date +%Y%m%d_%H%M%S)
   
   pg_dump -h localhost -U postgres ecom_api > $BACKUP_DIR/ecom_api_$DATE.sql
   
   # Keep only last 7 days
   find $BACKUP_DIR -name "ecom_api_*.sql" -mtime +7 -delete
   ```

2. **Application Backups:**
   ```bash
   # Backup application and config
   tar -czf /home/ubuntu/backups/app_backup_$(date +%Y%m%d).tar.gz \
       /home/ubuntu/ecom-api \
       /etc/systemd/system/ecom-api.service \
       /etc/nginx/sites-available/ecom-api
   ```

## Monitoring and Logging

### Application Logs

The application uses structured logging with Logrus:

```bash
# View logs
sudo journalctl -u ecom-api -f

# Filter by level
sudo journalctl -u ecom-api -p err

# View logs from specific time
sudo journalctl -u ecom-api --since "2024-01-15 10:00:00"
```

### Metrics Collection

Optional Prometheus integration:

```go
// Add to main.go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Add metrics endpoint
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### Alerting

Setup basic alerting with email notifications:

```bash
# Install mailutils
sudo apt install mailutils -y

# Create alert script
#!/bin/bash
if ! curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "E-commerce API is down!" | mail -s "API Alert" admin@example.com
fi
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors:**
   ```bash
   # Check PostgreSQL status
   sudo systemctl status postgresql
   
   # Check database logs
   sudo tail -f /var/log/postgresql/postgresql-*.log
   
   # Test connection
   psql -h localhost -U postgres -d ecom_api
   ```

2. **Port Already in Use:**
   ```bash
   # Find process using port 8080
   sudo lsof -i :8080
   
   # Kill process if needed
   sudo kill -9 <PID>
   ```

3. **Migration Failures:**
   ```bash
   # Check migration status
   migrate -path migrations -database "postgres://user:pass@host:port/db?sslmode=disable" version
   
   # Force migration version
   migrate -path migrations -database "postgres://user:pass@host:port/db?sslmode=disable" force <version>
   ```

4. **Memory Issues:**
   ```bash
   # Check memory usage
   free -h
   
   # Check process memory
   ps aux | grep api
   
   # Add swap if needed
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```

### Debug Mode

For troubleshooting, enable debug mode:

```env
SERVER_MODE=debug
LOG_LEVEL=debug
```

### Health Check Endpoints

The API provides health check endpoints:

```bash
# Basic health check
curl http://localhost:8080/health

# Detailed system info (debug mode only)
curl http://localhost:8080/health/detailed
```

---

This deployment guide covers the most common deployment scenarios. Choose the option that best fits your infrastructure requirements and technical expertise.