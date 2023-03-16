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
}

func NewRequestTime() RequestTime {
	return RequestTime{
		dnsLookup:        time.Duration(0),
		connectTime:      time.Duration(0),
		tlsHandshake:     time.Duration(0),
		serverProcessing: time.Duration(0),
		contentTransfer:  time.Duration(0),
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

func ExecuteRequest(siteUrl string) (RequestTime, error) {

	siteUrl, err := isValidUrl(siteUrl)
	if err != nil {
		return NewRequestTime(), err
	}

	req, err := http.NewRequest("GET", siteUrl, nil)
	if err != nil {
		return NewRequestTime(), err
	}

	var start, connect, dns, tlsHandshake time.Time
	requestTime := NewRequestTime()

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			requestTime.dnsLookup = time.Since(dns)
			//log.Printf("DNS Done: %v\n", time.Since(dns))
			requestTime.addrs = ddi.Addrs
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			requestTime.tlsHandshake = time.Since(tlsHandshake)
			//log.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
			//log.Printf(cs.ServerName)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			requestTime.connectTime = time.Since(connect)
			//log.Printf("Connect time: %v\n", time.Since(connect))
			requestTime.ip = addr
		},

		GotFirstResponseByte: func() {
			requestTime.serverProcessing = time.Since(start)
			//log.Printf("Time from start to first byte: %v\n", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	transport := http.Transport{
		DisableKeepAlives: true,

		MaxIdleConns:          1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return NewRequestTime(), err
	}

	requestTime.contentTransfer = time.Since(start)
	requestTime.status = resp.Status

	//log.Printf("Total time: %v\n", time.Since(start))
	return requestTime, nil
}
