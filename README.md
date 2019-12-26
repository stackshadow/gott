# GOTT

[![Go Report Card](https://goreportcard.com/badge/github.com/oimyounis/gott)](https://goreportcard.com/report/github.com/oimyounis/gott)
[![GoDoc](https://godoc.org/github.com/oimyounis/gott?status.svg)](https://godoc.org/github.com/oimyounis/gott)
![GitHub](https://img.shields.io/github/license/oimyounis/gott)

GOTT is a MQTT Broker written in pure Go. Aims to be a high-performance pluggable broker with most of the features that could be embedded in a broker available out of the box.  
  
It needs a lot of optimizations and improvements to be functional and usable for production.  
**Hopefully with your contribution we could build something great!**

## Project Status
### Under Development
- GOTT is currently in a very early stage and is under heavy development. It is not stable or complete enough to be used in production.
- Currently implementing the MQTT v3.1.1 spec.

### Planned for v1
- [x] Ping (client -> server)
- [x] Topic filtering with wildcards support
- [x] Subscriptions
- [x] QoS 0 messages
- [x] QoS 1 messages
- [x] QoS 2 messages
- [x] Unsubscribe a client when it disconnects (finished but still needs optimizations)
- [x] Retained messages
- [x] Will messages
- [x] Sessions
- [ ] Plugins

### Planned for v2
- [ ] MQTT v5
- [ ] TLS
- [ ] Clustering (maybe)
- [ ] WebSockets

### Known Issues
- Restarting the broker will reset existing subscriptions. They are not saved to disk (sessions are saved to disk but subscriptions are not restored on broker restart).

## Installation
1. Run `go get -u github.com/google/uuid`
2. Run `go get -u github.com/dgraph-io/badger`
3. Run `go get -u github.com/json-iterator/go`
4. Clone/download this repo and build/run `main/main.go`

## License
Apache License 2.0, see LICENSE.

## Contribution
You are very welcome to submit a new feature, fix a bug, an optimization to the code, report a bug or even a benchmark would be helpful.  
### To Contribute:  
Open an issue or:
1. Fork this repo.
2. Create a new branch with a descriptive name (example: *feature/some-new-function* or *fix/something-somewhere*).
3. Commit and push your code to your new branch.
4. Create a new Pull Request here.  
