package agent

import (
	"errors"
	"testing"
)

func TestSwitchAgentVersionRequiresCapability(t *testing.T) {
	const uuid = "switch-version-node"
	SetPresence(uuid, 1, true)
	MarkV2Client(uuid)
	t.Cleanup(func() {
		SetPresence(uuid, 1, false)
		DeleteConnectedClients(uuid)
	})

	if err := SwitchAgentVersion(uuid, "1.5.*"); !errors.Is(err, ErrSwitchVersionUnsupported) {
		t.Fatalf("unsupported agent error = %v", err)
	}
	SetV2Capabilities(uuid, []string{"switch_version"})
	if err := SwitchAgentVersion(uuid, "1.5.*"); err != nil {
		t.Fatalf("capable agent was rejected: %v", err)
	}
}
