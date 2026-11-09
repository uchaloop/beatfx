# Scheduler adapter

**Connect the beat scheduler to the Fx lifecycle.**

[Overview](../../README.md) · [Quick start](#quick-start) · [Lifecycle](#how-it-works) · [Configuration](#recommended-configuration)

See [installation requirements](../../README.md#installation).

## Quick start

```go
package main

import (
    "context"
    "time"

    "github.com/uchaloop/beat"
    "github.com/uchaloop/beatfx"
    "github.com/uchaloop/beat/job"
    "go.uber.org/fx"
)

func main() {
    fx.New(
        fx.Supply(beat.Config{Period: 5 * time.Minute, Jitter: 30 * time.Second}),
        fx.Provide(func() (*job.Runner, error) {
            return job.MakeRunner(
                job.Config{Timeout: 4 * time.Minute},
                func(ctx context.Context) (int64, error) {
                    // Replace with one batch of application work.
                    return 0, ctx.Err()
                },
            )
        }),
        beatfx.Module(beat.WithGracefulStop()),
        fx.StopTimeout(5 * time.Minute),
    ).Run()
}
```

The application can instead provide the Runner using `beatfx/jobfx.Module()`.
beatfx requires only the Runner value, not that specific adapter.

`Module` requires beat.Config and *job.Runner. beat.Handler and one ordered
beatfx.Options slice are optional. A container Handler is the default; explicit
WithHandler options override it. Static options precede container options.
Use one Module per application.

The complete program above belongs in an application module using beatfx. Run `GOWORK=off go run .`; attempts run every five
minutes with a replica offset. Press Ctrl+C to stop. The example returns zero
processed items and does not emit application result logs.

## How it works

```mermaid
flowchart LR
    A["app.Run"] --> B["OnStart → beat.Start"]
    B --> C["Scheduled attempts"]
    C -->|terminal loop error| D["Request Fx shutdown: code 1"]
    D --> E["OnStop → beat.Stop"]
    S["SIGINT / SIGTERM"] --> E
```

| Event | Behaviour |
|---|---|
| Attempt returns an error | Continue the schedule |
| Another cluster owns the point | Skip this point |
| Loop cannot continue | Request application shutdown with code 1 |
| Application shuts down | Stop scheduling; cancel or drain accepted work |

Startup starts the scheduler. Shutdown stops scheduling and drains or cancels
work according to beat's options. Ordinary attempt errors do not stop the daemon.
If the scheduling loop terminates unexpectedly, beatfx requests Fx shutdown with
exit code 1. app.Run handles this automatically. Applications driving Start/Stop
manually must receive app.Wait and call Stop themselves. The original loop error
is returned by beat.Stop through the shutdown hook.

## Documentation

[API reference on pkg.go.dev](https://pkg.go.dev/github.com/uchaloop/beatfx)

## Recommended configuration

> [!TIP]
> We recommend [confmaker](https://github.com/uchaloop/confmaker) for typed ENV
> configuration and [confx](https://github.com/uchaloop/confx) for its Fx integration.
> Configuration loading stays in the application; it is optional for the work libraries.

<details>
<summary><strong>Configure from ENV with confmaker / confx</strong></summary>

Replace `fx.Supply(beat.Config{...})` in the quick start with these Fx options:

```go
confx.Module(),
confx.Provide[beat.Config]("beat"),
```

Import `github.com/uchaloop/confx`. Include `confx.Module()` once, even when
registering several config types. Keep the work constructor and adapter Module.

`BEAT_PERIOD=5m BEAT_JITTER=30s` fills beat.Config. To load the Runner timeout
as well, add `confx.Provide[job.Config]("job")` and have its constructor accept
`cfg job.Config` instead of constructing that value inline.

</details>

## Related libraries

| Library | Purpose |
|---|---|
| [beat](https://github.com/uchaloop/beat) | Scheduling and shutdown contracts |
| [jobfx](../jobfx/README.md) | An optional way to provide the Runner |

## Acknowledgements

Thanks to the authors and maintainers of [Uber Fx](https://github.com/uber-go/fx)
for dependency injection and lifecycle primitives that make this adapter possible.

## License

[MIT](../../LICENSE)
