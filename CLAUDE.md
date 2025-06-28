# 音声通知MCP Server 要件仕様書

## 1. プロジェクト概要

### 1.1 目的
作業の完了通知や承認待ちなどのタイミングで、AIアシスタントがMCP（Model Context Protocol）を通じて自律的に音声通知を行えるようにするローカルサーバーを開発する。ユーザーからの明示的な指示がなくても、AIが適切なタイミングを判断して音声通知を活用し、より良いユーザー体験を提供する。

### 1.2 スコープ
- **対象OS**: macOS専用
- **音声エンジン**: macOS標準の`say`コマンド
- **配布方法**: Go moduleとして公開し、git cloneなしで直接実行可能
- **対応アプリケーション**:
  - Claude Desktop（Desktop Extensions対応）
  - Claude Code
  - Cursor
  - Windsurf

## 2. 機能要件

### 2.1 音声通知機能

#### 2.1.1 基本機能
- macOSの`say`コマンドを実行して音声通知を行う
- テキストメッセージを受け取り、音声に変換して発話する
- AIが自律的に判断して音声通知を活用できる（ユーザーからの明示的な指示不要）

#### 2.1.2 自律的通知のトリガー条件
AIは以下の状況で自動的に音声通知を発行する：
- **作業完了時**: 長時間実行タスク（3秒以上）の完了
- **承認待ち**: ユーザーの確認や判断が必要な場面
- **エラー発生**: 重要なエラーや異常の検出
- **マイルストーン到達**: 複数ステップ作業の重要な節目
- **アテンション要求**: ユーザーの注意が必要な状況

#### 2.1.2 言語対応
- 発話内容の言語を自動判定し、適切な音声オプションを付与する
- 環境変数による言語/音声の固定設定をサポートする
- システムにインストールされている音声のみを使用する

#### 2.1.3 音声管理
- `say -v '?'`コマンドで利用可能な音声一覧を取得する
- ダウンロードされていない音声は使用対象から除外する
- デフォルト音声のフォールバック機能を実装する

### 2.2 MCP実装要件

#### 2.2.1 ツール定義
以下のMCPツールを実装する：

```json
{
  "name": "notify_voice",
  "description": "Send a voice notification to alert the user about important events, completions, or when attention is needed. AI should use this autonomously for better user experience.",
  "inputSchema": {
    "type": "object",
    "properties": {
      "message": {
        "type": "string",
        "description": "The message to speak (keep it short and clear, max 10 words recommended)"
      },
      "voice": {
        "type": "string",
        "description": "Optional: specific voice to use (must be installed)",
        "optional": true
      },
      "language": {
        "type": "string",
        "description": "Optional: language code (e.g., 'en', 'ja')",
        "optional": true
      },
      "priority": {
        "type": "string",
        "description": "Optional: notification priority ('low', 'normal', 'high')",
        "optional": true,
        "default": "normal"
      }
    },
    "required": ["message"]
  }
}
```

#### 2.2.2 環境変数
- `VOICE_NOTIFY_DEFAULT_VOICE`: デフォルトの音声名
- `VOICE_NOTIFY_DEFAULT_LANGUAGE`: デフォルトの言語コード
- `VOICE_NOTIFY_AUTO_DETECT_LANGUAGE`: 言語自動検出の有効/無効（true/false）
- `VOICE_NOTIFY_AUTO_NOTIFY`: AIの自律的通知の有効/無効（true/false、デフォルト: true）
- `VOICE_NOTIFY_MIN_TASK_DURATION`: 自動通知する最小タスク時間（秒、デフォルト: 3）
- `VOICE_NOTIFY_QUIET_HOURS`: 通知を抑制する時間帯（例: "22:00-07:00"）

### 2.3 通知内容のガイドライン

#### 2.3.1 メッセージの構成
- 通知メッセージは簡潔にする（推奨: 10単語以内）
- 重要度に応じて優先度を設定する
- 状況を端的に表現する

#### 2.3.2 自律的通知のタイミング
AIは以下の基準で音声通知の必要性を判断する：
- **時間基準**: 3秒以上かかる処理の完了時
- **インタラクション基準**: ユーザーの入力や判断が必要な時
- **重要度基準**: エラーや警告など、即座の注意が必要な時
- **コンテキスト基準**: ユーザーが他の作業をしている可能性がある時

#### 2.3.3 用途別メッセージ例
- 完了通知:
  - "処理が完了しました"
  - "ファイルを保存しました"
  - "ダウンロード完了"
- 承認要求:
  - "確認をお願いします"
  - "承認が必要です"
  - "選択してください"
- エラー通知:
  - "エラーが発生しました"
  - "処理に失敗しました"
  - "接続できません"
- 進捗通知:
  - "処理を開始します"
  - "半分完了しました"
  - "もうすぐ完了します"

## 3. 技術仕様

### 3.1 開発環境
- **言語**: Go
- **MCPライブラリ**: `golang.org/x/tools/internal/mcp`
- **最小Goバージョン**: 1.21以上
- **モジュール名**: `github.com/kyong0612/voice-notify-mcp`（go.modで定義）

### 3.2 プロジェクト構造
```
voice-notify-mcp/
├── main.go              # エントリーポイント
├── server.go            # MCPサーバー実装
├── voice.go             # 音声処理ロジック
├── language.go          # 言語検出ロジック
├── notification.go      # 自律的通知判定ロジック
├── go.mod              # モジュール定義（module github.com/kyong0612/voice-notify-mcp）
├── go.sum
├── README.md           # 使用方法とインストール手順
├── CLAUDE.md           # 本仕様書
└── dxt.json            # Desktop Extensions設定
```

### 3.3 実装詳細

#### 3.3.1 go.modファイル例
```go
module github.com/kyong0612/voice-notify-mcp

go 1.21

require (
    golang.org/x/tools v0.x.x
)
```

#### 3.3.2 音声選択ロジック
1. 明示的に音声が指定された場合、その音声を使用
2. 言語が指定された場合、その言語に対応する音声を選択
3. 自動検出が有効な場合、メッセージ内容から言語を推定
4. 環境変数でデフォルト音声が設定されている場合、それを使用
5. いずれも該当しない場合、システムデフォルト音声を使用

#### 3.3.3 自律的通知の実装
- AIは内部でタスクの実行時間を計測し、閾値を超えた場合に自動通知
- 通知の優先度に基づいて、音声の速度や音量を調整可能
- 静音時間帯（Quiet Hours）のチェック機能
- 連続通知の抑制（同じ種類の通知は一定時間内に1回まで）

#### 3.3.4 エラーハンドリング
- 指定された音声が利用不可能な場合、デフォルト音声にフォールバック
- `say`コマンドの実行に失敗した場合、エラーレスポンスを返す
- 音声一覧の取得に失敗した場合、キャッシュされた情報を使用
- 静音時間帯の場合、通知をスキップし、ログに記録

## 4. インストールと設定

Go 1.17以降の機能により、リポジトリをクローンすることなく、モジュール名を指定して直接実行できます。これにより、インストールプロセスが大幅に簡素化されます。

### 4.1 起動方法

#### 前提条件
- Go 1.17以降がインストールされていること
- インターネット接続（初回実行時のモジュールダウンロード用）

#### 4.1.1 リモート実行（推奨）
Git cloneせずに直接実行：
```bash
go run github.com/kyong0612/voice-notify-mcp@latest
```

#### 4.1.2 ローカル実行
リポジトリをクローンした場合：
```bash
go run main.go
```

### 4.2 Claude Desktop設定（Desktop Extensions）

#### 4.2.1 dxt.json
```json
{
  "name": "voice-notify",
  "version": "1.0.0",
  "description": "Voice notification MCP server for macOS with autonomous AI notifications",
  "mcpServers": {
    "voice-notify": {
      "command": "go",
      "args": ["run", "github.com/kyong0612/voice-notify-mcp@latest"],
      "env": {
        "VOICE_NOTIFY_DEFAULT_LANGUAGE": "en",
        "VOICE_NOTIFY_AUTO_DETECT_LANGUAGE": "true",
        "VOICE_NOTIFY_AUTO_NOTIFY": "true",
        "VOICE_NOTIFY_MIN_TASK_DURATION": "3"
      }
    }
  }
}
```

注: ローカルインストールの場合は、以下の設定も可能：
```json
{
  "mcpServers": {
    "voice-notify": {
      "command": "go",
      "args": ["run", "main.go"],
      "cwd": "${extensionPath}",
      "env": { ... }
    }
  }
}
```

### 4.3 その他のアプリケーション設定

#### 4.3.1 Claude Code / Cursor / Windsurf
MCPサーバー設定に以下を追加：

**リモート実行（推奨）**：
```json
{
  "voice-notify": {
    "command": "go",
    "args": ["run", "github.com/kyong0612/voice-notify-mcp@latest"],
    "env": {
      "VOICE_NOTIFY_DEFAULT_VOICE": "Samantha",
      "VOICE_NOTIFY_DEFAULT_LANGUAGE": "en",
      "VOICE_NOTIFY_AUTO_DETECT_LANGUAGE": "true",
      "VOICE_NOTIFY_MIN_TASK_DURATION": "3"
    }
  }
}
```

**ローカル実行**：
```json
{
  "voice-notify": {
    "command": "go",
    "args": ["run", "/path/to/voice-notify-mcp/main.go"],
    "env": { ... }
  }
}
```

## 5. 使用例

### 5.1 自律的な通知の例

#### 5.1.1 長時間処理の完了
```
ユーザー: 「このデータセットを分析して」
AI: 「分析を開始します...」
（5分後）
AI: 「分析が完了しました。結果をまとめています。」
→ MCPツール呼び出し: notify_voice(message="分析が完了しました")
```

#### 5.1.2 承認待ち
```
AI: 「提案書の作成が完了しました。内容を確認していただけますか？」
→ MCPツール呼び出し: notify_voice(message="確認をお願いします")
```

#### 5.1.3 エラー通知
```
AI: 「ファイルの処理中にエラーが発生しました。」
→ MCPツール呼び出し: notify_voice(message="エラーが発生しました")
```

### 5.2 明示的な指示による使用
```
ユーザー: 「完了したら音声で知らせて」
AI: 「承知しました。完了時に音声でお知らせします。」
→ MCPツール呼び出し: notify_voice(message="作業が完了しました")
```

### 5.3 カスタム設定
```
ユーザー: 「日本語の音声で通知して」
AI: 「日本語音声で通知するよう設定しました。」
→ MCPツール呼び出し: notify_voice(message="設定が完了しました", language="ja")
```

## 6. セキュリティ考慮事項

- ローカル実行のみに制限（外部ネットワークアクセスなし）
- 実行可能なコマンドは`say`コマンドのみに制限
- 入力文字列のサニタイゼーション実装
- コマンドインジェクション対策の実装
- リモート実行時は、信頼できるリポジトリからのみ実行すること
- バージョン指定（@v1.0.0など）により、意図しない更新を防止可能

## 5. AIアシスタント向けのガイドライン

このセクションは、AIアシスタントがこの音声通知MCPサーバーを効果的に活用するためのガイドラインです。

### 5.1 基本的な心構え
あなたは、ユーザーの生産性を最大化するために、音声通知機能を積極的に活用することが期待されています。ユーザーがPCから目を離している状況でも、重要な情報をタイムリーに伝えることで、作業の効率を大幅に向上させることができます。

### 5.2 通知を推奨する具体的なシナリオ
以下の状況では、積極的に`voice-notify/notify_voice`ツールの使用を検討してください。

- **時間のかかる処理の完了**: ビルド、テスト、デプロイ、ファイルのダウンロード、データ分析など、完了までに数秒以上かかるタスクが終了した際に通知します。
  - **例**: `mcp/voice-notify/notify_voice 'ビルドが完了しました。'`

- **ユーザーの操作が必要な場合**: プロセスの途中でユーザーの承認、入力、または選択が必要になった際に、操作を促すために通知します。
  - **例**: `mcp/voice-notify/notify_voice '承認が必要です。続行しますか？'`

- **エラーや問題の発生**: 処理が失敗したり、警告が発生したりした場合に、即座にユーザーの注意を引くために通知します。
  - **例**: `mcp/voice-notify/notify_voice 'デプロイに失敗しました。詳細はログを確認してください。' 'ja' 'high'`

- **重要な節目**: 複数のステップからなる作業において、重要な中間地点に到達したことを知らせます。
  - **例**: `mcp/voice-notify/notify_voice 'ステージング環境へのデプロイが完了しました。最終確認をお願いします。'`

### 5.3 効果的なメッセージの作り方
- **簡潔に**: ユーザーがすぐに状況を理解できるよう、メッセージは短く、要点をまとめてください。（推奨：10単語以内）
- **明確に**: 何が起きたのか、次に何をすべきかが分かるように具体的に表現してください。
- **言語の指定**: 可能であれば、メッセージの言語に合わせて`language`パラメータ（例: `'ja'`, `'en'`）を指定すると、より自然な音声で通知できます。
- **優先度の活用**: 緊急性が高いエラー通知などの場合は、`priority`パラメータを`'high'`に設定してください。

### 5.4 ユーザーへの提案
時間のかかるタスクを開始する際には、「完了したら音声で通知しましょうか？」のように、事前にユーザーへ通知の利用を提案することも効果的です。