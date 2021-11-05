# A testing DB

This database is mainly for running within tests, it implements `db.Connection`

To understand the code it's best to read the `*_test.go` files

## Filters

Supported:

- Key has value filters (`{foo: 'bar'}`)
- Numeric filters $gt, $gte, $lt, $lte
- Logical filters $eq, $ne, $not, $and, $or
- Array filters $size
- Others $type _(only some types)_

Noteworthy Unsupported:

- Nested keys ({'foo.bar.bas': 'example'})
- Array filters $in
