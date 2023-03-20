package main

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

type RequestTime struct {
	id               int
	url              string
	ip               string
	status           string
	addrs            []net.IPAddr
	dnsLookup        time.Duration
	connectTime      time.Duration
	tlsHandshake     time.Duration
	serverProcessing time.Duration
	contentTransfer  time.Duration
	total            time.Duration
}

func NewRequestTime() RequestTime {
	return RequestTime{
		dnsLookup:        time.Duration(0),
		connectTime:      time.Duration(0),
		tlsHandshake:     time.Duration(0),
		serverProcessing: time.Duration(0),
		contentTransfer:  time.Duration(0),
		total:            time.Duration(0),
	}
}

func isValidUrl(urlToTest string) (string, error) {

	// add 'http' prefix
	if !strings.HasPrefix(urlToTest, "http") {
		urlToTest = "https://" + urlToTest
	}

	_, err := url.ParseRequestURI(urlToTest)
	if err != nil {
		return urlToTest, errors.New("invalid url")
	}

	u, err := url.Parse(urlToTest)
	_ = err
	if err != nil || u.Scheme == "" || u.Host == "" {
		return urlToTest, errors.New("invalid url")
	}

	return urlToTest, nil
}

var counter = 0

func ExecuteRequest(siteUrl string) (RequestTime, error) {

	siteUrl, err := isValidUrl(siteUrl)
	if err != nil {
		return NewRequestTime(), err
	}

	req, err := http.NewRequest("GET", siteUrl, nil)
	if err != nil {
		return NewRequestTime(), err
	}

	var start, end, connectStart, connectDone, dnsStart, dnsDone, tlsHandshakeStart, tlsHandshakeDone, gotConn, firstResponseByte time.Time
	requestTime := NewRequestTime()
	requestTime.url = siteUrl

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsDone = time.Now()
			requestTime.addrs = ddi.Addrs
		},

		ConnectStart: func(network, addr string) { connectStart = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connectDone = time.Now()
			requestTime.ip = addr
		},

		TLSHandshakeStart: func() { tlsHandshakeStart = time.Now() },
		TLSHandshakeDone:  func(cs tls.ConnectionState, err error) { tlsHandshakeDone = time.Now() },

		GotConn:              func(info httptrace.GotConnInfo) { gotConn = time.Now() },
		GotFirstResponseByte: func() { firstResponseByte = time.Now() },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	transport := http.Transport{
		DisableKeepAlives: true,

		MaxIdleConns:          1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Transport: &transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// always refuse to follow redirects, visit does that
			// manually if required.
			return http.ErrUseLastResponse
		},
	}
	start = time.Now()
	resp, err := client.Do(req)
	//resp, err := transport.RoundTrip(req)
	end = time.Now()
	if err != nil {
		return NewRequestTime(), err
	}

	requestTime.status = resp.Status

	counter++
	requestTime.id = counter

	requestTime.dnsLookup = dnsDone.Sub(dnsStart)
	requestTime.connectTime = connectDone.Sub(connectStart)
	requestTime.tlsHandshake = tlsHandshakeDone.Sub(tlsHandshakeStart)
	requestTime.serverProcessing = firstResponseByte.Sub(gotConn)
	requestTime.contentTransfer = end.Sub(firstResponseByte)
	requestTime.total = end.Sub(start)

	return requestTime, nil
}
