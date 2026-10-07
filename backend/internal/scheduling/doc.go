// Package scheduling is the core bounded context: appointments, holds,
// cancellation rules and availability.
//
// Layout (created step by step):
//   - domain/            aggregates, value objects, events (stdlib only)
//   - app/command/       write side
//   - app/query/         read side
//   - adapters/postgres/ persistence
//   - adapters/connect/  RPC handlers
//   - projectors/        read models built from events
package scheduling

// TODO(step-1): create domain/
