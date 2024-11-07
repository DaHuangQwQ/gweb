package gweb

type Middleware func(next HandleFunc) HandleFunc
