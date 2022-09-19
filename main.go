package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"time"
)

// NewProxy takes target host and creates a reverse proxy
func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url2.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	// Certificate check skip
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	proxy.FlushInterval = 10 * time.Millisecond

	// For modifying the request
	originalDirector := proxy.Director
	proxy.Director = func(request *http.Request) {
		//request.Host = url.Host
		request.URL.Scheme = url.Scheme
		request.URL.Host = url.Host

		originalDirector(request)
		modifyRequest(request)
	}

	//proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

// Modifying the request
func modifyRequest(req *http.Request) {
	remoteAddress, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Add("X-Forwarded-For", remoteAddress)
	req.Header.Set("Accept", "application/json")
}

// Error handler
func errorHandler() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		fmt.Printf("ErrorHanddler: %v \n", err)
		return
	}
}

//Modifying the response
func modifyResponse() func(*http.Response) error {
	return func(response *http.Response) error {
		response.Header.Set("Test", "test")
		return nil
	}
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		proxy.ServeHTTP(w, r)
	}
}

func logRequest(r *http.Request) {
	targetHost := "https://httpbin.org/"
	targetURI := r.RequestURI

	fmt.Printf("Redirecting to: %s%s\n", targetHost, targetURI)
}

func main() {
	// Initialize a reverse proxy and pass the actual backend server url here
	proxy, err := newProxy("https://httpbin.org")
	if err != nil {
		panic(err)
	}

	// Handle all requests to the server using the proxy
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
