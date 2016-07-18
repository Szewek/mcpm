package mcpmdb

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/Szewek/mcpm/util"
)

const (
	addonsCDN       = "addons.cursecdn.com"
	addonsOriginCDN = "addons-origin.cursecdn.com"
	uagent          = "CurseClient/7.1 (MCPM; +https://github.com/Szewek/mcpm)"
)

var client = &http.Client{
	Timeout: time.Minute,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		req.Header.Set("User-Agent", uagent)
		if req.Host == addonsCDN {
			req.Host = addonsOriginCDN
			req.URL.Host = addonsOriginCDN
		}
		return nil
	},
}

func init() {
	var er error
	client.Jar, er = cookiejar.New(nil)
	util.Must(er)
}

func get(ur string) (*http.Response, error) {
	u, ue := url.Parse(ur)
	util.Must(ue)
	req := &http.Request{
		Method:     "GET",
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Body:       nil,
		Host:       u.Host,
	}
	req.Header.Set("User-Agent", uagent)
	return client.Do(req)
}
