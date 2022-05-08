package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
        "flag"
)

type loggingStatusCode struct {
	status int
	http.ResponseWriter
}

func (logWirter *loggingStatusCode) WriteHeader(statusCode int) {
	fmt.Printf("status code => %d \n", statusCode)
	logWirter.WriteHeader(statusCode)
	logWirter.status = statusCode
}

func withAccessLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRespWriter := loggingStatusCode{status: 0, ResponseWriter: w}
		fn.ServeHTTP(&logRespWriter, r)
		fmt.Sprintf("%s, %d", r.RemoteAddr, logRespWriter.status)
		fmt.Printf("RemoteAddr: %s, statusCode: %d\n", r.RemoteAddr, logRespWriter.status)
	}
}

func withHeaderUpdate(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			w.Header().Set(k, v[0])
			fmt.Printf("%s, %s\n", k, v[0])
		}
		w.Header().Add("VERSION", os.Getenv("VERSION"))
		fn(w, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hello, World!")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}

func main() {
	http.HandleFunc("/", withAccessLog(withHeaderUpdate(handler)))
	http.HandleFunc("/healthz", withAccessLog(withHeaderUpdate(healthz)))

	var port string
        port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8099"
	}

	flag.StringVar(&port, "p", port, "Specify http server port. Default is 8099")
        flag.Parse()

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(http.ErrAbortHandler)
	}
}
