package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/efrenfuentes/ipinfo/internal/ipinfo_client"
)

func (api *API) GetIpHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
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

	format := get_param(r, "format", "json")

	if format == "json" {
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
	} else if format == "text" {
		var response_body string

		response_body = fmt.Sprintf("Your IP address is: %s\nCity: %s\nRegion: %s\nCountry: %s\nLoc: %s\nOrg: %s\nPostal: %s\nTimezone: %s\n",
			ip, details.City, details.Region, details.Country, details.Loc, details.Org, details.Postal, details.Timezone)

		err = api.WritePlainText(w, http.StatusOK, response_body, nil)
		if err != nil {
			api.ServerErrorResponse(w, r, err)
		}
	} else {
		api.ErrorResponse(w, r, http.StatusBadRequest, "Invalid format")
	}

	return
}

func getIP(r *http.Request) string {
	// We'll always grab the first IP address in the X-Forwarded-For header
	// list.  We do this because this is always the *origin* IP address, which
	// is the *true* IP of the user.  For more information on this, see the
	// Wikipedia page: https://en.wikipedia.org/wiki/X-Forwarded-For
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]).String()

	// If the IP address is empty, then we'll just return an empty string.
	if ip == "" {
		return ""
	}

	return ip
}

func get_param(r *http.Request, param string, default_value ...string) string {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	values, ok := r.Form[param]
	if ok && len(values) > 0 {
		return values[0]
	}

	if len(default_value) > 0 {
		return default_value[0]
	}

	return ""
}