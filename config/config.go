package config

import (
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

// OpStackConfig holds Optimism-specific chain configuration.
type OpStackConfig struct {
	OptimismPortalAddress string `yaml:"optimism_portal_address,omitempty"`
}

// ArbitrumConfig holds Arbitrum-specific chain configuration.
type ArbitrumConfig struct {
	DelayedInboxAddress string `yaml:"delayed_inbox_address,omitempty"`
}

// Chain holds the configuration for a specific blockchain.
type Chain struct {
	Name          string          `yaml:"name"`
	DisplayName   string          `yaml:"display_name,omitempty"`
	RPCURL        string          `yaml:"rpc_url"`
	OpStackConfig *OpStackConfig  `yaml:"opstack,omitempty"`
	Arbitrum      *ArbitrumConfig `yaml:"arbitrum,omitempty"`
}

const (
	OpStackChainType  = "opstack"
	ArbitrumChainType = "arbitrum"
)

// Type returns the type of the chain based on its configuration.
func (c *Chain) Type() (string, error) {
	hasOpStack := c.OpStackConfig != nil
	hasArbitrum := c.Arbitrum != nil

	if hasOpStack && hasArbitrum {
		return "", fmt.Errorf("chain '%s' cannot have both opstack and arbitrum configurations", c.Name)
	}

	if hasOpStack {
		return OpStackChainType, nil
	}

	if hasArbitrum {
		return ArbitrumChainType, nil
	}

	return "", fmt.Errorf("chain '%s' has no opstack or arbitrum configuration defined, unable to determine type", c.Name)
}

// AppConfig holds the overall application configuration, including all chains.
type AppConfig struct {
	Chains []Chain `yaml:"chains"`
}

// LoadConfig reads the YAML configuration file from the given path and unmarshals it.
func LoadConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg AppConfig
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %w", err)
	}

	if err = validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func validateConfig(cfg *AppConfig) error {
	if len(cfg.Chains) == 0 {
		return fmt.Errorf("no chains configured")
	}

	for i, chain := range cfg.Chains {
		if chain.Name == "" {
			return fmt.Errorf("chain at index %d missing name", i)
		}

		if chain.RPCURL == "" {
			return fmt.Errorf("chain '%s' missing rpc_url", chain.Name)
		}

		chainType, err := chain.Type()
		if err != nil {
			return fmt.Errorf("chain '%s': %w", chain.Name, err) // Error determining type
		}

		switch chainType {
		case OpStackChainType:
			// L2OutputOracleAddress is optional
			if chain.OpStackConfig.OptimismPortalAddress == "" {
				return fmt.Errorf("opstack chain '%s' missing optimism_portal_address", chain.Name)
			}
		case ArbitrumChainType:
			// ArbitrumConfig is already confirmed to be non-nil by chain.Type()
			if chain.Arbitrum.DelayedInboxAddress == "" {
				return fmt.Errorf("arbitrum chain '%s' missing delayed_inbox_address", chain.Name)
			}
		default:
			// we should never reach here thanks to the Type() method
			panic(fmt.Sprintf("unexpected chain type '%s' for chain '%s'", chainType, chain.Name))
		}
	}

	// sort chain by name for consistency
	sort.Slice(cfg.Chains, func(i, j int) bool {
		return cfg.Chains[i].Name < cfg.Chains[j].Name
	})

	return nil
}
