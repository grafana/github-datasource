package schema

import (
	_ "embed"
)

//go:embed config.json
var ConfigSchemaJSON []byte
