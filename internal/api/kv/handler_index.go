package kv

import (
	coreStorage "github.com/balerter/balerter/internal/corestorage"
	"go.uber.org/zap"
	"net/http"
)

// HandlerIndex handle API request GET /api/v1/kv
func HandlerIndex(storage coreStorage.CoreStorage, logger *zap.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		var err error

		data, err := storage.KV().All()
		if err != nil {
			logger.Error("error get kv data", zap.Error(err))
			rw.Header().Add("X-Error", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = newResource(data).render(rw)
		if err != nil {
			logger.Error("error write response", zap.Error(err))
			rw.Header().Add("X-Error", "error write response")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
