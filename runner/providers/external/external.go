package external

import (
	"context"

	"github.com/cloudbase/garm/config"
	"github.com/cloudbase/garm/runner/common"
	v010 "github.com/cloudbase/garm/runner/providers/v0.1.0"
	v011 "github.com/cloudbase/garm/runner/providers/v0.1.1"
)

// NewProvider selects the provider based on the interface version
func NewProvider(ctx context.Context, cfg *config.Provider, controllerID string) (common.Provider, error) {
	switch cfg.External.InterfaceVersion {
	case "v0.1.0":
		return v010.NewProvider(ctx, cfg, controllerID)
	case "v0.1.1":
		return v011.NewProvider(ctx, cfg, controllerID)
	default:
		// No version declared, assume legacy
		return v010.NewProvider(ctx, cfg, controllerID)
	}
}
