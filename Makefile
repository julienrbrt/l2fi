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

update-js:
	@echo "Updating htmx.min.js..."
	curl -L https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -o ./static/js/htmx.min.js
	@echo "Updating ethers.umd.min.js..."
	curl -L https://unpkg.com/ethers/dist/ethers.umd.min.js -o ./static/js/ethers.umd.min.js

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./... --fix