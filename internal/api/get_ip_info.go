package api

import (
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/efrenfuentes/ipinfo/internal/ipinfo_client"
)

func (api *API) GetIpInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the "ip" parameter from the request context.
	params := httprouter.ParamsFromContext(r.Context())

	ip := params.ByName("ip")

	if net.ParseIP(ip) == nil {
		api.ErrorResponse(w, r, http.StatusBadRequest, "Invalid IP address")
		return
	}

	ipinfo_api := ipinfo_client.NewClient()
	ipinfo_api.SetAccessToken(api.Config.IpinfoAccessToken)
	details, err := ipinfo_api.GetDetails(ip)
	if err != nil {
		api.ServerErrorResponse(w, r, err)
		return
	}

	var env Envelope
	
	env = Envelope{
		"ip": ip,
		"city": details.City,
		"region": details.Region,
		"country": details.Country,
		"loc": details.Loc,
		"org": details.Org,
		"postal": details.Postal,
		"timezone": details.Timezone,
	}

	err = api.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		api.ServerErrorResponse(w, r, err)
	}	

	return
}
