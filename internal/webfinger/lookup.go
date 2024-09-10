package webfinger

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func LookupResource(instance string, resource string) (webfinger *Webfinger, err error) {
	webfingerUrl := fmt.Sprintf("https://%s/.well-known/webfinger?resource=%s", instance, resource)
	response, err := http.Get(webfingerUrl)
	if err != nil {
		return
	}
	defer response.Body.Close()
	webfinger = &Webfinger{}
	err = json.NewDecoder(response.Body).Decode(webfinger)
	if err != nil {
		return
	}
	return
}
