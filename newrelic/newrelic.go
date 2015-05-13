// New Relic metrics reporter middleware.

package newrelic

import (
	"fmt"
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
