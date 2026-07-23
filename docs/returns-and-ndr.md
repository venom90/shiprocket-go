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

- [Create return order](examples/create-return-order/main.go)
- [Act on NDR](examples/act-on-ndr/main.go)
