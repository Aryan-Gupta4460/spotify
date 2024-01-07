package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"go.elastic.co/apm/module/apmgorilla/v2"
	"go.elastic.co/apm/module/apmhttp/v2"
	"go.elastic.co/apm/v2"
	"go.uber.org/zap"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HTTPReqInfo struct {
	method   string
	uri      string
	code     int
	size     int64
	duration time.Duration
}

type HTTPResponse struct {
	Err HTTPError `json:"error"`
}

type HTTPRoute struct {
	Method string
	Path   string
}

var LOGGER *zap.SugaredLogger

type HTTPController interface {
	SetRouteHandlers(r *mux.Router)
}

type WebServer struct {
	host string
	port string

	routes     *mux.Router
	srv        *http.Server
	rootRouter *mux.Router
	log        *zap.SugaredLogger
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				LOGGER.Errorf("Recovering from panic -> %v", r)
				LOGGER.Errorf("Stacktrace from panic -> \n%s", string(debug.Stack()))
			}
		}()
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		traceContext, _ := apmhttp.ParseTraceparentHeader(r.Header.Get("Traceparent"))
		traceContext.State, _ = apmhttp.ParseTracestateHeader(r.Header["Tracestate"]...)

		opts := apm.TransactionOptions{
			TraceContext: traceContext,
		}
		transaction := apm.DefaultTracer().StartTransactionOptions(r.URL.Path, r.Method, opts)
		defer transaction.End()

		defer r.Body.Close()
		next.ServeHTTP(w, r)
	})
}
func NewWebServer(logger *zap.SugaredLogger, host string, port string) *WebServer {
	LOGGER = logger

	rootRouter := mux.NewRouter()
	rootRouter.Use(apmgorilla.Middleware())
	rootRouter.Use(loggingMiddleware)
	prefix := "/api"
	ws := &WebServer{
		host: host, port: port, log: logger,
		routes: rootRouter.PathPrefix(prefix).Subrouter(),
	}
	ws.routes.Use(reqLogger)
	ws.rootRouter = rootRouter
	return ws
}

func (t *WebServer) Start() error {
	addr := fmt.Sprintf("%v:%v", t.host, t.port)
	t.log.Infof("Starting web server @ -> Host: %v, Port: %v", t.host, t.port)
	t.srv = &http.Server{
		Handler:      t.rootRouter,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  45 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return t.srv.ListenAndServe()
}

func (t *WebServer) Stop() error {
	return t.srv.Close()
}

func (t *WebServer) SetRoute(ctrl HTTPController) {
	ctrl.SetRouteHandlers(t.routes)
}

func reqLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method: r.Method,
			uri:    r.URL.String(),
		}
		m := httpsnoop.CaptureMetrics(h, w, r)
		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration

		durInMs := float64(ri.duration) / 1000000
		if durInMs > 200 {
			LOGGER.Info(ri.method, " ", ri.code, " ", ri.uri, " ", ri.duration, " ", ri.size, " Bytes  ***")
		} else {
			LOGGER.Info(ri.method, " ", ri.code, " ", ri.uri, " ", ri.duration, " ", ri.size, " Bytes")
		}
	}
	return http.HandlerFunc(fn)
}
