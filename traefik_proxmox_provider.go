// Package traefik_proxmox_provider is a plugin to use a proxmox cluster as a provider for Traefik.
package traefik_proxmox_provider

import (
	"context"
	"encoding/json"

	"github.com/NX211/traefik-proxmox-provider/provider"
)

// ClusterConfig represents a single Proxmox cluster/server configuration.
// This allows configuring multiple independent Proxmox API endpoints from a single plugin instance.
type ClusterConfig struct {
	Name          string `json:"name" yaml:"name" toml:"name"`
	ApiEndpoint   string `json:"apiEndpoint" yaml:"apiEndpoint" toml:"apiEndpoint"`
	ApiTokenId    string `json:"apiTokenId" yaml:"apiTokenId" toml:"apiTokenId"`
	ApiToken      string `json:"apiToken" yaml:"apiToken" toml:"apiToken"`
	ApiLogging    string `json:"apiLogging" yaml:"apiLogging" toml:"apiLogging"`
	ApiValidateSSL string `json:"apiValidateSSL" yaml:"apiValidateSSL" toml:"apiValidateSSL"`
}

// Config the plugin configuration.
// Either the legacy single-cluster fields (ApiEndpoint/ApiTokenId/ApiToken/...) can be used,
// or the multi-cluster `Clusters` list. If `Clusters` is empty, the legacy fields are used.
type Config struct {
	PollInterval   string          `json:"pollInterval" yaml:"pollInterval" toml:"pollInterval"`
	ApiEndpoint    string          `json:"apiEndpoint" yaml:"apiEndpoint" toml:"apiEndpoint"`
	ApiTokenId     string          `json:"apiTokenId" yaml:"apiTokenId" toml:"apiTokenId"`
	ApiToken       string          `json:"apiToken" yaml:"apiToken" toml:"apiToken"`
	ApiLogging     string          `json:"apiLogging" yaml:"apiLogging" toml:"apiLogging"`
	ApiValidateSSL string          `json:"apiValidateSSL" yaml:"apiValidateSSL" toml:"apiValidateSSL"`
	Clusters       []ClusterConfig `json:"clusters" yaml:"clusters" toml:"clusters"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	cfg := provider.CreateConfig()
	return &Config{
		PollInterval:   cfg.PollInterval,
		ApiEndpoint:    cfg.ApiEndpoint,
		ApiTokenId:     cfg.ApiTokenId,
		ApiToken:       cfg.ApiToken,
		ApiLogging:     cfg.ApiLogging,
		ApiValidateSSL: cfg.ApiValidateSSL,
		Clusters:       nil,
	}
}

// Provider a plugin.
type Provider struct {
	provider *provider.Provider
}

// New creates a new Provider plugin.
func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	providerConfig := &provider.Config{
		PollInterval:   config.PollInterval,
		ApiEndpoint:    config.ApiEndpoint,
		ApiTokenId:     config.ApiTokenId,
		ApiToken:       config.ApiToken,
		ApiLogging:     config.ApiLogging,
		ApiValidateSSL: config.ApiValidateSSL,
	}

	// If multi-cluster configuration is provided, map it through to the provider.
	if len(config.Clusters) > 0 {
		providerConfig.Clusters = make([]provider.ClusterConfig, 0, len(config.Clusters))
		for _, c := range config.Clusters {
			pc := provider.ClusterConfig{
				Name:           c.Name,
				ApiEndpoint:    c.ApiEndpoint,
				ApiTokenId:     c.ApiTokenId,
				ApiToken:       c.ApiToken,
				ApiLogging:     c.ApiLogging,
				ApiValidateSSL: c.ApiValidateSSL,
			}
			providerConfig.Clusters = append(providerConfig.Clusters, pc)
		}
	}

	innerProvider, err := provider.New(ctx, providerConfig, name)
	if err != nil {
		return nil, err
	}

	return &Provider{
		provider: innerProvider,
	}, nil
}

// Init initializes the provider.
func (p *Provider) Init() error {
	return p.provider.Init()
}

// Provide creates and sends dynamic configuration.
func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	return p.provider.Provide(cfgChan)
}

// Stop the provider.
func (p *Provider) Stop() error {
	return p.provider.Stop()
} 