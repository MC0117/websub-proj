# WebSub Hub

A simple WebSub hub written in Go. Handles subscriptions with intent verification and distributes published content to subscribers with HMAC-signed payloads.

## Running

Usig Docker

```
docker compose up --build
```

The hub runs on port 8080. A test client (`modfin/websub-client`) is included in the compose file and runs on port 8081.

## Endpoints

### POST /subscribe

Accepts form-encoded data with the following fields:

- `hub.callback` - subscriber's callback URL
- `hub.secret` - shared secret for HMAC signing
- `hub.topic` - topic to subscribe to
- `hub.mode` - should be `subscribe`

The hub sends a verification GET request to the callback URL with a challenge parameter. If the callback responds with 200, the subscription is stored.

### POST /publish

Accepts a JSON body and a `hub.topic` query parameter. The hub forwards the payload only to subscribers registered for that topic, with the body signed using HMAC-SHA256. The signature is sent in the `X-Hub-Signature` header. Outbound requests have a 15 second timeout to avoid hanging on unresponsive subscribers.

## Testing manually

Publish a message (with the hub running):

```
curl -X POST "http://localhost:8080/publish?hub.topic=/a/topic" -H "Content-Type: application/json" -d '{"message": "hello world"}'
```

## Implementation notes

- The subscriber store uses `sync.RWMutex` to handle concurrent reads and writes safely, since raceconditions are possible here
- Delivery to subscribers is done in parallel using goroutines. Errors from individual deliveries are logged rather than silently dropped.
- Subscribe requests are validated before processing, empty required fields are rejected with a 400.
- Intent verification follows the WebSub spec: the hub sends a challenge to the callback URL and only stores the subscription if the subscriber responds with 200.
- HMAC-SHA256 signing uses the subscriber-provided secret, so each subscriber can independently verify payload authenticity.

## Extra implementations beyond Case Req
- Support for multi-topic publishing. Only the subscribers that are subscribing to a certain topic will receive updates
- Outbound delivery requests have a 15 second timeout to prevent goroutine leaks from unresponsive subscribers.

## Bonus / future improvements

- Unsubscribe support via `hub.mode=unsubscribe`
- Request context propagation for outbound calls
- Persistent storage using Postgres or Redis (the store is in-memory only right now)

The reason these were not implemented is because control over the Websub client is restricted and the testing procedure had to be used differently (therefore in risk of not following the case instructions).

## Project structure

```
hub/
  cmd/hub/          - entrypoint
  internal/
    handlers/       - HTTP handlers for subscribe and publish
    delivery/       - outbound payload delivery with HMAC signing
    subscription/   - in-memory subscriber store (mutex-protected)
```

