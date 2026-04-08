// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package broker

import (
	"crypto/tls"
	"fmt"
	"log"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
)

// MessageHandler is a callback for received MQTT messages.
type MessageHandler func(topic string, payload []byte)

// Broker wraps the embedded mochi-mqtt server.
type Broker struct {
	server    *mqtt.Server
	store     *db.Store
	tlsConfig *tls.Config
	port      int
}

// New creates a new embedded MQTT broker.
func New(port int, tlsConfig *tls.Config, store *db.Store) *Broker {
	opts := &mqtt.Options{
		InlineClient: true,
	}
	return &Broker{
		server:    mqtt.New(opts),
		store:     store,
		tlsConfig: tlsConfig,
		port:      port,
	}
}

// Start initializes and starts the MQTT broker.
func (b *Broker) Start() error {
	// Add auth hook
	authHook := &AuthHook{store: b.store}
	if err := b.server.AddHook(authHook, nil); err != nil {
		return fmt.Errorf("add auth hook: %w", err)
	}

	// Add ACL hook
	aclHook := &ACLHook{}
	if err := b.server.AddHook(aclHook, nil); err != nil {
		return fmt.Errorf("add ACL hook: %w", err)
	}

	// Create TLS listener
	addr := fmt.Sprintf(":%d", b.port)
	tlsListener := listeners.NewTCP(listeners.Config{
		ID:        "tls",
		Address:   addr,
		TLSConfig: b.tlsConfig,
	})
	if err := b.server.AddListener(tlsListener); err != nil {
		return fmt.Errorf("add TLS listener: %w", err)
	}

	log.Printf("[mqtt] broker starting on port %d", b.port)
	return b.server.Serve()
}

// Stop gracefully shuts down the broker.
func (b *Broker) Stop() error {
	return b.server.Close()
}

// Publish sends a message to a topic via the inline client.
func (b *Broker) Publish(topic string, payload []byte, retain bool) error {
	log.Printf("[mqtt-broker] Publishing topic=%q retain=%v payloadLen=%d", topic, retain, len(payload))

	// Log connected clients and their subscriptions for debugging
	for id, cl := range b.server.Clients.GetAll() {
		inline := "external"
		if cl.Net.Inline {
			inline = "inline"
		}
		subs := cl.State.Subscriptions.GetAll()
		subTopics := make([]string, 0, len(subs))
		for filter := range subs {
			subTopics = append(subTopics, filter)
		}
		log.Printf("[mqtt-broker]   client=%q (%s) subscriptions=%v", id, inline, subTopics)
	}

	err := b.server.Publish(topic, payload, retain, 1)
	if err != nil {
		log.Printf("[mqtt-broker] Publish ERROR: %v", err)
	} else {
		log.Printf("[mqtt-broker] Publish OK (no error from server)")
	}
	return err
}

// Subscribe registers an internal inline handler for a topic filter.
func (b *Broker) Subscribe(filter string, handler MessageHandler) {
	if err := b.server.Subscribe(filter, 1, func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		handler(pk.TopicName, pk.Payload)
	}); err != nil {
		log.Printf("[mqtt] subscribe error for %s: %v", filter, err)
	}
}
