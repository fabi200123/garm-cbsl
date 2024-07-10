package external

import (
	"context"

	"github.com/cloudbase/garm/config"
	"github.com/cloudbase/garm/params"
	"github.com/cloudbase/garm/runner/common"
	v010 "github.com/cloudbase/garm/runner/providers/v0.1.0"
	v011 "github.com/cloudbase/garm/runner/providers/v0.1.1"
)

// NewProvider selects based on the version, which provider to create.
func NewProvider(ctx context.Context, cfg *config.Provider, controlerInfo params.ControllerInfo) (common.Provider, error) {
	switch cfg.External.InterfaceVersion {
	case "v0.1.0":
		return v010.NewProvider(ctx, cfg, controlerInfo.ControllerID.String())
	case "v0.1.1":
		return v011.NewProvider(ctx, cfg, controlerInfo)
	default:
		// No version declared, assume legacy
		return v010.NewProvider(ctx, cfg, controlerInfo.ControllerID.String())
	}
}
