# WoT implementation in Golang

Based on Web of Things (WoT) Thing Description v1.1
https://w3c.github.io/wot-thing-description/

## Check dns-ds
# macOS
dns-sd -B _wot._tcp local.
dns-sd -L "<instance>" _wot._tcp local.

# Linux (Avahi)
avahi-browse -r _wot._tcp