# WoT implementation in Golang

Based on Web of Things (WoT) Thing Description v1.1
https://w3c.github.io/wot-thing-description/

## DNS-SD discovery (`protocolDnssd`)

Implements [W3C WoT Discovery](https://www.w3.org/TR/wot-discovery/#introduction-dns-sec)
over mDNS/DNS-SD, service type `_wot._tcp`, TXT records `td` (absolute TD
path, required), `type` (`Thing`|`Directory`), `scheme` (`http`|`https`).

### Announce (producer side)

The `Advertiser` is a regular `producer.ProtocolServer`: add it to the same
Producer as the `HttpServer` (which is the one actually serving the TDs).

```go
advertiser, _ := protocolDnssd.NewAdvertiser(protocolDnssd.Options{Port: 8888})
producer.AddServer(httpServer)
producer.AddServer(advertiser) // announces every exposed Thing
```

### Discover (consumer side)

```go
entries, _ := protocolDnssd.Browse(ctx, protocolDnssd.BrowseOptions{
    PreferIPv4:   true,
    AllowPrivate: true, // required on a LAN: private addresses are rejected by default (SSRF)
})
for entry := range entries {
    td, _ := http.Get(entry.TDURL.String())
    // consumer.Consume(...)
}
```

### Networks

Several independent houses can share a LAN: `Options.Network` and
`BrowseOptions.Network` name the network, which selects the service type
(`ServiceType`). The default network (`""`) is the standard `_wot._tcp`;
a named one uses `_wot-<network>._tcp` (1 to 11 characters among
`[a-z0-9-]`). Advertisers and browsers of different networks never see each
other.

```go
protocolDnssd.NewAdvertiser(protocolDnssd.Options{Port: 8888, Network: "dummy"}) // _wot-dummy._tcp
```

Notes:
- A process that both announces and browses sees its own announcements —
  filter by instance name.
- Instance-name conflicts are not resolved (RFC 6762 §9): keep TD ids unique.
- On macOS, a third-party mDNS resolver may not hear another third-party
  responder *on the same host* (mDNSResponder owns port 5353). Cross-host
  discovery — the normal topology — is unaffected, and `dns-sd` always sees
  the announcements.

## Check dns-sd
# macOS
dns-sd -B _wot._tcp local.
dns-sd -L "<instance>" _wot._tcp local.
dns-sd -B _wot-dummy._tcp local.   # a named network

# Linux (Avahi)
avahi-browse -r _wot._tcp
