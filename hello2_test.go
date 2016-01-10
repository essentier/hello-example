package main

import (
	"net/http"
	"testing"

	"github.com/essentier/testutil"
)

func TestNonexistentAPI(t *testing.T) {
	t.Parallel()
	helloService := testutil.CreateRestService("hello-example", t)
	defer helloService.Release()

	res := helloService.Resource("nonexistent").Get(nil)

	if res.Resource.Raw.StatusCode != http.StatusNotFound {
		t.Errorf("Calling a nonexistent API should return status code 404. The status code was %v", res.Resource.Raw.StatusCode)
	}
}
