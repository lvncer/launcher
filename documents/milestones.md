# Launcher Milestones

## 最終ゴール

Homebrew などでインストールすれば、ユーザーがリポジトリをクローンしたり、手でビルドしたり、yabai/skhd/Warp の細かい設定を調べたりしなくても、すぐランチャーを起動できる状態にする。

理想の体験:

```sh
brew install kihhi/tap/launcher
launcher
```

必要なら追加で:

```sh
launcher install-hotkey
launcher doctor
launcher config
```

## 現状の課題

- 実行にはリポジトリのクローンとGoビルドが必要。
- Warp、yabai、skhd を前提にした起動フローが強く、初見ユーザーにはセットアップ手順が多い。
- 起動スクリプトがローカル環境向けで、配布物としてはまだ汎用化されていない。
- ウィンドウ制御が yabai 前提なので、yabai を使っていないユーザーにはそのまま提供しづらい。
- 利用履歴やアイコンキャッシュが `/tmp` 固定で、ユーザー設定・永続化・削除の考え方がまだ弱い。
- アプリ探索が `/Applications` 中心で、`~/Applications` やSpotlight相当の探索には未対応。

## 方針

まずは「開発者が自分で使いやすい構成」に整え、その後に「外部ユーザーが安全に導入できる配布物」に育てる。yabai/skhd/Warp 連携は便利機能として残しつつ、コアの `launcher` は単体で動くCLI/TUIアプリとして独立させる。

## Milestone 1: リポジトリ構成を配布前提に整える

目的: 開発用ファイル、スクリプト、ビルド成果物、ドキュメントの責務を明確にする。

- `bin/` にビルド成果物を集約し、git管理外にする。
- `scripts/` に開発・起動補助スクリプトを集約する。
- `scripts/launcher.config.sh` はスクリプト専用設定として扱う。
- README に通常実行、ビルド、Warp/skhd 経由の起動方法を明記する。
- `documents/` に設計メモ、ロードマップ、セットアップ方針を残す。

完了条件:

- ルート直下にビルド済みバイナリや一時スクリプトが散らばっていない。
- `go test ./...` と `go build -o bin/launcher .` が通る。
- SKHD から呼ぶべきスクリプトパスがREADMEに明記されている。

## Milestone 2: ランチャー本体を環境依存から切り離す

目的: `launcher` 単体を、Warp/yabai/skhd がなくても使えるようにする。

- `launcher` 本体は標準のターミナルで起動できるTUIとして維持する。
- Warp/yabai/skhd 連携は `scripts/` またはサブコマンドに分離する。
- ウィンドウ制御が失敗しても、ランチャー本体の実行には影響しないようにする。
- `/Applications` に加えて `~/Applications` も探索対象にする。
- 利用履歴を `/tmp` ではなく、`os.UserConfigDir()` または `os.UserCacheDir()` 配下へ移す。
- アイコンキャッシュも `os.UserCacheDir()` 配下へ移す。

完了条件:

- `launcher` 単体実行でアプリ検索・実行・履歴保存ができる。
- yabai が入っていない環境でも起動自体は壊れない。
- ユーザーデータの保存先がmacOS標準の場所になっている。

## Milestone 3: 設定と診断コマンドを用意する

目的: ユーザーが自分の環境を把握し、問題を解決しやすくする。

- `launcher doctor` を追加し、Goビルド後の実行環境を診断する。
- `doctor` で確認する項目:
  - macOSかどうか
  - `/Applications` と `~/Applications` の読み取り可否
  - `sips` の有無
  - yabai/skhd/Warp の有無
  - 設定ディレクトリとキャッシュディレクトリの書き込み可否
- `launcher config path` で設定ファイルの場所を表示する。
- `launcher cache clear` でアイコンキャッシュを削除できるようにする。

完了条件:

- セットアップに失敗したユーザーが `launcher doctor` の出力を貼れば状況が分かる。
- キャッシュや設定の場所をユーザーが自分で探さなくてよい。

## Milestone 4: ホットキー起動をインストール可能にする

目的: SKHD を使う人にも、使わない人にも起動しやすい道を用意する。

- `launcher install-skhd` を追加し、skhd設定例を生成する。
- 既存の `~/.skhdrc` を上書きせず、追記案を表示するか確認付きで追記する。
- yabai/Warp連携はオプション扱いにする。
- yabaiがない場合は、通常のターミナル起動やmacOS Automator/Shortcuts案を提示する。
- 将来的にはLaunchAgentや小さなメニューバーアプリも検討する。

完了条件:

- skhd利用者はコマンドひとつで設定の雛形を得られる。
- skhd未利用者にも代替の起動手段がREADMEにある。

## Milestone 5: Homebrew で配布する

目的: ユーザーがGo環境なしでインストールできるようにする。

- GitHub Releases でmacOS向けビルド済みバイナリを配布する。
- `goreleaser` の導入を検討する。
- Homebrew Tap を作る。
- Formula で `bin.install "launcher"` を提供する。
- 必要に応じて補助スクリプトやサンプル設定も `pkgshare` に配置する。

理想の導入手順:

```sh
brew tap kihhi/tap
brew install launcher
launcher doctor
```

完了条件:

- Go未インストールのmacOS環境で `brew install` だけで `launcher` が使える。
- `launcher doctor` が導入後の状態を説明できる。
- READMEにHomebrewでの導入手順がある。

## Milestone 6: 初回体験を磨く

目的: インストール直後に迷わず使えるようにする。

- 初回起動時に簡単な使い方を表示する。
- `Recent` が空のときの表示を自然にする。
- アプリアイコン変換中でもUIが固まりすぎないようにする。
- 検索結果が空のときにヒントを表示する。
- `>` コマンドモードや `/` ファイルモードの実装範囲を整理する。

完了条件:

- 初めて起動したユーザーが、何を入力すればよいか分かる。
- アイコン生成やキャッシュ作成が失敗しても、ランチャーとして使える。

## Milestone 7: 配布後の品質を保つ

目的: 使う人が増えても壊れにくく、改善しやすい状態にする。

- GitHub Actions で `go test ./...` と `go build` を実行する。
- リリース作成を自動化する。
- README にトラブルシューティングを追加する。
- `documents/` に設計判断や既知の制約を残す。
- バグ報告テンプレートを用意する。

完了条件:

- mainブランチの基本品質がCIで守られる。
- リリース手順が手作業に依存しすぎない。

## 優先順位

1. `launcher` 本体を単体で安定動作させる。
2. 保存先・キャッシュ先をmacOS標準ディレクトリに移す。
3. `doctor` と設定確認コマンドを作る。
4. skhd/yabai/Warp連携をオプション化する。
5. GitHub Releases とHomebrew Tapで配布する。

## メモ

現時点ではyabai/skhd/Warp連携を急いで捨てる必要はない。lvncerの実運用では便利なので残す。ただし、配布時には「必須依存」ではなく「便利な起動方法のひとつ」に下げるのがよい。
