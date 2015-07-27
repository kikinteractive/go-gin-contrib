package expvar

import (
	"encoding/json"
	"expvar"

	"github.com/gin-gonic/gin"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		vars := map[string]string{}
		expvar.Do(func(kv expvar.KeyValue) {
			vars[kv.Key] = kv.Value.String()
		})

		out, err := json.MarshalIndent(vars, "", "\t")
		if err != nil {
			c.Error(err)
		}

		w := c.Writer
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(out)
	}
}
