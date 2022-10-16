package api

import (
	"net"
	"net/http"
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
	// ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]).String()

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		api.ServerErrorResponse(w, r, err)
		return
	} 

	// If the user specifies a 'format' querystring, we'll try to return the
	// user's IP address in the specified format.
	if format, ok := r.Form["format"]; ok && len(format) > 0 {
		env := Envelope{"ip": ip}

		switch format[0] {
		case "json":
			err = api.WriteJSON(w, http.StatusOK, env, nil)
			if err != nil {
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
	}

	err = api.WritePlainText(w, http.StatusOK, ip, nil)
	if err != nil {
		api.ServerErrorResponse(w, r, err)
	}
}
