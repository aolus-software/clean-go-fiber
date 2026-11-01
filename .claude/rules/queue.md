# Rule: Queue (asynq)

Applies to `internal/queue/` (client, `tasks/`) and `internal/worker/` (server, `handlers/`).

## One task, one file each side

| Step | Where |
| --- | --- |
| task type constant + payload struct + `New<Name>Task` | `internal/queue/tasks/<name>.go` |
| handler | `internal/worker/handlers/<name>.go` |
| registration on the mux | `internal/worker/server.go` |
| enqueue | a service, via the `*asynq.Client` built in `internal/queue/client.go` |

Task types are constants (`TypeEmailSend = "email:send"`), never raw strings. Payloads are
JSON-serialisable structs.

## Handlers

- Return the error on failure so asynq retries; log both success and failure with zap.
- Never log secrets or token-bearing content from the payload.
- Handlers are idempotent where they can be — a retry must not double an effect that matters.

## Enqueueing

- Every enqueue sets `asynq.MaxRetry(n)` explicitly.
- Check and log the enqueue error; never discard it.
- Enqueue after the database transaction commits.
- No goroutines in request handlers — background work goes through the queue.

The worker and the API share Redis settings (`REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`,
`REDIS_DB`) from `internal/config/redis.go`.
