package awsserviceproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
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
	setupGlobalRequestHandler(ServiceConfiguration{
		ServiceName: "a",
		Port:        1,
	}, fmt.Sprintf("http://localhost:%d", port))
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

func setupGlobalRequestHandler(conf ServiceConfiguration, to string) {
	go func() {
		logger := log.Default()
		origin, _ := url.Parse(to)

		director := func(req *http.Request) {
			host := strings.Replace(req.Host, "tfmc", "service", 1)
			req.Header.Add("host", host)
			req.Host = host
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}

		proxy := &httputil.ReverseProxy{Director: director}
		server := http.NewServeMux()

		server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			logger.Printf("%s <- %s %s %s\n", conf.ServiceName, r.Method, r.URL.Path, r.URL.RawQuery)
			defer r.Body.Close()
			fmt.Println("ReqUri" + r.RequestURI)
			if r.Host == "app.service" {
				root := "static"
				switch r.Method {
				case "GET":
					if r.URL.Path == "" || r.URL.Path == "/" {
						http.ServeFile(w, r, path.Join(root, "index.html"))
					} else {
						http.ServeFile(w, r, path.Join(root, r.URL.Path))
					}
				}
			} else {
				proxy.ServeHTTP(w, r)
			}

		})

		logger.Printf("Starting server at port: %d\n", 80)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", 80), server); err != nil {
			log.Fatal("Failed. Try to execute `lsof -t -i tcp:80 | xargs kill`.")
		}

		logger.Printf("Started: %d\n", 80)
	}()
}
