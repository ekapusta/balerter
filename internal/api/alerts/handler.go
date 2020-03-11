package alerts

import (
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"go.uber.org/zap"
	"net/http"
)

type alertManagerAPIer interface {
	GetAlerts() []*alertManager.AlertInfo
}

// Handler handle API request GET /api/v1/alerts
//
// Endpoint receive arguments:
// name=<NAME1>,<NAME2> - filter by name
// level=error,success - filter by alert level
//
// Examples:
// GET /api/v1/alerts?level=error
// GET /api/v1/alerts?level=error,warn&name=foo
// GET /api/v1/alerts?level=error,warn&name=foo,bar
func Handler(alertManager alertManagerAPIer, logger *zap.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		data := alertManager.GetAlerts()

		data, err = filter(req, data)
		if err != nil {
			logger.Error("error filter alerts", zap.Error(err))
			rw.Header().Add("X-Error", err.Error())
			rw.WriteHeader(400)
			return
		}

		err = newResource(data).render(rw)
		if err != nil {
			logger.Error("error write response", zap.Error(err))
			rw.Header().Add("X-Error", "error write response")
			rw.WriteHeader(500)
			return
		}
	}
}
