package prometheus

import (
	"github.com/DaHuangQwQ/gweb"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func NewMiddlewareBuilder(namespace string, subsystem string, name string, help string) *MiddlewareBuilder {
	return &MiddlewareBuilder{Namespace: namespace, Subsystem: subsystem, Name: name, Help: help}
}

func (m *MiddlewareBuilder) Build() gweb.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Namespace: m.Namespace,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})

	prometheus.MustRegister(vector)

	return func(next gweb.HandleFunc) gweb.HandleFunc {
		return func(ctx *gweb.Context) {
			startTime := time.Now()

			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()

				pattern := ctx.MatchedRoute
				if pattern == "" {
					pattern = "unknown"
				}

				vector.WithLabelValues(pattern, ctx.Req.Method, strconv.Itoa(ctx.RespStatusCode)).Observe(float64(duration))
			}()

			next(ctx)
		}
	}
}
