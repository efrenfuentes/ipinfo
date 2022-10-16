package api

import (
	"net/http"
)

func (api *API) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := Envelope{
		"status": "available",
		"system_info": map[string]string{
			"version": version,
		},
	}

	err := api.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		api.ServerErrorResponse(w, r, err)
	}
}
