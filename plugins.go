package gott

import (
	"net"
	"os"
	"path"
	"plugin"

	"github.com/oimyounis/gott/utils"
)

// gottPlugin represents an plugin
type gottPlugin struct {
	name                string
	plug                *plugin.Plugin
	onSocketOpen        func(conn net.Conn) bool
	onBeforeConnect     func(clientID, username, password string) bool
	onConnect           func(clientID, username, password string) bool
	onMessage           func(clientID, username string, topic, payload []byte, dup, qos byte, retain bool)
	onBeforePublish     func(clientID, username string, topic, payload []byte, dup, qos byte, retain bool) bool
	onPublish           func(clientID, username string, topic, payload []byte, dup, qos byte, retain bool)
	onBeforeSubscribe   func(clientID, username string, topic []byte, qos byte) bool
	onSubscribe         func(clientID, username string, topic []byte, qos byte)
	onBeforeUnsubscribe func(clientID, username string, topic []byte) bool
	onUnsubscribe       func(clientID, username string, topic []byte)
	onDisconnect        func(clientID, username string, graceful bool)
}

func (b *Broker) bootstrapPlugins() {
	if !utils.PathExists(pluginDir) {
		b.logger.Info().Msg("Plugins directory does not exist. Creating a new one.")
		if err := os.Mkdir(pluginDir, 0775); err != nil {
			b.logger.Error().Err(err).Msg("Failed to create the plugins directory")
		}
	}
	for _, pstring := range b.config.Plugins {
		p, err := plugin.Open(path.Join(pluginDir, pstring))
		if err != nil {
			b.logger.Error().Err(err).Str("name", pstring).Msg("Skipping loading plugin")
			continue
		}

		pluginObj := gottPlugin{
			name: pstring,
			plug: p,
		}

		if bootstrap, err := p.Lookup("Bootstrap"); err == nil {
			if b, ok := bootstrap.(func()); ok {
				b()
			}
		}

		h, err := p.Lookup("OnSocketOpen")
		if err == nil {
			f, ok := h.(func(conn net.Conn) bool)
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnSocketOpen")
			if ok {
				pluginObj.onSocketOpen = f
			}
		}

		h, err = p.Lookup("OnBeforeConnect")
		if err == nil {
			f, ok := h.(func(clientID, username, password string) bool)
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnBeforeConnect")
			if ok {
				pluginObj.onBeforeConnect = f
			}
		}

		if h, err = p.Lookup("OnConnect"); err == nil {
			f, ok := h.(func(clientID, username, password string) bool)
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnConnect")
			if ok {
				pluginObj.onConnect = f
			}
		}

		if h, err = p.Lookup("OnMessage"); err == nil {
			f, ok := h.(func(clientID, username string, topic, payload []byte, dup, qos byte, retain bool))
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnMessage")
			if ok {
				pluginObj.onMessage = f
			}
		}

		if h, err = p.Lookup("OnBeforePublish"); err == nil {
			f, ok := h.(func(clientID, username string, topic, payload []byte, dup, qos byte, retain bool) bool)
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnBeforePublish")
			if ok {
				pluginObj.onBeforePublish = f
			}
		}

		if h, err = p.Lookup("OnPublish"); err == nil {
			f, ok := h.(func(clientID, username string, topic, payload []byte, dup, qos byte, retain bool))
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnPublish")
			if ok {
				pluginObj.onPublish = f
			}
		}

		if h, err = p.Lookup("OnBeforeSubscribe"); err == nil {
			f, ok := h.(func(clientID, username string, topic []byte, qos byte) bool)
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnBeforeSubscribe")
			if ok {
				pluginObj.onBeforeSubscribe = f
			}
		}

		if h, err = p.Lookup("OnSubscribe"); err == nil {
			f, ok := h.(func(clientID, username string, topic []byte, qos byte))
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnSubscribe")
			if ok {
				pluginObj.onSubscribe = f
			}
		}

		if h, err = p.Lookup("OnBeforeUnsubscribe"); err == nil {
			f, ok := h.(func(clientID, username string, topic []byte) bool)
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnBeforeUnsubscribe")
			if ok {
				pluginObj.onBeforeUnsubscribe = f
			}
		}

		if h, err = p.Lookup("OnUnsubscribe"); err == nil {
			f, ok := h.(func(clientID, username string, topic []byte))
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnUnsubscribe")
			if ok {
				pluginObj.onUnsubscribe = f
			}
		}

		if h, err = p.Lookup("OnDisconnect"); err == nil {
			f, ok := h.(func(clientID, username string, graceful bool))
			GOTT.logger.Debug().
				Str("name", pstring).
				Bool("loaded", ok).
				Msg("plugin loader OnDisconnect")
			if ok {
				pluginObj.onDisconnect = f
			}
		}

		b.plugins = append(b.plugins, pluginObj)

		GOTT.logger.Debug().
			Str("name", pstring).
			Msg("plugin loaded")

	}
}

func (b *Broker) invokeOnSocketOpen(conn net.Conn) bool {
	for _, p := range b.plugins {
		if p.onSocketOpen != nil {
			if !p.onSocketOpen(conn) {
				return false
			}
		}
	}

	return true
}

func (b *Broker) invokeOnBeforeConnect(clientID, username, password string) bool {
	for _, p := range b.plugins {
		if p.onBeforeConnect != nil {
			if !p.onBeforeConnect(clientID, username, password) {
				return false
			}
		}
	}

	return true
}

func (b *Broker) invokeOnConnect(clientID, username, password string) bool {
	for _, p := range b.plugins {
		if p.onConnect != nil {
			if !p.onConnect(clientID, username, password) {
				return false
			}
		}
	}

	return true
}

func (b *Broker) invokeOnMessage(clientID, username string, topic, payload []byte, dup, qos byte, retain bool) {
	for _, p := range b.plugins {
		if p.onMessage != nil {
			p.onMessage(clientID, username, topic, payload, dup, qos, retain)
		}
	}
}

func (b *Broker) invokeOnBeforePublish(clientID, username string, topic, payload []byte, dup, qos byte, retain bool) bool {
	for _, p := range b.plugins {
		if p.onBeforePublish != nil {
			if !p.onBeforePublish(clientID, username, topic, payload, dup, qos, retain) {
				return false
			}
		}
	}

	return true
}

func (b *Broker) invokeOnPublish(clientID, username string, topic, payload []byte, dup, qos byte, retain bool) {
	for _, p := range b.plugins {
		if p.onPublish != nil {
			p.onPublish(clientID, username, topic, payload, dup, qos, retain)
		}
	}
}

func (b *Broker) invokeOnBeforeSubscribe(clientID, username string, topic []byte, qos byte) bool {
	for _, p := range b.plugins {
		if p.onBeforeSubscribe != nil {
			if !p.onBeforeSubscribe(clientID, username, topic, qos) {
				return false
			}
		}
	}

	return true
}

func (b *Broker) invokeOnSubscribe(clientID, username string, topic []byte, qos byte) {
	for _, p := range b.plugins {
		if p.onSubscribe != nil {
			p.onSubscribe(clientID, username, topic, qos)
		}
	}
}

func (b *Broker) invokeOnBeforeUnsubscribe(clientID, username string, topic []byte) bool {
	for _, p := range b.plugins {
		if p.onBeforeUnsubscribe != nil {
			if !p.onBeforeUnsubscribe(clientID, username, topic) {
				return false
			}
		}
	}

	return true
}

func (b *Broker) invokeOnUnsubscribe(clientID, username string, topic []byte) {
	for _, p := range b.plugins {
		if p.onUnsubscribe != nil {
			p.onUnsubscribe(clientID, username, topic)
		}
	}
}

func (b *Broker) invokeOnDisconnect(clientID, username string, graceful bool) {
	for _, p := range b.plugins {
		if p.onDisconnect != nil {
			p.onDisconnect(clientID, username, graceful)
		}
	}
}
