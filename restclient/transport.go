/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package restclient

import (
	"crypto/tls"
	"net/http"

	"github.com/commcos/utils/restclient/transport"
)

// TLSConfigFor returns a tls.Config that will provide the transport level security defined
// by the provided Config. Will return nil if no transport level security is requested.
func TLSConfigFor(config *Config) (*tls.Config, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.TLSConfigFor(cfg)
}

// TransportFor returns an http.RoundTripper that will provide the authentication
// or transport level security defined by the provided Config. Will return the
// default http.DefaultTransport if no special case behavior is needed.
func TransportFor(config *Config) (http.RoundTripper, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.New(cfg)
}

// HTTPWrappersForConfig wraps a round tripper with any relevant layered behavior from the
// config. Exposed to allow more clients that need HTTP-like behavior but then must hijack
// the underlying connection (like WebSocket or HTTP2 clients). Pure HTTP clients should use
// the higher level TransportFor or RESTClientFor methods.
func HTTPWrappersForConfig(config *Config, rt http.RoundTripper) (http.RoundTripper, error) {
	cfg, err := config.TransportConfig()
	if err != nil {
		return nil, err
	}
	return transport.HTTPWrappersForConfig(cfg, rt)
}

// TransportConfig converts a client config to an appropriate transport config.
func (c *Config) TransportConfig() (*transport.Config, error) {
	conf := &transport.Config{
		UserAgent:     c.UserAgent,
		Transport:     c.Transport,
		WrapTransport: c.WrapTransport,
		TLS: transport.TLSConfig{
			Insecure:   c.Insecure,
			ServerName: c.ServerName,
			CAFile:     c.CAFile,
			CAData:     c.CAData,
			CertFile:   c.CertFile,
			CertData:   c.CertData,
			KeyFile:    c.KeyFile,
			KeyData:    c.KeyData,
		},
		Username:    c.Username,
		Password:    c.Password,
		BearerToken: c.BearerToken,
		Impersonate: transport.ImpersonationConfig{
			UserName: c.Impersonate.UserName,
			Groups:   c.Impersonate.Groups,
			Extra:    c.Impersonate.Extra,
		},
		Dial: c.Dial,
	}

	return conf, nil
}

// Wrap adds a transport middleware function that will give the caller
// an opportunity to wrap the underlying http.RoundTripper prior to the
// first API call being made. The provided function is invoked after any
// existing transport wrappers are invoked.
func (c *Config) Wrap(fn transport.WrapperFunc) {
	c.WrapTransport = transport.Wrappers(c.WrapTransport, fn)
}
