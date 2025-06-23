# Providers Endpoints

This repository contains the source code for a Go application that interacts with various cloud service providers to fetch their available regions. The application now includes **caching with Turso DB** and **Slack webhook notifications** for enhanced reliability and monitoring.

## Features

- **Multi-Provider Support**: Fetches regions from 13+ cloud providers including AWS, DigitalOcean, Google Cloud, Vultr, Linode, and more
- **Intelligent Caching**: Uses Turso DB to cache region data for 24 hours, reducing API calls and improving performance
- **Slack Notifications**: Sends notifications when:
  - Region fetching fails for any provider
  - Region data changes (new regions added/removed)
- **Concurrent Processing**: Fetches regions from multiple providers simultaneously for optimal performance
- **Fallback Support**: Falls back to cached data when fresh fetching fails
- **Vercel Integration**: Ready for deployment on Vercel platform

## Structure

- `main.go` - Application entry point with environment configuration
- `service/` - Individual provider implementations (AWS, DigitalOcean, Google Cloud, etc.)
- `lib/` - Core functionality including caching, notifications, and service orchestration
  - `lib.go` - Original region fetching logic
  - `turso.go` - Turso DB integration for caching
  - `slack.go` - Slack webhook notifications
  - `cached_service.go` - Cached service wrapper with notifications

## Configuration

The application uses environment variables for configuration. Copy `config.example.env` to create your own configuration:

```bash
# Turso Database Configuration (Required for caching)
TURSO_DATABASE_URL=libsql://your-database-name.turso.io
TURSO_AUTH_TOKEN=your-turso-auth-token

# Slack Webhook Configuration (Required for notifications)
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK

# Optional: Specific channels for different notification types
SLACK_ERROR_CHANNEL=#alerts
SLACK_CHANGES_CHANNEL=#infrastructure-changes
```

### Setting up Turso DB

1. Install Turso CLI: `curl -sSfL https://get.tur.so/install.sh | bash`
2. Create a database: `turso db create providers-cache`
3. Get the database URL: `turso db show --url providers-cache`
4. Create an auth token: `turso db tokens create providers-cache`

### Setting up Slack Notifications

1. Create a Slack app in your workspace
2. Add an Incoming Webhook integration
3. Copy the webhook URL to your environment configuration

#### Channel Selection Options

The application supports sending notifications to different channels:

1. **Default Channel**: Messages go to the channel configured in your webhook
2. **Error Channel**: Set `SLACK_ERROR_CHANNEL` to send error notifications to a specific channel (e.g., "#alerts")
3. **Changes Channel**: Set `SLACK_CHANGES_CHANNEL` to send region change notifications to a specific channel (e.g., "#infrastructure-changes")

**Note**: Channel override requires your webhook to have permissions to post to different channels, or you may need to use a bot token instead of webhooks for full channel control.

## Dependencies

This project uses several dependencies, including:

- `github.com/PuerkitoBio/goquery` for parsing HTML
- `github.com/tbxark/g4vercel` for Vercel integration
- `github.com/tursodatabase/libsql-client-go/libsql` for Turso DB connectivity

## Building and Running

To build and run this project, you need to have Go installed. Then, you can use the `go build` and `go run` commands.

```bash
# Install dependencies
go mod tidy

# Run with environment variables
TURSO_DATABASE_URL="your-url" TURSO_AUTH_TOKEN="your-token" go run main.go

# Or build and run
go build -o providers-endpoints
./providers-endpoints

# Test caching functionality (requires Turso DB configuration)
TURSO_DATABASE_URL="file:test.db" go run cmd/test_cache.go
```

## How It Works

1. **Cache Check**: First checks Turso DB for cached region data (valid for 24 hours)
2. **Fresh Fetch**: If cache miss or expired, fetches fresh data from providers
3. **Change Detection**: Compares new data with cached data to detect changes
4. **Notifications**: Sends Slack notifications for failures or changes
5. **Cache Update**: Updates cache with new data
6. **Fallback**: Returns cached data if fresh fetch fails

## Supported Providers

- Amazon AWS (S3 & EC2)
- Amazon Lightsail
- DigitalOcean (Spaces & Droplets)
- UpCloud
- Exoscale
- Google Cloud (Storage & Compute)
- Backblaze B2
- Linode
- Outscale
- Storj
- Vultr
- Hetzner
- Synology C2
