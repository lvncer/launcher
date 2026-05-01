# Launcher

Goで作る“自作rofi（ランチャー）”

## 機能

- アプリ自動収集（/Applications）
- アプリ名の左に `.app` のアイコン画像を小さく変換して表示
- 並び順
  - `Recent`: 直近使った5件を上に表示
  - `Apps`: 直近利用したものも含めて、全アプリを使用履歴に関係なく名前順で表示
  - `Apps` は一度に最大10件だけ表示し、選択を `↑`/`↓` で動かすとウィンドウ内がスクロールする（Raycast のイメージ）。上下に ▲/▼ の状態表示を固定表示
- モード切替（>, /, 通常）
- プレビュー表示（右側）
- 実行＆履歴保存

## 使い方

```sh
go run .
```

通常モードでは `/Applications` 配下のアプリを検索します。

- 何も付けない: アプリ検索
- `>` から始める: コマンド検索
- `/` から始める: ファイル検索用モード

アプリアイコンは、各 `.app` の `Info.plist` から `CFBundleIconFile` を読み、`.icns` を macOS の `sips` でPNGに変換してから、ターミナル上で2行のANSIカラー表示にしています。変換後のPNGは `/tmp/launcher_icon_cache` にキャッシュされます。

## ビルド方法

コードを変更した後は、必ずビルドし直してください。
`bin/launcher` はビルド済みバイナリなので、Go のコードを書き換えただけでは動きは変わりません。

```sh
mkdir -p bin
go build -o bin/launcher .
```

ビルド後に実行します。

```sh
./bin/launcher
```

テストもまとめて確認する場合は次のコマンドです。

```sh
go test ./... && mkdir -p bin && go build -o bin/launcher .
```

## Warp から起動する場合

`scripts/launch-warp-launcher.sh` は、起動前に `bin/launcher` をビルドしてから Warp の浮動ウィンドウで実行します。ウィンドウサイズは `scripts/launcher.config.sh` の `WIN_WIDTH` / `WIN_HEIGHT` で調整できます。狭いとリストが縦に入りきらないので、足りなければ高さを大きくしてください。

```sh
./scripts/launch-warp-launcher.sh
```

SKHD からこのスクリプトを実行している場合も、起動時に次のビルドが走ります。

```sh
go build -o /Users/kihhi/gitrepos/launcher/bin/launcher .
```

そのため、普段は Go のコードを変更したあとに手動ビルドし忘れても、SKHD 経由の起動時に最新化されます。

SKHD から呼ぶ場合は、起動コマンドを `scripts/launch-warp-launcher.sh` に向けてください。

もし挙動が変わっていないように見える場合は、`/tmp/launcher_usage.json` の履歴も確認してください。たとえば履歴が `Discord` だけなら、直近5件枠として `Discord` が一番上に表示されるのは正常です。

## 構成

- `main.go`: エントリーポイント
- `scripts/`: 起動用スクリプト
- `scripts/launcher.config.sh`: Warp 起動時のウィンドウサイズ
- `bin/`: ビルド成果物（git管理外）
- `internal/launcher/`: ランチャー本体
- `internal/launcher/usage.go`: 利用履歴と並び順
- `internal/launcher/model.go`: 入力、検索、キー操作
- `internal/launcher/view.go`: 画面表示
- `internal/launcher/window.go`: Warp/yabai 連携
