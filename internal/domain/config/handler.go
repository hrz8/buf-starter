package config

import (
	"context"

	"connectrpc.com/connect"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// ConfigProvider provides configuration values.
type ConfigProvider interface {
	GetDashboardBrandingName() string
	GetAuthServerBrandingName() string
}

// Handler implements the ConfigService.
type Handler struct {
	cfg ConfigProvider
}

// NewHandler creates a new config handler.
func NewHandler(cfg ConfigProvider) *Handler {
	return &Handler{cfg: cfg}
}

// GetPublicConfig returns public configuration including branding.
func (h *Handler) GetPublicConfig(
	ctx context.Context,
	req *connect.Request[altalunev1.GetPublicConfigRequest],
) (*connect.Response[altalunev1.GetPublicConfigResponse], error) {
	response := &altalunev1.GetPublicConfigResponse{
		Branding: &altalunev1.BrandingConfig{
			DashboardName:  h.cfg.GetDashboardBrandingName(),
			AuthServerName: h.cfg.GetAuthServerBrandingName(),
		},
	}
	return connect.NewResponse(response), nil
}
