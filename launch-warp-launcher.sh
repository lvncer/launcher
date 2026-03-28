#!/bin/bash

# Warpが起動してるかチェック
if pgrep -x "Warp" > /dev/null; then
  # すでに起動してる → 新規ウィンドウ
  osascript -e 'tell application "Warp" to activate'
  sleep 0.2
  osascript -e 'tell application "System Events" to keystroke "n" using command down'
else
  # 起動してない → 普通に起動（これで1ウィンドウ出る）
  open -a Warp
  sleep 0.5
fi

# コマンド貼り付け
echo "/Users/kihhi/gitrepos/launcher/launcher" | pbcopy
osascript -e 'tell application "System Events" to keystroke "v" using command down'

# 実行
osascript -e 'tell application "System Events" to key code 36'
