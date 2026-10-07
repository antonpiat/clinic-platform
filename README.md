# clinic-platform

Appointment booking platform for clinics. Go modular monolith with DDD,
CQRS, transactional outbox and Kafka, plus a Next.js frontend.

> Work in progress.

## Quick start

Requirements: Go 1.27+, Docker, [Task](https://taskfile.dev),
[golangci-lint](https://golangci-lint.run) v2.

```bash
cp .env.example .env
task up        # postgres, redpanda, console at http://localhost:8080
task ci        # lint, vet, build, test
task todo      # remaining TODOs by step
```

## Architecture

Decisions and their reasoning live in [docs/adr](docs/adr).