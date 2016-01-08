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

	res, _ := helloService.Resource("nonexistent").Get(nil)

	if res.Resource.Raw.StatusCode != http.StatusNotFound {
		t.Errorf("Calling a nonexistent API should return status code 404. The status code was %v", res.Resource.Raw.StatusCode)
	}

	// t.Parallel()
	// helloService := getService("hello-nomock", t)
	// defer provider.Release(helloService)

	// hostUrl := helloService.GetHttpUrl()
	// api := gopencils.Api(hostUrl)
	// var result interface{}
	// res := api.Res("nonexistent", &result)
	// _, err := res.Get()
	// if err != nil {
	// 	t.Fatalf("Failed to call the API. Error is: %v.", err)
	// }

	// if res.Raw.StatusCode != http.StatusNotFound {
	// 	t.Errorf("Calling a nonexistent API should return status code 404. The status code was %v", res.Raw.StatusCode)
	// }
}
