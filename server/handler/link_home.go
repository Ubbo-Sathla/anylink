package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/Ubbo-Sathla/anylink/admin"
)

func LinkHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LinkHome Get", r.RemoteAddr)
	hu, _ := httputil.DumpRequest(r, true)
	fmt.Println("DumpHome: ", string(hu))

	connection := strings.ToLower(r.Header.Get("Connection"))
	userAgent := strings.ToLower(r.UserAgent())
	if connection == "close" && (strings.Contains(userAgent, "anyconnect") || strings.Contains(userAgent, "openconnect")) {
		w.Header().Set("Connection", "close")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "hello world!")
}

func LinkOtpQr(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	idS := r.FormValue("id")
	jwtToken := r.FormValue("jwt")
	data, err := admin.GetJwtData(jwtToken)
	if err != nil || idS != fmt.Sprint(data["id"]) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	admin.UserOtpQr(w, r)
}
