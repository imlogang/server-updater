#!/bin/bash
STABLE_VERSION=$(curl -s https://piston-meta.mojang.com/mc/game/version_manifest_v2.json | jq -r '.latest.release')

function notify_discord() {
    echo "Notifying Discord"
    local message="$1"
    curl -s -X POST "${DISCORD_WEBHOOK_LINK}" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        --data "{\"content\": \"${message}\"}"
}

function notify_server() {
    echo "Notifiting sever"
    local message="$1"
    curl -s -X POST "http://minecraftdell2.logangodsey.com/api/client/servers/${serverValue}/command" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${PteroToken}" \
        -d "{
            \"command\": \"say ${message}\"
        }"
}

function update_discord_and_server() {
    notify_discord "The server will be updated to version ${STABLE_VERSION} in 5 minutes. Please update your client!"
    notify_server "The server will be updating in 5 minutes to ${STABLE_VERSION}!"
    
    sleep 4m 

    notify_discord "The server will be updated to version ${STABLE_VERSION} in 1 minute. Please update your client!"
    notify_server "The server will be updating in 1 minute to ${STABLE_VERSION}!"
    
    sleep 55

    notify_server "The server will be updating in 5 seconds to ${STABLE_VERSION}!"
    sleep 5
}

function update_server() {
    echo "Updating server"
    curl -s -X POST "http://minecraftdell2.logangodsey.com/api/client/servers/${serverValue}/settings/reinstall" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${PteroToken}"
}

update_discord_and_server
update_server
