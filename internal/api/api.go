package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// IPAddress is a struct we use to represent JSON API responses.
type IPAddress struct {
	IP string `json:"ip"`
}

func GetIpAddress(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// We'll always grab the first IP address in the X-Forwarded-For header
	// list.  We do this because this is always the *origin* IP address, which
	// is the *true* IP of the user.  For more information on this, see the
	// Wikipedia page: https://en.wikipedia.org/wiki/X-Forwarded-For
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]).String()

	// If the user specifies a 'format' querystring, we'll try to return the
	// user's IP address in the specified format.
	if format, ok := r.Form["format"]; ok && len(format) > 0 {
		jsonStr, _ := json.Marshal(IPAddress{ip})

		switch format[0] {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(jsonStr))
			return
		}
	}

	// IP in plain text.
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, ip)
}
