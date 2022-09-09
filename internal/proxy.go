package awsserviceproxy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/elliotchance/sshtunnel"
)

type RunConfiguration struct {
	Tunnel   TunnelConfiguration
	Services []ServiceConfiguration
}

type TunnelConfiguration struct {
	UserAndHost         string
	CertificateLocation string
	Destination         string
}

type ServiceConfiguration struct {
	ServiceName string
	Port        int
}

func Start(env string) {
	var wg sync.WaitGroup
	wg.Add(1)
	tunnel := NewTunnelConfiguration(env)
	port := setupTunnel(tunnel)
	setupGlobalRequestHandler(fmt.Sprintf("http://localhost:%d", port))
	wg.Wait()
}

func NewTunnelConfiguration(env string) TunnelConfiguration {
	return TunnelConfiguration{
		CertificateLocation: FileName(env),
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
	reqBody, respBody := "", ""
	if req.Body != nil {
		reqBody, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}
	if resp.Body != nil {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	sugar, _ := newProductionZaplogger(req.Host)
	message := fmt.Sprintf("%s %s %d %s %s", req.Host, req.Method, resp.StatusCode, req.URL.Path, req.URL.RawQuery)
	sugar.Infow(
		message,
		"service", req.Host,
		"method", req.Method,
		"path", req.URL.Path,
		"query", req.URL.RawQuery,
		"requestBody", string(reqBody),
		"responseBody", string(respBody),
		"status", resp.StatusCode,
		"requestHeaders", req.Header,
		"responseHeaders", resp.Header,
	)
	log.Println(message)

	return resp, nil
}
