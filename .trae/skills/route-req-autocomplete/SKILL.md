---
name: "route-req-autocomplete"
description: "Generates route Req* handlers from examples/protos/XX.go. Invoke when user says '新增XX路由' or when completing Req* methods in examples/route."
---

# Route Req Autocomplete

Use this skill to auto-generate route handlers in `examples/route` based on `Req*` structs in `examples/protos/XX.go`.

## When To Invoke

- User says `新增XX路由` or `补充XX路由` `.
- User asks to generate all Req handlers for a route module.
- You are completing `Req*` methods in `examples/route`.

## Command Intent

Input command:

- `新增skin路由`
- `新增equip路由`

Expected behavior:

1. Locate `examples/protos/XX.go`.
2. Find all structs named `Req*` in that file.
3. For each `ReqXxx`, infer paired response `ResXxx`.
4. Open or create `examples/route/XX.go`.
5. Generate missing route methods with project style.
6. Do not duplicate existing methods.

## Target Method Pattern

Reference style (`signin.go`):

```go
func (ps *SignInRoute) ReqSignIn(s *network.Session, index int32, msg *protos.ReqSignIn) *protos.ResSignIn {
	player := player.GetPlayerService().GetPlayerBySession(s)
	err := ps.service.SignIn(player)
	if err != nil {
		return &protos.ResSignIn{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignIn{}
}
```

Generated style (`ReqXxx` -> `ResXxx`):

```go
func (r *XxRoute) ReqXxx(s *network.Session, index int32, msg *protos.ReqXxx) *protos.ResXxx {
	p := player.GetPlayerService().GetPlayerBySession(s)
	code := r.service.Xxx(p, msg)
	return &protos.ResXxx{
		Code: code,
	}
}
```

## Generation Rules (Strict)

1. Method name must start with `Req`.
2. Signature should use:
   - `s *network.Session`
   - `index int32`
   - `msg *protos.ReqXxx`
   - return `*protos.ResXxx`
3. Resolve player from session first using:
   - `p := player.GetPlayerService().GetPlayerBySession(s)`
4. Delegate to module service (`r.service`) using existing naming conventions.
5. Response mapping priority:
   - If service returns `int32` code: fill `Code`.
   - If service returns `*protos.ResXxx`: return directly.
   - If service returns `BusinessError`: map to `Code`.
6. Keep style aligned with existing module files (`equip.go`, `signin.go`, `rune.go`).
7. Keep imports minimal and correct (`network`, `protos`, module service, player service, optional event/context imports).

## File-Level Behavior

- If `examples/route/XX.go` does not exist:
  - Create `XXRoute`, `NewXXRoute`, `Init()`, `service` field.
- If file exists:
  - Preserve current structure and only append missing `Req*` methods.
- Always run formatting-compatible output.

## Output Checklist

- All `Req*` in `examples/protos/XX.go` are covered by route methods.
- Existing methods are not duplicated.
- Code compiles against current module service/protos contracts where possible.
