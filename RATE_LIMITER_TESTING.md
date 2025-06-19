# Rate Limiter Testing Guide

This guide explains how to test the rate limiter functionality using autocannon.

## Prerequisites

1. **Node.js and npm** - For running autocannon
2. **Go API server** - Your DVD rental API must be running

## Setup

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Start your Go API server:**
   ```bash
   go run cmd/api/main.go
   ```
   
   Make sure your server is running on `localhost:8080` (or update the URLs in the test scripts).

## Rate Limiter Configuration

The rate limiter is configured with these default values:
- **Requests per time frame**: 100 (configurable via `RATE_LIMITER_REQUEST_PER_TIME_FRAME` env var)
- **Time frame**: 60 seconds (configurable via `RATE_LIMITER_TIME_FRAME` env var)

## Running Tests

### Quick Burst Test
For a quick test to see if rate limiting is working:
```bash
npm run test:ratelimiter:burst
```

### Comprehensive Test Suite
For a full test suite with multiple scenarios:
```bash
npm run test:ratelimiter
```

## Test Scenarios

### 1. Normal Load Test
- **Duration**: 5 seconds
- **Connections**: 10
- **Expected**: All requests should succeed (2xx responses)

### 2. High Load Test
- **Duration**: 10 seconds
- **Connections**: 20
- **Expected**: Some requests should be rate limited (4xx responses)

### 3. Burst Load Test
- **Duration**: 5 seconds
- **Connections**: 50
- **Expected**: Many requests should be rate limited (4xx responses)

### 4. Sustained Load Test
- **Duration**: 65 seconds (slightly longer than the time frame)
- **Connections**: 15
- **Expected**: Should see rate limiting behavior over time

## Understanding Results

- **2xx responses**: Successful requests (not rate limited)
- **4xx responses**: Rate limited requests (should be 429 Too Many Requests)
- **5xx responses**: Server errors (should be 0)

## Customizing Tests

You can modify the test parameters in the JavaScript files:

- `connections`: Number of concurrent connections
- `duration`: Test duration in seconds
- `pipelining`: Number of pipelined requests per connection

## Troubleshooting

### Server Not Responding
If you get "Server not responding" errors:
1. Make sure your Go API server is running
2. Check that it's listening on `localhost:8080`
3. Verify the `/v1/test` endpoint is accessible

### No Rate Limiting Detected
If you don't see any 4xx responses:
1. Check your rate limiter configuration
2. Ensure the rate limiter middleware is properly applied
3. Verify the rate limit values are set correctly

### Rate Limiting Too Aggressive
If too many requests are being blocked:
1. Increase the `RATE_LIMITER_REQUEST_PER_TIME_FRAME` value
2. Increase the `RATE_LIMITER_TIME_FRAME` value
3. Restart your server after changing environment variables

## Environment Variables

You can customize the rate limiter behavior by setting these environment variables:

```bash
export RATE_LIMITER_REQUEST_PER_TIME_FRAME=50  # Lower for stricter limiting
export RATE_LIMITER_TIME_FRAME=30s             # Shorter time frame
```

## Manual Testing

You can also test manually using curl:

```bash
# Test normal request
curl http://localhost:8080/v1/test

# Test rate limiting with multiple rapid requests
for i in {1..150}; do
  curl -w "%{http_code}\n" http://localhost:8080/v1/test
done
```

Look for 429 status codes in the output to confirm rate limiting is working. 