# `schema` Creates a json schema from a struct

This package can generate a json schema for a struct that is compatible with https://json-schema.org/

Supported struct tags:
- `json:`
    - `"-"` Ignores the field
    - `"other_name"` Renames the field
- `jsonSchema:`
    - `"notRequired"` Set the field are not required (by default all fields with the exeption of `ptr`, `array`, `slice` and `map` are set as required)
    - `"required"` Set the field as required
