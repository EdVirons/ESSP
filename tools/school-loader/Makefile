.PHONY: build run dry-run clean test

BINARY=school-loader
SQLITE_PATH ?= /home/pato/scrapper/data/schools.db
SSOT_URL ?= http://localhost:8081
TENANT_ID ?= tenant_001
BATCH_SIZE ?= 1000

build:
	go build -o $(BINARY) ./cmd/loader

run: build
	./$(BINARY) \
		--sqlite-path=$(SQLITE_PATH) \
		--ssot-url=$(SSOT_URL) \
		--tenant-id=$(TENANT_ID) \
		--batch-size=$(BATCH_SIZE)

dry-run: build
	./$(BINARY) \
		--sqlite-path=$(SQLITE_PATH) \
		--tenant-id=$(TENANT_ID) \
		--dry-run \
		--verbose

clean:
	rm -f $(BINARY)

test:
	go test -v ./...
