package assets

import (
	"embed"
)

//go:embed css
//go:embed images
//go:embed html
//go:embed js
var Assets embed.FS
