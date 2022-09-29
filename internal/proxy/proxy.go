package proxy

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
	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
)

func StartProxy(env domain.Environment, certLocation string) {
	var wg sync.WaitGroup
	wg.Add(1)
	tunnel := newTunnelConfiguration(env, certLocation)
	port := setupTunnel(tunnel)

	setupGlobalRequestHandler(fmt.Sprintf("http://localhost:%d", port))
	startLocalWebServer()
	wg.Wait()

}

func newTunnelConfiguration(env domain.Environment, certLocation string) TunnelConfiguration {
	return TunnelConfiguration{
		CertificateLocation: certLocation, //awswebproxy.FilePathToBastionKey(env),
		UserAndHost:         fmt.Sprintf("ec2-user@bastion.%sservices.technipfmc.com", strings.ReplaceAll(env.String()+".", "prod.", "")),
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

func logRoundtrip(req *http.Request, resp *http.Response) {
	if req == nil {
		return
	}

	if resp == nil {
		return
	}

	shouldReadReqBody := strings.Contains(req.Header.Clone().Get("Content-Type"), "application/json")
	shouldReadRespBody := strings.Contains(resp.Header.Clone().Get("Content-Type"), "application/json")

	reqBody, err := httputil.DumpRequest(req, shouldReadReqBody)
	if err != nil {
		log.Println("Failed to dump request", err)
	}

	respBody, err := httputil.DumpResponse(resp, shouldReadRespBody)
	if err != nil {
		log.Println("Failed to dump response", err)
	}

	logEntry := domain.LogEntry{
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

	GetLogEntryHandler(req.Host).Submit(logEntry)
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	go logRoundtrip(req, resp)
	return resp, err
}
