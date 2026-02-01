.PHONY: build install clean test installers installer-deb installer-windows installer-dmg

VERSION := 0.1.0
BINARY := transport
BUILD_DIR := build
DEB_DIR := $(BUILD_DIR)/deb/$(BINARY)_$(VERSION)_amd64
DMG_DIR := $(BUILD_DIR)/dmg_content

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

# Build flags
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

all: build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY) ./cmd/transport/

install: build
	cp $(BINARY) $(HOME)/.local/bin/

uninstall:
	rm -f $(HOME)/.local/bin/$(BINARY)

clean:
	$(GOCLEAN)
	rm -f $(BINARY)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v ./...

# Cross-compilation
build-all: build-linux build-darwin build-windows

build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/transport/
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/transport/

build-darwin:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 ./cmd/transport/
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 ./cmd/transport/

build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/transport/

run:
	$(GOBUILD) -o $(BINARY) ./cmd/transport/
	./$(BINARY) $(ARGS)

# Installers
installers: installer-deb installer-windows installer-dmg

installer-deb: build-linux
	mkdir -p $(DEB_DIR)/DEBIAN
	mkdir -p $(DEB_DIR)/usr/local/bin
	cp $(BUILD_DIR)/$(BINARY)-linux-amd64 $(DEB_DIR)/usr/local/bin/$(BINARY)
	chmod 755 $(DEB_DIR)/usr/local/bin/$(BINARY)
	echo "Package: $(BINARY)\nVersion: $(VERSION)\nSection: utils\nPriority: optional\nArchitecture: amd64\nMaintainer: Transport\nDescription: Transport application" > $(DEB_DIR)/DEBIAN/control
	dpkg-deb --build $(DEB_DIR) $(BUILD_DIR)/$(BINARY)_$(VERSION)_amd64.deb

installer-windows: build-windows
	sed 's/VERSION/$(VERSION)/g' installer/windows.nsi.template > $(BUILD_DIR)/$(BINARY).nsi
	cd $(BUILD_DIR) && makensis $(BINARY).nsi

installer-dmg: build-darwin
	mkdir -p $(DMG_DIR)
	cp $(BUILD_DIR)/$(BINARY)-darwin-arm64 $(DMG_DIR)/$(BINARY)
	cp $(BUILD_DIR)/$(BINARY)-darwin-amd64 $(DMG_DIR)/$(BINARY)-intel
	chmod 755 $(DMG_DIR)/$(BINARY) $(DMG_DIR)/$(BINARY)-intel
	echo "Transport v$(VERSION)\n\nCopy 'transport' (Apple Silicon) or 'transport-intel' (Intel) to /usr/local/bin/" > $(DMG_DIR)/README.txt
	genisoimage -V "Transport" -D -R -apple -no-pad -o $(BUILD_DIR)/$(BINARY)_$(VERSION)_macos.dmg $(DMG_DIR)/
