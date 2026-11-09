# beatfx

[![Go Reference](https://pkg.go.dev/badge/github.com/uchaloop/beatfx.svg)](https://pkg.go.dev/github.com/uchaloop/beatfx) [![CI](https://github.com/uchaloop/beatfx/actions/workflows/ci.yml/badge.svg)](https://github.com/uchaloop/beatfx/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/v/tag/uchaloop/beatfx?label=release)](https://github.com/uchaloop/beatfx/tags) [![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Run recurring or one-shot beat work with Uber Fx.**
The beat module owns execution, scheduling and cluster assignment. Beatfx owns
Fx wiring and lifecycle. Configuration loading through confx is optional.

[Install](#installation) · [Choose an adapter](#quick-start) · [Reference](#reference)

## Installation

Requires Go 1.27.0 or later and uses beat v0.8.0.

```sh
go get github.com/uchaloop/beatfx@latest
```

The module includes the scheduler adapter and the jobfx subpackage for one-shot
execution. Both use `github.com/uchaloop/beat/job` types. Types from the former
standalone job module are not interchangeable with these types.

## Quick start

| Scenario | Entry point | Guide |
|---|---|---|
| Long-lived scheduler | `beatfx.Module()` inside `fx.New(...).Run()` | [Scheduler adapter](docs/beatfx/README.md) |
| Kubernetes Job / CronJob or another one-shot process | `jobfx.Execute(...)` | [One-shot execution](docs/jobfx/README.md) |
| Manual one-shot lifecycle | `jobfx.MakeApp` or `jobfx.New`, then Start/Stop | [App lifecycle](docs/jobfx/README.md#run-the-process) |

Import `github.com/uchaloop/beatfx/jobfx` for one-shot helpers. It is a package
inside this module, not the former standalone jobfx module. Execute accepts
ordinary Fx options, a lazy confx.App, or a configured jobfx.App; it exits only
after cleanup. Use manual Start/Stop when caller-owned cleanup or error handling
is needed. For a daemon, use Fx.Run or confx.New(...).Run with beatfx.Module.

## Reference

- [beat](https://github.com/uchaloop/beat): execution and scheduling contracts
- [API reference](https://pkg.go.dev/github.com/uchaloop/beatfx)
- [Scheduler example](example_test.go) and [one-shot example](jobfx/example_test.go)
- [CHANGELOG](CHANGELOG.md)

## License

[MIT](LICENSE)
