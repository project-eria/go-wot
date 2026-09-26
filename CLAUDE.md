# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Go library implementing the W3C Web of Things (WoT) Thing Description v1.1 specification (`github.com/project-eria/go-wot`). Provides both producer (expose Things) and consumer (consume Things) sides with pluggable protocol bindings (HTTP via Fiber, WebSocket).

## Commands

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./dataSchema/
go test ./tests/

# Run a single test
go test ./dataSchema/ -run TestBooleanValidate

# Run tests with verbose output
go test -v ./...

# Build check (library, no main)
go build ./...
```

## Architecture

### Layered Design

```
thing/              Core Thing Description model (W3C TD v1.1)
dataSchema/         Data validation types (Boolean, Integer, Number, String, Object, Array)
interaction/        Interaction Affordances (Property, Action, Event, Form)
securityScheme/     Security definitions (currently only NoSecurity)
├── producer/       Server-side: expose Things with handlers
├── consumer/       Client-side: consume Things, read/write/observe
├── protocolHttp/   HTTP protocol binding (Fiber v2)
├── protocolWebSocket/  WebSocket binding (extends HTTP, for observe/events)
└── protocolDnssd/  DNS-SD discovery binding (_wot._tcp): Advertiser (a
                    ProtocolServer announcing exposed Things via mDNS) +
                    Browse (consumer-side discovery). TXT records td/type/
                    scheme validated per W3C WoT Discovery. Browse rejects
                    private addresses by default (BrowseOptions.AllowPrivate).
```

### Key Patterns

**Options Pattern** — All constructors use functional options:
```go
dataSchema.NewString(dataSchema.StringMinLength(5), dataSchema.StringMaxLength(100))
interaction.NewProperty("key", "Title", "Desc", schema, interaction.PropertyReadOnly(), interaction.PropertyObservable())
interaction.NewAction("key", "Title", "Desc", interaction.ActionInput(inputSchema))
```
Nil options are safely ignored.

**Protocol Binding Interfaces** — `producer.ProtocolServer` and `consumer.ProtocolClient` abstract protocol details. HTTP and WebSocket implementations are provided; new protocols plug in via these interfaces.

**Handler Pattern (Producer)** — Properties, actions, and events each have typed handler interfaces:
- `PropertyReadHandler` / `PropertyWriteHandler` / `PropertyObserveHandler`
- `ActionHandler`
- `EventSubscriptionHandler` / `EventListenerHandler`

Handlers receive `map[string]interface{}` for URI variables and return typed values.

**DataSchema Interface** — All schema types implement three methods: `FromString()`, `Validate()`, `IsNullOrEmpty()`. The `SimpleType` constraint (`bool | int | float64 | string`) is used for generics in typed schemas.

### JSON Marshalling

Custom `MarshalJSON`/`UnmarshalJSON` throughout for W3C compliance: `@context` can be string or array, `security` can be string or array, forms include protocol-specific supplements, etc.

### Concurrency

Producer uses `sync.RWMutex` for thread-safe access to exposed things, properties, and observer/event connections.

## Testing

- **Unit tests**: Per-package in `*_test.go` files (testify assertions)
- **Integration tests**: In `tests/` package, use `httpexpect` against a real Fiber server with HTTP + WebSocket bindings
- **Mocks**: In `mocks/` directory, generated with mockery

## Dependencies

- **HTTP server**: `gofiber/fiber/v2`
- **WebSocket**: `gofiber/websocket/v2` (server) + `gorilla/websocket` (client)
- **mDNS/DNS-SD**: `brutella/dnssd` (responder answers host A queries and probes for name conflicts; Browse runs in 10 s rounds, see `dnssdBrowser.go`)
- **Logging**: `rs/zerolog`
- **Testing**: `stretchr/testify` + `gavv/httpexpect/v2`
