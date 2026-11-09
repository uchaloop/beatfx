# Changelog

## [0.2.0] - 2026-11-09

### Breaking changes

- Use the shared `github.com/uchaloop/beat/assignment` package for cluster selection.
- Require Go 1.27.0 or later, matching beat v0.8.0.

- Use `github.com/uchaloop/beat/job` types throughout the scheduler adapter.
  Existing providers of standalone job types must migrate their imports.
- Require beat v0.8.0; remove the standalone job dependency.

### Fixed

- Copy static scheduler options when Module is created so caller slice changes
  cannot replace the captured configuration; referenced option state stays shared.
- Correct the assignment package path in GoDoc and align internal names.

### Changed

- Split scheduler and one-shot guides under docs with a compact root README.
- Gate release publication on reusable CI including confx integration; validate
  tags and changelog sections before publication and query Go proxy afterwards.

### Added

- Move jobfx v0.2.1 into `github.com/uchaloop/beatfx/jobfx`, preserving Module,
  lazy App Start/Stop, Execute, CLI handling and multicluster assignment.
- Preserve 30s/10s/30s lifecycle defaults, silent Execute and exit after cleanup.
- Include one-shot tests and separate confx v0.7.0 integration checks.

## [0.1.2] - 2026-10-01

- Update dependencies

## [0.1.1] - 2026-09-22

- Update README examples, executable examples and tests for job's `int64`
  processed-item counts. DI wiring and lifecycle behavior are unchanged.

## [0.1.0] - 2026-09-20

- Initial standalone release of the beatfx Uber Fx adapter.
- Independent module, dependencies, tests and release workflow.

[0.2.0]: https://github.com/uchaloop/beatfx/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/uchaloop/beatfx/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/uchaloop/beatfx/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/uchaloop/beatfx/releases/tag/v0.1.0
