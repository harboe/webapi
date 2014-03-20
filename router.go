package webapi

import (
	"log"
	"fmt"
	"regexp"
)

type WebApiRouter struct {
	basePath 		string
	Routes 			[]*WebApiRoute
}

func (router *WebApiRouter) Map(method, path, accept string, handler WebApiHandler) *WebApiRoute {

	path = router.basePath + path;

	r := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern := r.ReplaceAllStringFunc(path, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	});
	
	var route = &WebApiRoute{method, accept, path, handler, regexp.MustCompile(pattern)};
	router.Routes = append(router.Routes, route);
	
	log.Println("Route:", route.Method, "=>",route.Path);
	return route;
}

func (router *WebApiRouter) Match(method, path string) (h WebApiHandler, p map[string]string) {

	for _, route := range router.Routes {
		if match, params := route.Match(method, path); match == true {
			return route.handler, params;
		}
	}

	return nil, nil;
}

type WebApiRoute struct {
	Method			string
	Accept 			string
	Path			string
	handler     	WebApiHandler
	regex 			*regexp.Regexp
}

func (r *WebApiRoute) MatchMethod(method string) bool {
	return r.Method == "*" || method == r.Method || (method == "HEAD" && r.Method == "GET");
}

func (r *WebApiRoute) Match(method, path string) (bool, map[string]string) {
	
	// add Any method matching support
	if !r.MatchMethod(method) {
		return false, nil
	}

	matches := r.regex.FindStringSubmatch(path);
	if len(matches) > 0 && matches[0] == path {
		params := make(map[string]string)
		for i, name := range r.regex.SubexpNames() {
			if len(name) > 0 {
				params[name] = matches[i]
			}
		}
		return true, params
	}
	return false, nil
}

func newRouter(basePath string) *WebApiRouter {
	r := new(WebApiRouter);
	r.basePath = basePath;

	return r;
}