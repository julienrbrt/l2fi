package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/julienrbrt/l2fi/config"
	"github.com/julienrbrt/l2fi/l2"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	appConfig, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Template setup
	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = tmpl.ExecuteTemplate(w, "base.html", appConfig)
	})

	r.Get("/bsod", func(w http.ResponseWriter, r *http.Request) {
		errorType := r.URL.Query().Get("type")
		errorMessage := r.URL.Query().Get("message")

		if errorType == "" {
			errorType = "WALLET_CONNECTION_ERROR"
		}
		if errorMessage == "" {
			errorMessage = "An unknown wallet or signing error occurred."
		}

		data := map[string]any{
			"ErrorType":    errorType,
			"ErrorMessage": errorMessage,
			"ErrorCode":    "0x000000EA",
			"ErrorDetails": "WALLET_DRIVER_IRQL_NOT_LESS_OR_EQUAL",
			"ModuleName":   "web3wallet.sys",
			"ErrorAddress": "0xFFFFF800`12345678",
			"BaseAddress":  "0xFFFFF800`12300000",
			"DateStamp":    "0x12345678",
		}

		_ = tmpl.ExecuteTemplate(w, "bsod.html", data)
	})

	r.Post("/force-inclusion", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Invalid form"})
			return
		}

		fromAddress := r.FormValue("from_address")
		toAddress := r.FormValue("to_address")
		valueStr := r.FormValue("value")
		value := new(big.Int)
		if valueStr != "" {
			var ok bool
			value, ok = new(big.Int).SetString(valueStr, 10)
			if !ok {
				_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Invalid value format"})
				return
			}
		}

		data := r.FormValue("data")
		if data == "" {
			data = "0x"
		}

		gasLimitStr := r.FormValue("gas_limit")
		gasLimit, err := strconv.ParseUint(gasLimitStr, 10, 64)
		if err != nil {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Invalid gas limit"})
			return
		}

		// Validate addresses
		if !common.IsHexAddress(fromAddress) {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Invalid from address"})
			return
		}
		if !common.IsHexAddress(toAddress) {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Invalid to address"})
			return
		}

		chainNameFromForm := r.FormValue("l2")
		var selectedChainConfig *config.Chain
		for _, chainCfg := range appConfig.Chains {
			if chainCfg.Name == chainNameFromForm {
				temp := chainCfg // Create a new variable to take the address of
				selectedChainConfig = &temp
				break
			}
		}

		if selectedChainConfig == nil {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Unsupported L2: " + chainNameFromForm})
			return
		}

		var chainClient l2.L2
		chainType, typeErr := selectedChainConfig.Type()
		if typeErr != nil {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Could not determine chain type: " + typeErr.Error()})
			return
		}

		switch chainType {
		case config.OpStackChainType:
			client, clientErr := getOpStackClient(*selectedChainConfig)
			if clientErr != nil {
				_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Failed to create OpStack client: " + clientErr.Error()})
				return
			}
			chainClient = client
		case config.ArbitrumChainType:
			client, clientErr := getArbitrumClient(*selectedChainConfig)
			if clientErr != nil {
				_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Failed to create Arbitrum client: " + clientErr.Error()})
				return
			}
			chainClient = client
		default:
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": "Unsupported chain type for L2: " + chainNameFromForm})
			return
		}

		txJSON, err := chainClient.BuildForceInclusionTx(fromAddress, toAddress, data, value, gasLimit)
		if err != nil {
			_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Error": err.Error()})
			return
		}

		_ = tmpl.ExecuteTemplate(w, "result.html", map[string]any{"Tx": txJSON})
	})

	log.Printf("Server started on http://localhost:%d", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getOpStackClient(chainConfig config.Chain) (*l2.OpStackClient, error) {
	if chainConfig.OpStackConfig == nil {
		return nil, fmt.Errorf("opstack config is nil for chain %s", chainConfig.Name)
	}

	return l2.NewOpStackClient(
		chainConfig.RPCURL,
		chainConfig.OpStackConfig.OptimismPortalAddress,
	)
}

func getArbitrumClient(chainConfig config.Chain) (*l2.ArbitrumClient, error) {
	if chainConfig.Arbitrum == nil {
		return nil, fmt.Errorf("arbitrum config is nil for chain %s", chainConfig.Name)
	}

	return l2.NewArbitrumClient(
		chainConfig.RPCURL,
		chainConfig.Arbitrum.DelayedInboxAddress,
	)
}
