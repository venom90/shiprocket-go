# Returns And NDR

## Returns

Covered operations:

- Create return order
- Create exchange order
- Update return order
- List return orders
- Return-specific serviceability
- Return-specific AWB assignment

## NDR

Covered operations:

- List all NDR records
- Fetch NDR by AWB
- Act on an NDR

Supported action constants:

- `ndr.ActionReattempt`
- `ndr.ActionRTO`
- `ndr.ActionReturn`

Runnable examples:

- [Create return order](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/create-return-order/main.go)
- [Act on NDR](/Users/tirumalrao/workspace/venom90/shiprocket-go/docs/examples/act-on-ndr/main.go)
