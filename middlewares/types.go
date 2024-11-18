package middlewares

import (
	"github.com/DaHuangQwQ/gweb/internal/types"
)

type Middleware func(next types.HandleFunc) types.HandleFunc
