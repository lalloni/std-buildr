package httpe

import (
	"io"
	"net/http"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"gitlab.cloudint.afip.gob.ar/std/std-buildr/credentials"
)

func Put(client *http.Client, creds *credentials.Credentials, url string, body io.Reader, expect int) (int, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return 0, errors.Wrap(err, "creating http request")
	}
	if creds != nil {
		req.SetBasicAuth(creds.Username, creds.Password)
	}
	log.Debugf("sending http request: %+v", req)
	res, err := client.Do(req)
	if err != nil {
		log.Debugf("got error from http client: %+v", err)
		return 0, errors.Wrap(err, "executing http request")
	}
	log.Debugf("received http response: %+v", res)
	defer res.Body.Close()
	if expect != 0 && expect != res.StatusCode {
		return 0, errors.Errorf("unexpected http response status '%s' (expecting '%d')", res.Status, expect)
	}
	log.Debugf("http request was succesful")
	return res.StatusCode, nil
}
