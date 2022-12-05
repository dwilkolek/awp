package proxy

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
	"github.com/tfmcdigital/aws-web-proxy/internal/tools/aws"
)

const WEB_SERVER_PORT = 2137

func StartProxy(env domain.Environment) {
	var wg sync.WaitGroup
	wg.Add(1)
	go aws.GetAwsClient().StartBastionProxy(env)
	go globalRequestHandler(env)
	go localWebServer()

	wg.Wait()

}

func globalRequestHandler(env domain.Environment) {
	origin, _ := url.Parse(fmt.Sprintf("http://localhost:%d", domain.SSM_PROXY_PORT))
	director := func(req *http.Request) {
		host := req.Host
		req.Header.Add("host", host)
		req.Host = host
		req.URL.Scheme = "http"

		if host == "awp" {
			originAwp, _ := url.Parse(fmt.Sprintf("http://localhost:%d", WEB_SERVER_PORT))
			req.URL.Host = originAwp.Host
		} else {
			req.URL.Host = origin.Host
			if env != domain.PROD {
				if domain.GetConfig().HeaderOverwrites[host] != nil {
					for header, value := range domain.GetConfig().HeaderOverwrites[host] {
						if strings.HasPrefix(value, "toBase64:") {
							req.Header.Add(header, base64.StdEncoding.EncodeToString([]byte(strings.Replace(value, "toBase64:", "", 1))))
						} else {
							req.Header.Add(header, value)
						}
					}
				}
			}
		}

	}

	proxy := &httputil.ReverseProxy{Director: director}

	proxy.ModifyResponse = func(r *http.Response) error {
		logRoundtrip(r.Request, r)

		return nil
	}
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", 80), server); err != nil {
		log.Fatal("Failed. Try to execute `lsof -t -i tcp:80 | xargs kill`.")
	}

	log.Printf("Started: %d\n", 80)
}

func logRoundtrip(req *http.Request, resp *http.Response) {
	if req == nil {
		return
	}

	if resp == nil {
		return
	}

	if req.Host == "" {
		return
	}

	shouldReadReqBody := strings.Contains(req.Header.Clone().Get("Content-Type"), "application/json")
	shouldReadRespBody := strings.Contains(resp.Header.Clone().Get("Content-Type"), "application/json")
	isRespGzip := strings.Contains(resp.Header.Clone().Get("Content-Encoding"), "gzip")

	reqBody, err := httputil.DumpRequest(req, shouldReadReqBody)
	if err != nil {
		log.Println("Failed to dump request "+req.RequestURI, err)
	}

	respBody, err := httputil.DumpResponse(resp, shouldReadRespBody && !isRespGzip)
	if err != nil {
		log.Println("Failed to dump response", err)
	}

	if shouldReadRespBody && isRespGzip {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close() //  must close
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// fmt.Printf("----------------\n%s\n----------------\n", string(bodyBytes))
		reader := bytes.NewReader([]byte(bodyBytes))
		gzreader, err := gzip.NewReader(reader)
		if err != nil {
			log.Println("Failed to dump response", err)
		} else {
			output, err := ioutil.ReadAll(gzreader)
			if err != nil {
				log.Println("Failed to dump response", err)
			} else {
				respBody = append(respBody, output...)
			}
		}
	}

	logEntry := domain.LogEntry{
		Timestamp:       time.Now().UnixMilli(),
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

	GetLogEntryHandler(req.Host).Submit(&logEntry)

}
