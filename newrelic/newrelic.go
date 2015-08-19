// New Relic metrics reporter middleware.

package newrelic

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yvasiyarov/gorelic"
)

var agent *gorelic.Agent

// InitAgent creates and inits a newrelic agent.
func InitAgent(license, appName, hostName string, verbose bool) error {

	// Sanity checks.
	switch {
	case license == "":
		return fmt.Errorf("empty newrelic license")
	case appName == "":
		return fmt.Errorf("empty newrelic app name")
	case hostName == "":
		return fmt.Errorf("empty newrelic hostname")
	}

	agent = gorelic.NewAgent()

	// Init agent identification parameters (company license, application name and hostname).
	agent.NewrelicLicense = license
	agent.AgentGUID = appName
	agent.NewrelicName = hostName

	// Init http metric collectors.
	agent.CollectHTTPStat = true
	agent.CollectHTTPStatuses = true
	agent.Verbose = verbose

	// Add custom metrics.
	agent.AddCustomMetric(&openSocketsMetrica{})

	// Start agent.
	agent.Run()
	return nil
}

// Handler is a gin middleware handler. It wraps the default handler and updates metric timers and counters.
func Handler(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	if agent != nil {
		agent.HTTPTimer.UpdateSince(startTime)
		agent.HTTPStatusCounters[c.Writer.Status()].Inc(1)
	}
}

// Our custom metrics go here.

// openSocketsMetrica reports the number of open sockets by the current process.
type openSocketsMetrica struct{}

func (metrica *openSocketsMetrica) GetName() string {
	return "network/openSockets"
}

func (metrica *openSocketsMetrica) GetUnits() string {
	return "count"
}

func (metrica *openSocketsMetrica) GetValue() (float64, error) {
	return float64(countOpenSockets()), nil
}

// countOpenSockets finds out the number of open sockets in the current process
// using the system function lsof with some custom parameters:
// -a: AND filters
// -iTCP: select only TCP sockets
// -n: inhibits the conversion of network numbers to host names, speeds up output
// -P: nhibits the conversion of port numbers to port names, speeds up output
// -p <PID>: outputs only for this process
// -w: disables warnings
func countOpenSockets() int {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -a -iTCP -n -P -p %v -w | wc -l", os.Getpid())).Output()
	if err != nil {
		return 0
	}
	val := strings.TrimSpace(string(out))
	num, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return num
}
