package handler

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/Ubbo-Sathla/anylink/base"
	"github.com/Ubbo-Sathla/anylink/pkg/proxyproto"
	"github.com/gorilla/mux"
)

func startTls() {

	var (
		err error

		addr     = base.Cfg.ServerAddr
		certFile = base.Cfg.CertFile
		keyFile  = base.Cfg.CertKey
		certs    = make([]tls.Certificate, 1)
		ln       net.Listener
	)

	// 判断证书文件
	// _, err = os.Stat(certFile)
	// if errors.Is(err, os.ErrNotExist) {
	//	// 自动生成证书
	//	certs[0], err = selfsign.GenerateSelfSignedWithDNS("vpn.anylink")
	// } else {
	//	// 使用自定义证书
	//	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	// }
	caBytes, err := ioutil.ReadFile("/app/conf/ca.crt")
	if err != nil {
		panic("Unable to read /app/conf/ca.crt")
	}
	caPool := x509.NewCertPool()
	ok := caPool.AppendCertsFromPEM(caBytes)
	if !ok {
		panic("failed to parse root certificate")
	}

	certs[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	// 设置tls信息
	tlsConfig := &tls.Config{
		NextProtos:         []string{"http/1.1"},
		MinVersion:         tls.VersionTLS12,
		Certificates:       certs,
		ClientAuth:         tls.NoClientCert,
		RootCAs:            caPool,
		InsecureSkipVerify: true,
	}
	fmt.Printf("%#v\n", tlsConfig)
	srv := &http.Server{
		Addr:      addr,
		Handler:   initRoute(),
		TLSConfig: tlsConfig,
		ErrorLog:  base.GetBaseLog(),
	}

	ln, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	if base.Cfg.ProxyProtocol {
		ln = &proxyproto.Listener{Listener: ln, ProxyHeaderTimeout: time.Second * 5}
	}

	base.Info("listen server", addr)
	err = srv.ServeTLS(ln, "", "")
	if err != nil {
		base.Fatal(err)
	}
}

func initRoute() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", LinkHome).Methods(http.MethodGet)
	r.HandleFunc("/", LinkAuth).Methods(http.MethodPost)
	r.HandleFunc("/CSCOSSLC/tunnel", LinkTunnel).Methods(http.MethodConnect)
	r.HandleFunc("/otp_qr", LinkOtpQr).Methods(http.MethodGet)
	r.HandleFunc("/profile.xml", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("profile.xml Get", r.RemoteAddr)
		hu, _ := httputil.DumpRequest(r, true)
		fmt.Println("DumpHome: ", string(hu))
		b, _ := os.ReadFile(base.Cfg.Profile)
		w.Write(b)
	}).Methods(http.MethodGet)
	r.PathPrefix("/files/").Handler(
		http.StripPrefix("/files/",
			http.FileServer(http.Dir(base.Cfg.FilesPath)),
		),
	)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	return r
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RemoteAddr)
	hu, _ := httputil.DumpRequest(r, true)
	fmt.Println("NotFound: ", string(hu))

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 page not found")
}
