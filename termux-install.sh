#!/bin/bash
# PbxSpamGo — Termux Install Script
# Run: bash termux-install.sh

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  PbxSpamGo — Termux Setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Update packages
pkg update -y && pkg upgrade -y

# Install Go
pkg install golang git -y

# Check Go version
go version

# Clone repo (change to your repo URL)
git clone https://github.com/badmunda05/SpamGo.git
cd SpamGo

# Create .env file
echo ""
echo "Enter your details:"
read -p "APP_ID: " APP_ID
read -p "APP_HASH: " APP_HASH
read -p "OWNER_ID: " OWNER_ID
read -p "BOT_TOKEN1: " BOT_TOKEN1
read -p "MONGO_URL (press enter to skip): " MONGO_URL

cat > .env << ENVEOF
APP_ID=$APP_ID
APP_HASH=$APP_HASH
OWNER_ID=$OWNER_ID
BOT_TOKEN1=$BOT_TOKEN1
MONGO_URL=$MONGO_URL
START_PIC=https://files.tgvibes.online/5JreGgKB.jpg
HELP_PIC=https://files.tgvibes.online/5JreGgKB.jpg
ENVEOF

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Building..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━"
go mod tidy
go build -o pbxspamgo .

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  Starting bot..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./pbxspamgo
