# `schema` Creates a json schema from a struct

This package can generate a json schema for a struct that is compatible with https://json-schema.org/

The schemas generated are mainly used to generate OpenAPI v3 documentation at /api/v1/schema/openAPI

Default applied rules:

- A struct field is labeled as required when the data cannot be nil so `strings`,`bool`,`int`,`float`,`struct`, etc.. are required and types like `[]string`, `[8]int`, `*int`, `map[string]string` are not required. You can overwrite this behavior by using `jsonSchema` struct tag

Supported struct tags:

- `json:`
  - `"-"` Ignores the field
  - `"other_name"` Renames the field
- `jsonSchema:`
  - `"notRequired"` Set the field are not required (by default all fields with the exeption of `ptr`, `array`, `slice` and `map` are set as required)
  - `"required"` Set the field as required
  - `"deprecated"` Mark the field as deprecated
  - `"uniqueItems"` Every array entry must be unique _(Only for arrays)_
  - `"hidden"` Do not expose field in the schema
  - `"min=123"` Set the minimum value or array length for the field
  - `"max=123"` Set the maximum value or array length for the field

You can also chain jsonSchema tags using `,`, for example: `jsonSchema:"notRequired,deprecated"`
