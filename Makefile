.PHONY: build-all test-all lint tidy clean

SERVICES := template-svc

build-all:
	@for svc in $(SERVICES); do \
		echo "Building $$svc..."; \
		cd services/$$svc && go build ./cmd/...; \
	done

test-all:
	@for svc in $(SERVICES); do \
		echo "Testing $$svc..."; \
		cd services/$$svc && go test -v -race -count=1 ./...; \
	done

lint:
	@for svc in $(SERVICES); do \
		echo "Linting $$svc..."; \
		cd services/$$svc && golangci-lint run ./...; \
	done

tidy:
	@for svc in $(SERVICES); do \
		echo "Tidying $$svc..."; \
		cd services/$$svc && go mod tidy; \
	done
	cd pkg && go mod tidy 2>/dev/null || true

vet:
	@for svc in $(SERVICES); do \
		echo "Vetting $$svc..."; \
		cd services/$$svc && go vet ./...; \
	done

clean:
	@for svc in $(SERVICES); do \
		rm -rf services/$$svc/bin; \
	done

.PHONY: proto-gen
proto-gen:
	@echo "Run: protoc --go_out=pkg/proto --go_opt=paths=source_relative ..."
