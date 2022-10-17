package api

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/efrenfuentes/ipinfo/internal/ipinfo_client"
)

func (api *API) GetIpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// We'll always grab the first IP address in the X-Forwarded-For header
	// list.  We do this because this is always the *origin* IP address, which
	// is the *true* IP of the user.  For more information on this, see the
	// Wikipedia page: https://en.wikipedia.org/wiki/X-Forwarded-For
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]).String()

	// If the IP address is empty, then we'll just return an empty string.
	if ip == "" {
		api.ServerErrorResponse(w, r, err)
		return
	}

	show_details := get_param(r, "details") == "true"

	var ipinfo_api *ipinfo_client.Client
	var ip_details *ipinfo_client.Details

	if show_details {
		ipinfo_api = ipinfo_client.NewClient()
		ipinfo_api.SetAccessToken(api.Config.IpinfoAccessToken)
		ip_details, err = ipinfo_api.GetDetails(ip)
		if err != nil {
			api.ServerErrorResponse(w, r, err)
			return
		}
	}

	format := get_param(r, "format", "text")

	switch format {
	case "json":
		var env Envelope
	
		if show_details {
			env = Envelope{
				"ip_info": ip_details,
			}
		} else {
			env = Envelope{
				"ip": ip,
			}
		}

		err = api.WriteJSON(w, http.StatusOK, env, nil)
		if err != nil {
			api.ServerErrorResponse(w, r, err)
		}
		return
	default:
		var response_body string

		if show_details {
			response_body = fmt.Sprintf("Your IP address is: %s\nCity: %s\nRegion: %s\nCountry: %s\nLoc: %s\nOrg: %s\nPostal: %s\nTimezone: %s\n",
				ip, ip_details.City, ip_details.Region, ip_details.Country, ip_details.Loc, ip_details.Org, ip_details.Postal, ip_details.Timezone)
		} else {
			response_body = ip
		}

		err = api.WritePlainText(w, http.StatusOK, response_body, nil)
		if err != nil {
			api.ServerErrorResponse(w, r, err)
		}
	}
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