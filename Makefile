# binary aliases
TF = terraform -chdir=./infrastructure/
SHELL:=/bin/bash -O globstar

# variables
PROJECT=littleurl-api
ENVIRONMENT=dev

##@ Dependencies
.PHONY: install install-tools

install: install-tools ## Install all dependencies

install-tools: ## Install tooling
	asdf install

##@ Lint
.PHONY: lint lint-fix lint-tf lint-tf-fix

lint: lint-tf ## Run all linters

lint-fix: lint-tf-fix ## Attempt to automagically fiax all linting issues

lint-tf: ## Lint terraform code
	@echo "Linting Terraform"
	$(TF) validate
	$(TF) init -backend=false
	$(TF) fmt -check -recursive

lint-tf-fix: ## Auto fix terraform linting errors
	@echo "Fixing Terraform linting"
	$(TF) init -backend=false
	$(TF) fmt -recursive

##@ Build
.PHONY: build build-functions build-templates

build: build-templates build-functions ## Build everything

build-functions: ## Build lambda functions
	@echo "Building lambda functions"
	rm -rf ./build/functions/*
	GOOS=linux GOARCH=arm64 go build -o=./build/functions/ -tags=lambda.norpc ./functions/... 
	cd ./build/functions \
		&& for i in *; do mv "$$i" bootstrap && zip "$$i.zip" bootstrap; done \
		&& find . -type f ! -name '*.zip' -delete

build-templates: ## Compile MJML to HTML
	@echo "Compiling MJML"
	cd ./internal/templates && \
	for i in **/*.mjml; do ../../node_modules/.bin/mjml $$i -o $${i%.mjml*}.html; done

##@ Deployment
.PHONY: deploy upload-functions tf-init tf-plan tf-apply

deploy: tf-plan tf-apply build upload-functions ## Full deployment

upload-functions: ## Upload functions to S3
	@echo "Uploading lambda deployment packages"
	$(TF) workspace select ${ENVIRONMENT}
	aws s3 cp ./build/functions s3://$$($(TF) output -raw functions_bucket) --recursive

tf-init: ## Initialise terraform
	@echo "Initialising terraform"
	$(TF) init

tf-plan: ## Plan terraform changeset
	@echo "Creating terraform plan"
	$(TF) workspace select ${ENVIRONMENT}
	$(TF) plan -input=false -out=./${PROJECT}.tfplan

tf-apply: ## Apply terraform changeset
	@echo "Applying terraform plan"
	$(TF) workspace select ${ENVIRONMENT}
	$(TF) apply -input=false -auto-approve ./${PROJECT}.tfplan

##@ Helpers
.PHONY: help

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
