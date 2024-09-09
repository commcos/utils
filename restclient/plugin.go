/*

Copyright 2021-2022 This Project Authors.

Author:  seanchann <seanchann@foxmail.com>

See docs/ for more information about the  project.

*/

package restclient

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/commcos/utils/logger"
)

type AuthProvider interface {
	// WrapTransport allows the plugin to create a modified RoundTripper that
	// attaches authorization headers (or other info) to requests.
	WrapTransport(http.RoundTripper) http.RoundTripper
	// Login allows the plugin to initialize its configuration. It must not
	// require direct user interaction.
	Login() error
}

// Factory generates an AuthProvider plugin.
//
//	clusterAddress is the address of the current cluster.
//	config is the initial configuration for this plugin.
//	persister allows the plugin to save updated configuration.
type Factory func(clusterAddress string, config map[string]string, persister AuthProviderConfigPersister) (AuthProvider, error)

// AuthProviderConfigPersister allows a plugin to persist configuration info
// for just itself.
type AuthProviderConfigPersister interface {
	Persist(map[string]string) error
}

// All registered auth provider plugins.
var pluginsLock sync.Mutex
var plugins = make(map[string]Factory)

func RegisterAuthProviderPlugin(name string, plugin Factory) error {
	pluginsLock.Lock()
	defer pluginsLock.Unlock()
	if _, found := plugins[name]; found {
		return fmt.Errorf("Auth Provider Plugin %q was registered twice", name)
	}
	logger.Log(logger.DebugLevel, "Registered Auth Provider Plugin %q", name)
	plugins[name] = plugin
	return nil
}
