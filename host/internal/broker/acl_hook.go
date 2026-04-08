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
	"bytes"
	"log"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
)

// ACLHook enforces per-device topic access control.
type ACLHook struct {
	mqtt.HookBase
}

// ID returns the hook identifier.
func (h *ACLHook) ID() string {
	return "acl-hook"
}

// Provides indicates which hook methods this hook implements.
func (h *ACLHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnACLCheck,
	}, []byte{b})
}

// OnACLCheck validates that a client can only access its own device topics.
func (h *ACLHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	// Internal inline subscriptions are always allowed
	if cl.Net.Inline {
		return true
	}

	// Device clients can only access cdd/{their-device-id}/*
	deviceID := string(cl.ID)
	if deviceID == "" {
		log.Printf("[acl] DENIED: empty client ID for topic=%s write=%v", topic, write)
		return false
	}

	parts := strings.SplitN(topic, "/", 3)
	if len(parts) < 2 || parts[0] != "cdd" {
		log.Printf("[acl] DENIED: bad topic format topic=%s clientID=%s write=%v", topic, deviceID, write)
		return false
	}
	allowed := parts[1] == deviceID
	if !allowed {
		log.Printf("[acl] DENIED: topic deviceID=%q != clientID=%q topic=%s write=%v", parts[1], deviceID, topic, write)
	}
	return allowed
}
