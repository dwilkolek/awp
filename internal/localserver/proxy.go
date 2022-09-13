package localserver

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/elliotchance/sshtunnel"
	awswebproxy "github.com/tfmcdigital/aws-web-proxy/internal"
)

var logChan chan LogEntry = make(chan LogEntry)

func Start(env string) {
	var wg sync.WaitGroup
	wg.Add(1)

	tunnel := newTunnelConfiguration(env)
	port := setupTunnel(tunnel)
	setupGlobalRequestHandler(fmt.Sprintf("http://localhost:%d", port))

	startLocalWebServer()
	wg.Wait()

}

func newTunnelConfiguration(env string) TunnelConfiguration {
	return TunnelConfiguration{
		CertificateLocation: awswebproxy.FilePathToBastionKey(env),
		UserAndHost:         fmt.Sprintf("ec2-user@bastion.%sservices.technipfmc.com", strings.ReplaceAll(env+".", "prod.", "")),
		Destination:         "service.service:80",
	}

}

func setupTunnel(tunnelConfig TunnelConfiguration) int {
	tunnel := sshtunnel.NewSSHTunnel(
		tunnelConfig.UserAndHost,
		sshtunnel.PrivateKeyFile(tunnelConfig.CertificateLocation),
		tunnelConfig.Destination,
		"0",
	)

	tunnel.Log = log.Default()

	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	tunnel.Log.Printf("Started and exposed on port: %d\n", tunnel.Local.Port)

	return tunnel.Local.Port
}
func setupGlobalRequestHandler(to string) {
	go func() {

		origin, _ := url.Parse(to)
		director := func(req *http.Request) {
			host := req.Host
			req.Header.Add("host", host)
			req.Host = host
			req.URL.Scheme = "http"
			if host == "awp" {
				originAwp, _ := url.Parse("http://localhost:2137")
				req.URL.Host = originAwp.Host
			} else {
				req.URL.Host = origin.Host
			}

		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.Transport = &transport{http.DefaultTransport}
		server := http.NewServeMux()

		server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})

		// logger.Printf("Starting server at port: %d\n", 80)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", 80), server); err != nil {
			log.Fatal("Failed. Try to execute `lsof -t -i tcp:80 | xargs kill`.")
		}

		log.Printf("Started: %d\n", 80)
	}()

}

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, _ = t.RoundTripper.RoundTrip(req)

	reqBody, _ := httputil.DumpRequest(req, true)
	respBody, _ := httputil.DumpResponse(resp, true)

	sugar, _ := newProductionZaplogger(req.Host)

	logEntry := LogEntry{
		Message:         fmt.Sprintf("%s %s %d %s %s", req.Host, req.Method, resp.StatusCode, req.URL.Path, req.URL.RawQuery),
		Service:         req.Host,
		Method:          req.Method,
		Path:            req.URL.Path,
		Query:           req.URL.RawQuery,
		Request:         string(reqBody),
		Response:        string(respBody),
		Status:          resp.StatusCode,
		RequestHeaders:  req.Header,
		ResponseHeaders: resp.Header,
	}
	sugar.Infow(
		logEntry.Message,
		"service", logEntry.Service,
		"method", logEntry.Method,
		"path", logEntry.Path,
		"query", logEntry.Query,
		"request", logEntry.Request,
		"response", logEntry.Response,
		"status", logEntry.Status,
		"requestHeaders", logEntry.RequestHeaders,
		"responseHeaders", logEntry.ResponseHeaders,
	)
	log.Println(logEntry.Message)
	logChan <- logEntry
	return resp, nil
}

type LogEntry struct {
	Message         string              `json:"message"`
	Service         string              `json:"service"`
	Method          string              `json:"method"`
	Path            string              `json:"path"`
	Query           string              `json:"query"`
	Request         string              `json:"request"`
	Response        string              `json:"response"`
	Status          int                 `json:"status"`
	RequestHeaders  map[string][]string `json:"requestHeaders"`
	ResponseHeaders map[string][]string `json:"responseHeaders"`
}
