package jsonrpc

import (
	"context"
	"testing"

	"github.com/komari-monitor/komari/pkg/rpc"
)

func TestAgentControlMethodsAreSensitive(t *testing.T) {
	for _, method := range []string{
		"admin:getAgentStartupConfig",
		"admin:switchAgentVersion",
	} {
		if !rpc.IsSensitive(method) {
			t.Fatalf("%s must require sensitive-operation verification", method)
		}
	}
}

func TestSwitchAgentVersionRejectsMissingSensitiveVerification(t *testing.T) {
	meta := &rpc.ContextMeta{
		Principal:  rpc.NewUserPrincipal(""),
		Permission: rpc.RoleAdmin,
	}
	req := &rpc.JsonRpcRequest{
		Version: rpc.RPC_VERSION,
		ID:      1,
		Method:  "admin:switchAgentVersion",
		Params:  map[string]any{"uuid": "node", "version": "1.5.*"},
	}
	response := dispatchWithSensitive(context.Background(), nil, meta, req)
	if response.Error == nil || response.Error.Code != rpc.PermissionDenied {
		t.Fatalf("missing 2FA response = %#v", response)
	}
}
