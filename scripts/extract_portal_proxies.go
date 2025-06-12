package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"

	"github.com/julienrbrt/l2fi/config"
)

type ChainConfig struct {
	Name      string            `toml:"name"`
	DAType    string            `toml:"data_availability_type"`
	Addresses map[string]string `toml:"addresses"`
}

func main() {
	dir := flag.String("dir", ".", "Directory containing opstack chain configs")
	flag.Parse()

	files, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var entries []config.Chain
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".toml") {
			continue
		}
		path := filepath.Join(*dir, f.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", path, err)
		}
		var cfg ChainConfig
		if err := toml.Unmarshal(b, &cfg); err != nil {
			log.Fatal("Failed to unmarshal TOML:", err)
		}

		if cfg.DAType != "eth-da" { // only include chains with eth-da DA type
			continue
		}

		proxy, ok := cfg.Addresses["OptimismPortalProxy"]
		if !ok {
			continue
		}

		entries = append(entries, config.Chain{
			Name:        strings.ToLower(strings.ReplaceAll(cfg.Name, " ", "-")),
			DisplayName: strings.Title(cfg.Name),
			OpStackConfig: &config.OpStackConfig{
				OptimismPortalAddress: proxy,
			},
		})
	}

	bz, err := yaml.Marshal(&config.AppConfig{
		RPCURL: "https://eth.llamarpc.com",
		Chains: entries,
	})
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	if err := os.WriteFile("config.yaml", bz, 0o644); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}
