#!/bin/bash

APP="Warp"
CMD="LAUNCHER_CLOSE_WARP_FLOAT=1 /Users/kihhi/gitrepos/launcher/launcher"

SCREEN_WIDTH=2560
SCREEN_HEIGHT=1600

WIN_WIDTH=1000
WIN_HEIGHT=600

POS_X=$(( (SCREEN_WIDTH - WIN_WIDTH) / 2 ))
POS_Y=$(( (SCREEN_HEIGHT - WIN_HEIGHT) / 2 ))

if ! pgrep -x "$APP" > /dev/null; then
  open -a "$APP"
  sleep 0.5
fi

osascript -e 'tell application "Warp" to activate'
sleep 0.2
osascript -e 'tell application "System Events" to keystroke "n" using command down'
sleep 0.3

yabai -m window --focus recent
yabai -m window --toggle float

yabai -m window --resize abs:$WIN_WIDTH:$WIN_HEIGHT
yabai -m window --move abs:$POS_X:$POS_Y

echo "$CMD" | pbcopy
osascript -e 'tell application "System Events" to keystroke "v" using command down'
osascript -e 'tell application "System Events" to key code 36'
