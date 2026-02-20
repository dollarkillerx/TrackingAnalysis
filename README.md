# Data Tracking System MVP

A Go-based data tracking system supporting ad tracking (click redirect / JS collection) and web analytics (pageview / event SDK). Uses JSON-RPC 2.0 over Gin, PostgreSQL via GORM, Redis for rate limiting / anti-replay / dedup, and hybrid RSA+AES encryption for secure event reporting.

## Prerequisites

- Go 1.23+
- PostgreSQL
- Redis

## Quick Start

1. **Create the database:**

```sql
CREATE DATABASE tracking;
```

2. **Configure** — edit `config.toml` with your database and Redis connection details.

3. **Run:**

```bash
go run cmd/server/main.go -config config.toml
```

RSA keys are auto-generated on first startup in the `keys/` directory.

## Configuration

See `config.toml` for all options:

| Section | Key | Description |
|---------|-----|-------------|
| `server` | `addr` | Listen address (default `:8080`) |
| `server` | `export_url` | Public-facing URL for links/SDK |
| `db` | `dsn` | PostgreSQL connection string |
| `redis` | `addr` | Redis address |
| `admin` | `username/password` | Admin credentials (hardcoded MVP) |
| `security` | `token_secret` | HMAC secret for tracking tokens |
| `security` | `ts_window_seconds` | Timestamp validity window (300s) |
| `security` | `dedup_seconds` | Click dedup window (10s) |
| `rate_limit` | `per_ip_per_minute` | Rate limit per IP (60) |
| `bot` | `block_threshold` | Bot score to block (80) |

## API Reference (JSON-RPC 2.0)

All RPC calls go to `POST /rpc` with `Content-Type: application/json`.

### Admin Methods

**Login:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.login",
    "params": {"username": "admin", "password": "changeme"},
    "id": 1
  }'
```

Response: `{"jsonrpc":"2.0","result":{"admin_token":"..."},"id":1}`

**Create Tracker (ad type):**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.tracker.create",
    "params": {
      "admin_token": "TOKEN_FROM_LOGIN",
      "name": "My Ad Tracker",
      "type": "ad",
      "mode": "302"
    },
    "id": 2
  }'
```

**Create Campaign:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.campaign.create",
    "params": {
      "admin_token": "TOKEN",
      "tracker_id": "TRACKER_ID",
      "name": "Summer Campaign"
    },
    "id": 3
  }'
```

**Create Channel:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.channel.create",
    "params": {
      "admin_token": "TOKEN",
      "tracker_id": "TRACKER_ID",
      "campaign_id": "CAMPAIGN_ID",
      "name": "Facebook Ads",
      "source": "facebook",
      "medium": "cpc"
    },
    "id": 4
  }'
```

**Create Target:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.target.create",
    "params": {
      "admin_token": "TOKEN",
      "tracker_id": "TRACKER_ID",
      "url": "https://example.com/landing"
    },
    "id": 5
  }'
```

**Generate Tracking Token:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.token.generate",
    "params": {
      "admin_token": "TOKEN",
      "tracker_id": "TRACKER_ID",
      "campaign_id": "CAMPAIGN_ID",
      "channel_id": "CHANNEL_ID",
      "target_id": "TARGET_ID",
      "mode": "302",
      "exp_seconds": 86400
    },
    "id": 6
  }'
```

**Other admin methods:** `admin.tracker.list`, `admin.tracker.update`, `admin.tracker.delete`, `admin.campaign.list`, `admin.channel.list`, `admin.channel.batchImport`, `admin.target.list`, `admin.site.create`, `admin.site.list`

### Track Methods

- `track.collectClick` — encrypted click tracking via JSON-RPC
- `track.collectEvents` — encrypted event batch via JSON-RPC

Both require hybrid RSA+AES encrypted payloads. See the JS SDK or the encryption example below.

## HTTP Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/rpc` | POST | JSON-RPC 2.0 endpoint |
| `/r/:token` | GET | 302 redirect click tracking |
| `/t/:token` | GET | JS-based click tracking page |
| `/sdk/track.js` | GET | Web analytics JS SDK |
| `/public-keys.json` | GET | RSA public key for encryption |

### Testing 302 Redirect

```bash
# Get the token from admin.token.generate, then:
curl -v "http://localhost:8080/r/TOKEN_HERE"
# Should return 302 redirect to the target URL
```

### Testing JS Track Page

```bash
curl "http://localhost:8080/t/TOKEN_HERE"
# Returns HTML page with embedded JS that collects browser env and reports via encrypted JSON-RPC
```

## Web Analytics SDK

**Create a site:**

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "admin.site.create",
    "params": {
      "admin_token": "TOKEN",
      "name": "My Website",
      "domain": "example.com"
    },
    "id": 7
  }'
```

**Add the SDK to your page:**

```html
<script src="http://localhost:8080/sdk/track.js"></script>
<script>
  TrackSDK.init("SITE_KEY_FROM_RESPONSE");
  // Pageview is tracked automatically on init
  // Track custom events:
  TrackSDK.trackEvent("button_click", { button_id: "signup" });
</script>
```

## Frontend JS Encryption Example

```html
<!DOCTYPE html>
<html>
<head><title>Encryption Test</title></head>
<body>
<script>
async function testEncryptedRPC() {
  // 1. Fetch the public key
  const keysResp = await fetch("/public-keys.json");
  const keys = await keysResp.json();

  // 2. Import RSA public key (SPKI)
  const pemContents = keys.public_key
    .replace("-----BEGIN PUBLIC KEY-----", "")
    .replace("-----END PUBLIC KEY-----", "")
    .replace(/\n/g, "");
  const binaryDer = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0));
  const pubKey = await crypto.subtle.importKey("spki", binaryDer.buffer,
    { name: "RSA-OAEP", hash: "SHA-256" }, false, ["encrypt"]);

  // 3. Generate random 32-byte AES key
  const dataKey = crypto.getRandomValues(new Uint8Array(32));

  // 4. RSA-OAEP encrypt the AES key
  const ek = await crypto.subtle.encrypt({ name: "RSA-OAEP" }, pubKey, dataKey);

  // 5. AES-GCM encrypt the payload
  const aesKey = await crypto.subtle.importKey("raw", dataKey, "AES-GCM", false, ["encrypt"]);
  const nonce = crypto.getRandomValues(new Uint8Array(12));
  const payload = { site_key: "YOUR_SITE_KEY", visitor_id: "test", session_id: "test",
    events: [{ type: "pageview", url: "https://example.com", title: "Test", referrer: "" }] };
  const ct = await crypto.subtle.encrypt({ name: "AES-GCM", iv: nonce }, aesKey,
    new TextEncoder().encode(JSON.stringify(payload)));

  // 6. Base64 encode and build JSON-RPC request
  function b64(buf) {
    return btoa(String.fromCharCode(...new Uint8Array(buf)));
  }

  const resp = await fetch("/rpc", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      jsonrpc: "2.0",
      method: "track.collectEvents",
      params: {
        ek: b64(ek), nonce: b64(nonce), ct: b64(ct),
        ts: Math.floor(Date.now() / 1000),
        nonce2: crypto.randomUUID(),
        kid: keys.kid
      },
      id: "test-1"
    })
  });
  console.log(await resp.json());
}
testEncryptedRPC();
</script>
</body>
</html>
```

## Architecture Notes

- **Events table**: Consider partitioning by month for production workloads
- **RSA keys**: Auto-generated on first startup; back up `keys/` directory
- **Rate limiting**: Redis-based sliding window, fail-open on Redis errors
- **Bot detection**: UA-based scoring with configurable thresholds
- **Click dedup**: Redis SETNX with configurable TTL window
