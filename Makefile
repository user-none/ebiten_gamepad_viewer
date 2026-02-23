BUILD_DIR     := build
FRAMEWORK     := $(BUILD_DIR)/Iosgamepadviewer.xcframework
DESKTOP_BIN   := $(BUILD_DIR)/ebiten_gamepad_viewer

.PHONY: all desktop ios clean

all: desktop ios

desktop: $(DESKTOP_BIN)

ios: $(FRAMEWORK)

$(DESKTOP_BIN):
	mkdir -p $(BUILD_DIR)
	go build -o $(DESKTOP_BIN) .

$(FRAMEWORK):
	mkdir -p $(BUILD_DIR)
	go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile bind -target ios -o $(FRAMEWORK) ./ios/

clean:
	rm -rf $(BUILD_DIR)
