## Summary

A simple To Do REST API using [Go](https://go.dev/), [GORM](https://gorm.io/) and [Gin](https://gin-gonic.com/).  In theory, should be
compatible with (or at least close to compatible with) the SPA from my other [ToDo learnings](https://github.com/csetera/todo-full-stack/tree/master/todo-client).

In addition, this implementation is a testbed for a couple of additional technologies:
* Role-based access control (RBAC) implemented using OIDC and [Zitdadel](https://zitadel.com/)
* [OpenTelemetry](https://opentelemetry.io/docs/languages/go/) for application instrumentation
* [Signoz](https://signoz.io/) in Docker for telemetry collection and visualization
