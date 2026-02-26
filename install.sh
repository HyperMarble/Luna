#!/bin/bash

# Luna AI CA Agent - Install/Update Script

LUNA_VERSION="0.1.0"
LUNA_REPO="HyperMarble/Luna"
LUNA_BINARY="luna"

install_luna() {
    echo "Installing Luna v$LUNA_VERSION..."
    
    # Get latest release URL
    LUNA_URL=$(curl -sL "https://api.github.com/repos/$LUNA_REPO/releases/latest" | grep -o "https.*luna-darwin-arm64" | head -1)
    
    if [ -z "$LUNA_URL" ]; then
        echo "Error: Could not find Luna binary"
        exit 1
    fi
    
    # Download
    curl -L "$LUNA_URL" -o "$LUNA_BINARY"
    chmod +x "$LUNA_BINARY"
    
    # Move to /usr/local/bin
    sudo mv "$LUNA_BINARY" /usr/local/bin/luna
    
    echo "Luna installed successfully!"
    echo "Run 'luna' to start"
}

update_luna() {
    install_luna
}

# Check if argument is provided
case "${1:-install}" in
    install)
        install_luna
        ;;
    update)
        update_luna
        ;;
    *)
        echo "Usage: curl -sL https://hypermarble.github.io/Luna/install.sh | bash"
        echo "   or: curl -sL https://hypermarble.github.io/Luna/install.sh | bash -s update"
        ;;
esac
