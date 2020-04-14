#!/bin/sh

echo 'Stopping daemon'
launchctl unload ~/Library/LaunchAgents/com.github.dncrht.kwik.plist

echo 'Removing launch script'
rm ~/Library/LaunchAgents/com.github.dncrht.kwik.plist
