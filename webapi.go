package webapi

import (
	"log"
	"net/http"
	"strings"
	"regexp"
)

const (
	GET 	= "GET"
	POST 	= "POST"
	PUT 	= "PUT"
	DELETE 	= "DELETE"
)

// The regex to check for the requested format (allows an optional trailing
// slash).
var rxExt = regexp.MustCompile(`(\.(?:xml|txt|json))\/?$`);

type WebApi struct{
	Router 			*WebApiRouter
	handlers		[]http.HandlerFunc
	encoders 		map[string]WebApiEncoder
}

type WebApiHandler func(r *WebApiRequest) (int, interface{});

func (api *WebApi) Get(path string, handler WebApiHandler) {
	api.Map(GET, path, "*", handler);
}

func (api *WebApi) Post(path, accept string, handler WebApiHandler) {
	api.Map(POST, path, accept, handler);
}

func (api *WebApi) Put(path, accept string, handler WebApiHandler) {
	api.Map(PUT, path, accept, handler);
}

func (api *WebApi) Delete(path string, handler WebApiHandler) {
	api.Map(DELETE, path, "*", handler);
}

func (api *WebApi) Map(method, path, accept string, handler WebApiHandler) {
	api.Router.Map(method, path, accept, handler);
}      

func (api *WebApi) Use(handler http.HandlerFunc) {
	api.handlers = append(api.handlers, handler);
}

func (api *WebApi) Encoder(contentType string, handler WebApiEncoder) {
	api.encoders[contentType] = handler;
}

func New(basePath string) *WebApi {
	log.Println("Initializing new webapi controller...");

	if len(basePath) == 0 {
		basePath = "/api/";
	}

	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath;
	}

	if !strings.HasSuffix(basePath, "/") {
		basePath = basePath + "/"
	}

	var api = new(WebApi);
	api.encoders = make(map[string]WebApiEncoder);
	api.Router = newRouter(basePath);

	api.Encoder("application/json", JsonEncoder);
	api.Encoder("application/xml", XmlEncoder);
	api.Encoder("text/plain", TextEncoder);
	api.Use(api.defaultContentNegotiation);

	http.HandleFunc(basePath, api.requestHandler());
	return api;
}

func (api *WebApi) defaultContentNegotiation(w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ""
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}

	switch ft {
		case ".xml":
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		case ".txt":
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		default:
			var acceptHeader = r.Header["Accept"][0];
	
			for _, x := range strings.Split(acceptHeader, ",") {
				// ignore encoding... :(
				if index := strings.IndexRune(x, ';'); index > 0 {
					x = x[0:index];
				}

				for key, _ := range api.encoders {
					if x == key {
						w.Header().Set("Content-Type", key);
						return;
					}
				}			
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8");
	}
}

func (api *WebApi) requestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// invoke custom handlers
		for _, handler := range api.handlers {
			handler(w, r);
		}

		if handler, params := api.Router.Match(r.Method, r.URL.Path); handler != nil {
			// build webapi request.
			webRequest := &WebApiRequest{r, params};
			// handling the request.
			statusCode, result := handler(webRequest);

			// encoding the response.
			contentType := w.Header().Get("Content-Type");
			body, _ := api.encode(contentType, result);

			w.WriteHeader(statusCode);
			w.Write(body);

		} else {
			http.NotFound(w, r);
		}
	}
}

func (api *WebApi) encode(ct string, v interface{}) ([]byte, error){

	if i := strings.IndexRune(ct, ';'); i > 0 {
		ct = ct[0:i];
	}

	if encoder := api.encoders[ct]; encoder != nil {
		return encoder(v);
	}

	return nil, nil;
}
