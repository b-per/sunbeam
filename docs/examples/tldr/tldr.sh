#!/bin/bash

# @sunbeam.title List TLDR Pages

set -euo pipefail

PLATFORM="${1:-osx}"

# shellcheck disable=SC2016
tldr --list --platform="$PLATFORM" | sunbeam query --arg platform="$PLATFORM" -R '{
    title: .,
    detail: {
      command: ["tldr", "--raw", "--platform", $platform, .]
    },
    actions: [
      {
        type: "open",
        title: "Open in browser",
        target: "https://tldr.inbrowser.app/pages/common/\(. | @uri)",
      }
    ]
}' | sunbeam list --json --title "tldr" --show-detail