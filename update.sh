#!/bin/bash

function notify_discord() {
    echo "Notifying Discord"
    local message="$1"
    curl -s -X POST "${DISCORD_WEBHOOK_LINK}" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        --data "{\"content\": \"${message}\"}"
}

function notify_minecraft_server() {
    echo "Notifiting sever"
    local message="$1"
    curl -s -X POST "http://minecraftdell2.logangodsey.com/api/client/servers/${serverValue}/command" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${PteroToken}" \
        -d "{\"command\": \"say ${message}\"}"
}

function update_discord_and_server() {
    notify_discord "The server will be updated to version ${STABLE_VERSION} in 5 minutes. Please update your client!"
    notify_minecraft_server "The server will be updating in 5 minutes to ${STABLE_VERSION}!"
    
    sleep 4m 

    notify_discord "The server will be updated to version ${STABLE_VERSION} in 1 minute. Please update your client!"
    notify_minecraft_server "The server will be updating in 1 minute to ${STABLE_VERSION}!"
    
    sleep 55

    notify_discord "The server will be updated to version ${STABLE_VERSION} in 5 seconds. Please update your client!"
    notify_minecraft_server "The server will be updating in 5 seconds to ${STABLE_VERSION}!"
    sleep 5
    notify_discord "The server is now updated to version ${STABLE_VERSION}."
}

function update_discord() {
    local server="$1"
    notify_discord "The ${server} server will be updated in 5 minutes. Please update your client!"
    sleep 4m 
    notify_discord "The server will be updated in 1 minute. Please update your client!"
    sleep 55
    notify_discord "The server will be updated in 5 seconds. Please update your client!"
    sleep 5
    notify_discord "The server is now updated."
}

function update_reinstall_server() {
    echo "Updating server by reinstalling"
    curl -s -X POST "http://minecraftdell2.logangodsey.com/api/client/servers/${serverValue}/settings/reinstall" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${PteroToken}"
}

function update_power_server() {
    local power_state="$1"
    echo "Updating server by ${power_state}"
    curl -s -X POST "http://minecraftdell2.logangodsey.com/api/client/servers/${serverValue}/power" \
        -H "Accept: application/json" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${PteroToken}" \
        -d '{"signal": "'"$power_state"'"}'
}


echo "serverValue is: $serverValue"

case "$serverValue" in
  "a6615eb7")
    echo "Running commands for SMP Vanilla"
    export DISCORD_WEBHOOK_LINK="${DISCORD_WEBHOOK_LINK_MINECRAFT}"
    STABLE_VERSION=$(curl -s https://piston-meta.mojang.com/mc/game/version_manifest_v2.json | jq -r '.latest.release')
    update_discord_and_server
    update_reinstall_server
    update_power_server "start"
    ;;

  "6b774df5")
    echo "Running commands Satisfactory"
    export DISCORD_WEBHOOK_LINK="${DISCORD_WEBHOOK_LINK_SATISFACTORY}"
    update_discord
    update_power_server "restart"
    ;;
  *)
    echo "Unknown serverValue: $serverValue"
    ;;
esac