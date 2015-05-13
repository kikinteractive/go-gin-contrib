package newrelic

import (
	"fmt"
	"net/http"

	"github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/newrelic_platform_go"
)

// Status counters registry.
var statusCounters map[int]metrics.Counter

// New metrica collector - counter per each http status code.
type counterByStatusMetrica struct {
	counter metrics.Counter
	name    string
	units   string
}

// metrics.IMetrica interface implementation.
func (m *counterByStatusMetrica) GetName() string { return m.name }

func (m *counterByStatusMetrica) GetUnits() string { return m.units }

func (m *counterByStatusMetrica) GetValue() (float64, error) { return float64(m.counter.Count()), nil }

// addHTTPStatusMetrics initializes counter metrics for all http statuses and adds them to the component.
func addHTTPStatusMetrics(component newrelic_platform_go.IComponent) {
	httpStatuses := []int{
		http.StatusContinue, http.StatusSwitchingProtocols,

		http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNonAuthoritativeInfo,
		http.StatusNoContent, http.StatusResetContent, http.StatusPartialContent,

		http.StatusMultipleChoices, http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther,
		http.StatusNotModified, http.StatusUseProxy, http.StatusTemporaryRedirect,

		http.StatusBadRequest, http.StatusUnauthorized, http.StatusPaymentRequired, http.StatusForbidden,
		http.StatusNotFound, http.StatusMethodNotAllowed, http.StatusNotAcceptable, http.StatusProxyAuthRequired,
		http.StatusRequestTimeout, http.StatusConflict, http.StatusGone, http.StatusLengthRequired,
		http.StatusPreconditionFailed, http.StatusRequestEntityTooLarge, http.StatusRequestURITooLong, http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable, http.StatusExpectationFailed, http.StatusTeapot,

		http.StatusInternalServerError, http.StatusNotImplemented, http.StatusBadGateway,
		http.StatusServiceUnavailable, http.StatusGatewayTimeout, http.StatusHTTPVersionNotSupported,
	}
	statusCounters = make(map[int]metrics.Counter, len(httpStatuses))

	for _, statusCode := range httpStatuses {
		counter := metrics.NewCounter()
		statusCounters[statusCode] = counter
		component.AddMetrica(&counterByStatusMetrica{
			counter: counter,
			name:    fmt.Sprintf("http/status/%d", statusCode),
			units:   "count",
		})
	}
}

// incStatusCounter increments the corresponding status counter.
func incStatusCounter(status int) {
	statusCounters[status].Inc(1)
}
