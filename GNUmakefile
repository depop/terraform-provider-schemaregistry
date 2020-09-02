TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=schemaregistry
TEST_COUNT?=1
ACCTEST_TIMEOUT?=1m
ACCTEST_PARALLELISM?=20
BINARY=terraform-provider-${PKG_NAME}
VERSION=0.1
OS_ARCH=darwin_amd64

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test $(TEST) $(TESTARGS) -timeout=120s -parallel=4

testacc: fmtcheck
	@if [ "$(TESTARGS)" = "-run=TestAccXXX" ]; then \
		echo ""; \
		echo "Error: Skipping example acceptance testing pattern. Update TESTARGS to match the test naming in the relevant *_test.go file."; \
		echo ""; \
		echo "For example if updating registry/registry_subject_schema.go, use the test names in registry/registry_subject_schema_test.go starting with TestAcc and up to the underscore:"; \
		echo "make testacc TESTARGS='-run=TestAccAWSAcmCertificate_'"; \
		exit 1; \
	fi
	TF_ACC=1 go test ./$(PKG_NAME) -v -count $(TEST_COUNT) -parallel $(ACCTEST_PARALLELISM) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT)

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

gencheck:
	@echo "==> Checking generated source code..."
	@$(MAKE) gen
	@git diff --compact-summary --exit-code || \
		(echo; echo "Unexpected difference in directories after code generation. Run 'make gen' command and commit."; exit 1)

depscheck:
	@echo "==> Checking source code with go mod tidy..."
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum || \
		(echo; echo "Unexpected difference in go.mod/go.sum files. Run 'go mod tidy' command or revert any go.mod/go.sum changes and commit."; exit 1)
	@echo "==> Checking source code with go mod vendor..."
	@go mod vendor
	@git diff --compact-summary --exit-code -- vendor || \
		(echo; echo "Unexpected difference in vendor/ directory. Run 'go mod vendor' command or revert any go.mod/go.sum/vendor changes and commit."; exit 1)

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

install:
	go build -o ${BINARY}
	mkdir -p ~/.terraform.d/plugins/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${OS_ARCH}/${BINARY}


.PHONY: build gen test testacc fmt fmtcheck lint test-compile depscheck docscheck
