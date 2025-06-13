#!/usr/bin/make -f

build: build-css
	go build .

install: build-css
	go install .

build-css-watch:
	@echo "Generating CSS files..."
	npx tailwindcss@v3 -i ./static/css/input.css -o ./static/css/output.css --watch

build-css:
	@echo "Generating CSS files..."
	npx tailwindcss@v3 -i ./static/css/input.css -o ./static/css/output.css

build-nitro-abi:
	@echo "Building Nitro ABI..."
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	git clone https://github.com/offchainlabs/nitro-contracts || true
	cd nitro-contracts && \
		yarn install && \
		yarn build && \
		solc --abi ./src/bridge/Inbox.sol --base-path . --include-path node_modules/ --include-path /usr/local/lib/node_modules/ -o nitro-abi --overwrite && \
		abigen --abi=nitro-abi/Inbox.abi --pkg=l2 --type=Inbox --out=../l2/arbitrum_inbox.go
	rm -rf ./nitro-contracts

update-js:
	@echo "Updating htmx.min.js..."
	curl -L https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -o ./static/js/htmx.min.js
	@echo "Updating ethers.umd.min.js..."
	curl -L https://unpkg.com/ethers/dist/ethers.umd.min.js -o ./static/js/ethers.umd.min.js

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./... --fix