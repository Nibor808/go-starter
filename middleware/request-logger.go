package middleware

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func LogRequest(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.SetOutput(os.Stdout)
		uri := r.URL.String()
		host := r.Host

		var ip string
		hdrRealIP := r.Header.Get("X-Real_Ip")
		hdrForwardedFor := r.Header.Get("X-Forwarded-For")

		if hdrRealIP == "" || hdrForwardedFor == "" {
			ip = r.RemoteAddr
		} else {

		}

		log.Println(r.Method, uri, host, ip, r.Header)
		h(w, r, p)
	}
}
