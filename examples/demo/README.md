# wfiber Demo (2-3 minutes)

This example is designed to showcase:

1. Go route types generate frontend SDK/types.
2. A backend shape change immediately surfaces TypeScript errors in frontend usage.

## Files

- Backend app: `examples/demo/backend/main.go`
- Generated client target: `examples/demo/frontend/src/api/api.gen.ts`
- Frontend usage sample: `examples/demo/frontend/src/demo.ts`
- Frontend HTTP adapter: `examples/demo/frontend/src/api/http.ts`

## 0) Install frontend tooling once

```bash
cd examples/demo/frontend
npm install
```

## 1) Start backend (generates SDK automatically)

From repo root:

```bash
go run ./examples/demo/backend
```

`wfiber` runs `CodeGen()` during `app.Listen`, so this creates/updates:

- `examples/demo/frontend/src/api/api.gen.ts`
- `examples/demo/backend/swagger/doc.json`

Keep the backend running in terminal 1.

Swagger docs are viewable in the browser at:

- `http://localhost:8081/swagger-ui` (Swagger UI)
- `http://localhost:8081/swagger/doc.json` (raw generated spec)

## 2) Show type-safe frontend usage

In terminal 2:

```bash
cd examples/demo/frontend
npm run typecheck
```

You should get a clean typecheck.

## 3) Trigger contract drift on purpose

Edit `examples/demo/backend/models/models.go`:

```go
type Todo struct {
	ID    string `json:"id"` // changed from int -> string
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
```

Restart backend:

```bash
go run ./examples/demo/backend
```

Now rerun typecheck in frontend:

```bash
cd examples/demo/frontend
npm run typecheck
```

Expected errors in `src/demo.ts` (old assumptions):

- `Type 'number' is not assignable to type 'string'.`
- Arithmetic on `payload.id` fails because it is now `string`.