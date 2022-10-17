package api

import (
	"net"
	"net/http"
	"strings"

	"github.com/efrenfuentes/ipinfo/internal/ipinfo_client"
)

func (api *API) GetIpHandler(w http.ResponseWriter, r *http.Request) {
	// We'll always grab the first IP address in the X-Forwarded-For header
	// list.  We do this because this is always the *origin* IP address, which
	// is the *true* IP of the user.  For more information on this, see the
	// Wikipedia page: https://en.wikipedia.org/wiki/X-Forwarded-For
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]).String()

	if ip == "" {
		api.ErrorResponse(w, r, http.StatusInternalServerError, "can't get the IP address")
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
