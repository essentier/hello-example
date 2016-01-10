package main

import (
	"testing"

	"github.com/essentier/testutil"
)

func TestHelloAPI(t *testing.T) {
	t.Parallel()
	helloService := testutil.CreateRestService("hello-example", t)
	defer helloService.Release()

	var helloResult map[string]string
	helloService.Resource("hello").Get(&helloResult)

	t.Logf("helloResult is %v", helloResult)
	expectedMessage := "Hello, World!!"
	if helloResult["message"] != expectedMessage {
		t.Errorf("hello message should be %v but is: %v", expectedMessage, helloResult["message"])
	}
}
