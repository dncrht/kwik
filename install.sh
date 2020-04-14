#!/bin/sh

echo 'Compiling binary'
go build -o bin/kwik

echo 'Preparing launch script'
sed "s|\$HOME|$HOME|g" com.github.dncrht.kwik.plist.TEMPLATE | sed "s|\MY_DIR|$(pwd)|g" > com.github.dncrht.kwik.plist

echo 'Installing launch script'
ln -s "$(pwd)/com.github.dncrht.kwik.plist" ~/Library/LaunchAgents

echo 'Launching daemon'
launchctl load ~/Library/LaunchAgents/com.github.dncrht.kwik.plist

echo 'Open webapp'
open http://localhost:2005
