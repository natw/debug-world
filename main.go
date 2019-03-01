package main

import (
	"fmt"
	"html"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "<html lang=\"en\">")

	// adding debug header to test (strong/weak) ETags in combination with NGINX
	w.Header().Set("ETag", "HelloWorld")

	fmt.Fprintln(w, "debug env vars:<br>")
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		key := pair[0]
		if strings.HasPrefix(key, "DEBUG_") {
			val := pair[1]
			fmt.Fprintf(w, "<b>%s</b>: %s<br>\n", key, val)
		}
	}

	var requestKeys []string
	for k := range r.Header {
		requestKeys = append(requestKeys, k)
	}
	sort.Strings(requestKeys)

	var responseKeys []string
	for k := range w.Header() {
		responseKeys = append(responseKeys, k)
	}
	sort.Strings(responseKeys)

	hosts := r.URL.Query()["host"]
	for _, host := range hosts {
		ips, time, _ := lookup(host)
		fmt.Fprintf(w, "lookup <b>%s</b>: %+v (%s)<br>\n", host, ips, time)
	}

	fmt.Fprintln(w, "<b>request.RequestURI:</b>", html.EscapeString(r.RequestURI), "<br>")
	fmt.Fprintln(w, "<b>request.RemoteAddr:</b>", html.EscapeString(r.RemoteAddr), "<br>")
	fmt.Fprintln(w, "<b>request.TLS:</b>", r.TLS, "<br>")
	fmt.Fprintln(w, "<b>Request Headers:</b><br>")
	for _, k := range requestKeys {
		fmt.Fprintln(w, k, ":", r.Header[k], "<br>")
	}

}

func lookup(host string) ([]string, string, error) {
	start := time.Now()
	ips, err := net.LookupHost(host)
	if err != nil {
		return []string{}, "0s", err
	}
	elapsed := time.Since(start)
	return ips, elapsed.String(), nil
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
