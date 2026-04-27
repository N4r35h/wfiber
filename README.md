# WFiber

<video src="https://raw.githubusercontent.com/N4r35h/wfiber/refs/heads/main/examples/demo/wFiberDemoProjectShowCase.webm" controls muted playsinline width="100%"></video>
(not loading ? view the showcase video on youtube @ [https://youtu.be/lji_UkhNAs0](https://youtu.be/lji_UkhNAs0))

WFiber (**W**rapped **Fiber**) is a thin wrapper over [gofiber](https://github.com/gofiber/fiber) focused on code generation.

Its main goal is to make your Go backend structs the source of truth and generate a typed TypeScript SDK for frontend API calls.

## Why wfiber

- Define request and response shapes once in Go.
- Generate TypeScript interfaces and API functions automatically.
- Catch stale frontend usage at compile time when backend shapes change.

WFiber uses [gos2tsi](https://github.com/N4r35h/gos2tsi) (built on `golang.org/x/tools/go/packages`) to parse struct metadata and generate TS without reflection-based runtime scanning.

## Features overview

- Typed TypeScript SDK generation from route input/output types
- Generated interfaces for backend structs used by routes
- Generated API methods (`GetX`, `PostX`, `PutXById`, etc.)
- Swagger document generation from the same route/type metadata

## Quickstart

Install:

```bash
go get github.com/N4r35h/wfiber
```

Minimal setup (based on `wfiber/wfiber_test.go`):

```go
app := wfiber.New(wfiber.WFiberAppConfig{
    APIPrefix:              "/api",
    FrontendFolder:         "frontend_test",
    SwaggerDocFolder:       "swagger_test",
    GeneratedAPIClientPath: "/src/api/api.gen.ts",
    GenerateClient:         true,
})

api := app.Group("/api")
tests := api.Group("/tests")
tests.Post("/", models_test.Tests{}, structs_test.SimpleAPIResponseWithData[models_test.Tests]{}, handler)
```

Then run code generation:

- Explicitly via `app.CodeGen()`
- Or automatically when calling `app.Listen(...)`

## Quick demo setup

Use the ready-made demo in `examples/demo`:

```bash
go run ./examples/demo/backend
```

Then in another terminal:

```bash
cd examples/demo/frontend
npm install
npm run typecheck
```

For the full backend-shape-change -> frontend type-error walkthrough, see:

- `examples/demo/README.md`

## Example: generated TypeScript SDK

From `wfiber/frontend_test/src/api/api.gen.ts`:

```ts
export interface Tests {
    ID: number
    name: string
}

import http from './http'

export const PostTests = (_ip: Tests, query?: string): Promise<SimpleAPIResponseWithData<Tests>> => {
    return new Promise((resolve, reject) => {
        http.post('/tests' + (query || ''), _ip).then(response => resolve(response.data)).catch(reject)
    })
}
```

The generated client imports `./http`, which is your manually maintained HTTP adapter (`wfiber/frontend_test/src/api/http.ts`).

## Backend shape change demo (contract drift protection)

This is the practical value of generated interfaces.

### 1) Baseline backend shape

Current model (`wfiber/models_test/test.go`):

```go
type Tests struct {
    ID   int
    Name string `json:"name"`
}
```

Generated frontend shape:

```ts
export interface Tests {
    ID: number
    name: string
}
```

### 2) Change backend contract

Change backend field type:

```go
type Tests struct {
    ID   string
    Name string `json:"name"`
}
```

Regenerate the SDK.

### 3) Frontend code now flags old assumptions

Old frontend usage:

```ts
const payload: Tests = { ID: 42, name: "alpha" }
const doubled = payload.ID * 2
```

After regeneration, TypeScript reports errors because `ID` is now `string`:

- `Type 'number' is not assignable to type 'string'.`
- `The left-hand side of an arithmetic operation must be of type 'any', 'number', 'bigint' or an enum type.`

This is exactly how wfiber helps prevent silent backend/frontend drift.

## Recommended workflow

1. Update Go structs and route contracts in backend code.
2. Regenerate artifacts (`app.CodeGen()` or run your app/tests that invoke it).
3. Commit generated SDK updates (`api.gen.ts`) with backend contract changes.
4. Run frontend type-check (`npx tsc --noEmit`) to catch outdated usage.
5. Fix reported type errors and ship.

## Repository example paths

- Route + generation flow: `wfiber/wfiber.go`
- Route/type example setup: `wfiber/wfiber_test.go`
- Backend model: `wfiber/models_test/test.go`
- Generated SDK: `wfiber/frontend_test/src/api/api.gen.ts`
- Manual HTTP client used by generated SDK: `wfiber/frontend_test/src/api/http.ts`

## Notes

- Do not manually edit generated files like `api.gen.ts`.
- Treat generated SDK changes as part of your API contract change.
- Keeping Go structs as source of truth removes duplicate contract maintenance across backend and frontend.