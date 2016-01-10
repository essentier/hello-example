package testutil

import (
	"io"
	"net/http"

	"github.com/essentier/gopencils"
)

type resWrapper struct {
	Resource   *gopencils.Resource
	errHandler RestErrorHandler
}

// Creates a new Resource.
func (rw *resWrapper) NewChildResource(resourceName string) *resWrapper {
	newRes := rw.Resource.NewChildResource(resourceName, nil)
	newRW := &resWrapper{errHandler: rw.errHandler, Resource: newRes}
	return newRW
}

// Same as Res() Method, but returns a Resource with url resource/:id
func (rw *resWrapper) NewChildIdResource(id string) *resWrapper {
	newRes := rw.Resource.NewChildIdResource(id)
	newRW := &resWrapper{errHandler: rw.errHandler, Resource: newRes}
	return newRW
}

// Sets QueryValues for current Resource
func (rw *resWrapper) SetQuery(querystring map[string]string) *resWrapper {
	rw.Resource.SetQuery(querystring)
	return rw
}

// Performs a GET request on given Resource
// Call SetQuery beforehand if you want to set the query string of the GET request.
func (rw *resWrapper) Get(responseBody interface{}) *resWrapper {
	rw.Resource.Response = responseBody
	_, err := rw.Resource.Get()
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST GET failed.")
	}
	return rw
}

// Performs a HEAD request on given Resource
// Call SetQuery beforehand if you want to set the query string of the HEAD request.
func (rw *resWrapper) Head() *resWrapper {
	_, err := rw.Resource.Head()
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST HEAD failed.")
	}
	return rw
}

// Performs a PUT request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (rw *resWrapper) Put(payload interface{}, responseBody interface{}) *resWrapper {
	rw.Resource.Response = responseBody
	_, err := rw.Resource.Put(payload)
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST PUT failed.")
	}
	return rw
}

// Performs a POST request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (rw *resWrapper) Post(payload interface{}, responseBody interface{}) *resWrapper {
	rw.Resource.Response = responseBody
	_, err := rw.Resource.Post(payload)
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST POST failed.")
	}
	return rw
}

// Performs a Delete request on given Resource.
// Call SetQuery beforehand if you want to set the query string of the DELETE request.
func (rw *resWrapper) Delete(responseBody interface{}) *resWrapper {
	rw.Resource.Response = responseBody
	_, err := rw.Resource.Delete()
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST DELETE failed.")
	}
	return rw
}

// Performs a OPTIONS request on given Resource.
// Call SetQuery beforehand if you want to set the query string of the OPTIONS request.
func (rw *resWrapper) Options(responseBody interface{}) *resWrapper {
	rw.Resource.Response = responseBody
	_, err := rw.Resource.Options()
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST OPTIONS failed.")
	}
	return rw
}

// Performs a PATCH request on given Resource.
// Accepts interface{} as parameter, will be used as payload.
func (rw *resWrapper) Patch(payload interface{}, responseBody interface{}) *resWrapper {
	rw.Resource.Response = responseBody
	_, err := rw.Resource.Patch(payload)
	if rw.errHandler != nil {
		rw.errHandler.HandleError(err, "REST PATCH failed.")
	}
	return rw
}

// Sets Payload for current Resource
func (rw *resWrapper) SetPayload(args interface{}) io.Reader {
	return rw.Resource.SetPayload(args)
}

// Sets Headers
func (rw *resWrapper) SetHeader(key string, value string) {
	rw.Resource.SetHeader(key, value)
}

// Overwrites the client that will be used for requests.
// For example if you want to use your own client with OAuth2
func (rw *resWrapper) SetClient(c *http.Client) {
	rw.Resource.SetClient(c)
}
