package main

import (
	"log"
	"testing"

	"github.com/essentier/spickspan"
	"github.com/essentier/spickspan/testutil"
)

var provider = spickspan.GetNomockProvider()

func init() {
	err := spickspan.BuildAll()
	if err != nil {
		log.Printf("Failed to build projects. The error is %v", err)
	}
}

func TestHelloAPI(t *testing.T) {
	t.Parallel()
	helloService, err := testutil.CreateRestService("hello-nomock1", provider)
	defer helloService.Release()

	var helloResult map[string]string
	_, err = helloService.Res("hello", &helloResult).Get()
	if err != nil {
		t.Fatalf("Failed to call the hello rest api. Error is: %v.", err)
	}
	t.Logf("hellResult is %v", helloResult)
	expectedMessage := "Hello, World!"
	if helloResult["message"] != expectedMessage {
		t.Errorf("hello message should be %v but is: %v", expectedMessage, helloResult["message"])
	}
}

// func TestNonexistentAPI(t *testing.T) {
// 	t.Parallel()
// 	helloService := getService("hello-nomock", t)
// 	defer provider.Release(helloService)

// 	hostUrl := helloService.GetHttpUrl()
// 	api := gopencils.Api(hostUrl)
// 	var result interface{}
// 	res := api.Res("nonexistent", &result)
// 	_, err := res.Get()
// 	if err != nil {
// 		t.Fatalf("Failed to call the API. Error is: %v.", err)
// 	}

// 	if res.Raw.StatusCode != http.StatusNotFound {
// 		t.Errorf("Calling a nonexistent API should return status code 404. The status code was %v", res.Raw.StatusCode)
// 	}
// }
