# goutil [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil) [![Go Report](https://goreportcard.com/badge/github.com/flier/goutil)](https://goreportcard.com/report/github.com/flier/goutil)

`goutil` is a set of toolkits designed to simplify golang development.

- [arena](#arena) A simple memory arena allocator.
- [either](#either) The Either with variants Left and Right is a general purpose sum type with two cases.
- [opt](#opt) Optional values.
- [res](#res) Error handling with the Result type.
- [tuple](#tuple) A finite heterogeneous sequence, (T0, T1, ..).
- [untrust](#untrust) Safe, fast, zero-panic, zero-crashing, zero-allocation parsing of untrusted inputs.
- [xiter](#xiter) Provides utilities for enhanced iteration patterns and helpers.

## Package

### arena [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/arena)

Package `arena` provides a simple memory arena allocator for Go, inspired by the article [Cheating the Reaper in Go](https://mcyoung.xyz/2025/04/21/go-arenas/).

### either [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/either)

The Either with variants Left and Right is a general purpose sum type with two cases.

### opt [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/opt)

Optional values.

### res [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/res)

Error handling with the Result type.

### tuple [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/tuple)

A finite heterogeneous sequence, (T0, T1, ..).

### untrust [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/untrust)

Safe, fast, zero-panic, zero-crashing, zero-allocation parsing of untrusted inputs.

### xiter [![Go Reference](https://pkg.go.dev/badge/github.com/flier/goutil)](https://pkg.go.dev/github.com/flier/goutil/pkg/xiter)

Provides utilities for enhanced iteration patterns and helpers.