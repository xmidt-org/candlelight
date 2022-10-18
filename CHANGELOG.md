# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.0.12]
- Upgraded Go to v1.19 and Open Telemetry libraries to v1.11.0, updated README [#70](https://github.com/xmidt-org/candlelight/pull/70)

## [v0.0.11]
- Fixed malformed linter config [#62](https://github.com/xmidt-org/candlelight/pull/62)
- Updated tracer, relaxed return type in function within Tracer.go [#60](https://github.com/xmidt-org/candlelight/pull/60)

## [v0.0.10]
- Update to use otel v1.7.0. [#44](https://github.com/xmidt-org/candlelight/pull/44)
- Update to use use common workflows for simplicity. [#45](https://github.com/xmidt-org/candlelight/pull/45)
- Update the dependencies to the latest. [#46](https://github.com/xmidt-org/candlelight/pull/46)


## [v0.0.9]
- Moved GenTID function from transaction package to candlelight. [#38](https://github.com/xmidt-org/candlelight/pull/38)

## [v0.0.8]
- Added header key for the WebPA transaction UUID. [#37](https://github.com/xmidt-org/candlelight/pull/37)

## [v0.0.7]
- upgraded webpa-common to v2.0.4

## [v0.0.6]
- Update webpa-common to v2. [#34](https://github.com/xmidt-org/candlelight/pull/34)

## [v0.0.5]
- Add tracing factory. [#29](https://github.com/xmidt-org/candlelight/pull/29)

## [v0.0.4]
- Add changes to allow tracing feature be optional. [#25](https://github.com/xmidt-org/candlelight/pull/25)

## [v0.0.3]
### Changed
- Add struct to unmarshal tracing info. [#19](https://github.com/xmidt-org/candlelight/pull/19) thanks to @Sachin4403
- We will be using defaults fields in logs and http Response.  [#22](https://github.com/xmidt-org/candlelight/pull/22) thanks to @Sachin4403
### Fixed
- TraceMiddleware was adding TraceID and SpanID to the response due to which we were seeing same traceID and multiple spanID in response, Now it will be adding the TraceID and SpanID only when the remote span is not present in context. [#21](https://github.com/xmidt-org/candlelight/pull/21) thanks to @Sachin4403

## [v0.0.2]
- Added setup and middleware for application tracing using opentelemetry. [#16](https://github.com/xmidt-org/candlelight/pull/16) thanks to @Sachin4403

## [v0.0.1]
- Updated the project configuration
- Initial creation

[Unreleased]: https://github.com/xmidt-org/candlelight/compare/v0.0.12..HEAD
[v0.0.12]: https://github.com/xmidt-org/candlelight/compare/v0.0.11..v0.0.12
[v0.0.11]: https://github.com/xmidt-org/candlelight/compare/v0.0.10..v0.0.11
[v0.0.10]: https://github.com/xmidt-org/candlelight/compare/v0.0.9..v0.0.10
[v0.0.9]: https://github.com/xmidt-org/candlelight/compare/v0.0.8..v0.0.9
[v0.0.8]: https://github.com/xmidt-org/candlelight/compare/v0.0.7..v0.0.8
[v0.0.7]: https://github.com/xmidt-org/candlelight/compare/v0.0.6..v0.0.7
[v0.0.6]: https://github.com/xmidt-org/candlelight/compare/v0.0.5..v0.0.6
[v0.0.5]: https://github.com/xmidt-org/candlelight/compare/v0.0.4..v0.0.5
[v0.0.4]: https://github.com/xmidt-org/candlelight/compare/v0.0.3..v0.0.4
[v0.0.3]: https://github.com/xmidt-org/candlelight/compare/v0.0.2..v0.0.3
[v0.0.2]: https://github.com/xmidt-org/candlelight/compare/v0.0.1..v0.0.2
[v0.0.1]: https://github.com/xmidt-org/candlelight/compare/v0.0.0..v0.0.1
