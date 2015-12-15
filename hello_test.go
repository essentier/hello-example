package main

import (
	"testing"

	"github.com/essentier/spickspan"
	"github.com/essentier/spickspan/testutil"
)

var provider = spickspan.GetNomockProvider()

func init() {
	//model.LoginToEssentier("http://104.196.21.169:8083", "cha_urwu@hotmail.com", "aaa")
	spickspan.BuildAll()
}

// func TestNonexistentAPI(t *testing.T) {
// 	t.Parallel()
// 	_, r := testutil.SendRestGetToService(t, "hello-nomock", provider, "hello")
// 	if r.StatusCode != http.StatusNotFound {
// 		t.Errorf("Calling a nonexistent API should return status code 404. The status code was %v", r.StatusCode)
// 	}
// }

func TestHelloAPI(t *testing.T) {
	t.Parallel()
	result, _ := testutil.SendRestGetToService(t, "hello-nomock", provider, "hello")
	expectedMessage := "Hello, World!"
	helloResponse := result.(map[string]string)
	if helloResponse["message"] != expectedMessage {
		t.Errorf("hello message should be %v but is: %v", expectedMessage, helloResponse["message"])
	}
}
