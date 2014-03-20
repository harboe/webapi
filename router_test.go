package webapi

import (
	"log"
	"testing"
)

var dummyHandler = func(r *WebApiRequest) (int, interface{}) { 
	log.Println("dummy..."); 
	return 500, nil;
};

func Test_simple_route(t *testing.T) {
	router := newRouter("api");
	r := router.Map(GET, "test/:id", "*", dummyHandler);

	if _, params := router.Match(GET, "/api/test/1"); params != nil {
		log.Println(params);
	} else {
		log.Println("bummer...")
	}

	log.Println(r)
	log.Println("testing....")
}