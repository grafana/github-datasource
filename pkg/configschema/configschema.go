package configschema

import (
	_ "embed"
)

//go:embed config_schema.json
var ConfigSchemaJSON []byte
