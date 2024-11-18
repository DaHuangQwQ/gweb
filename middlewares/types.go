package middlewares

import (
	"github.com/DaHuangQwQ/gweb/types"
)

type Middleware func(next types.HandleFunc) types.HandleFunc
