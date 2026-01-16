# Claude Skills

Claude Codeのagent skillを管理する共用リポジトリです。`skills` CLI（または `mise`）を使って各プロジェクトの`.claude/skills`やグローバルの`~/.claude/skills`に展開して利用します。

## 必要条件

- Git
- [Multipass](https://multipass.run/) - VM管理の場合
- [Docker](https://www.docker.com/) - コンテナ管理の場合
- Go 1.22+ - `skills` CLIのインストールに必要
- [mise](https://mise.jdx.dev/) - タスクランナーとして使用する場合（任意）

## skills CLI（推奨）

VM管理、コンテナ管理、スキル同期を統合したGo製のCLIツール `skills` の使用を推奨します。

### インストール

```bash
# Goを使ってインストール
go install github.com/takoeight0821/skills/jig/cmd/skills@latest
```

### クイックスタート

```bash
# 初期設定
skills config init

# VMの作成と起動
skills vm launch

# スキルの同期（グローバル）
skills sync --global --apply

# Claude Codeの実行
skills vm claude
```

詳細は [jig/README.md](jig/README.md) を参照してください。

## 使い方

### スキルの同期

`skills sync` コマンドを使用して、本リポジトリのスキルをローカル環境に同期します。

#### グローバルに同期（~/.claude/skills）

```bash
# プレビュー（dry-run）
skills sync --global --dry-run

# 実際に同期
skills sync --global --apply
```

#### プロジェクトに同期（.claude/skills）

```bash
# プレビュー（dry-run）
skills sync --project --dry-run

# 実際に同期
skills sync --project --apply
```

### 開発環境の管理

#### Multipass VM

```bash
skills vm launch   # 作成と起動
skills vm ssh      # SSH接続
skills vm stop     # 停止
```

#### Docker Container

```bash
skills docker launch   # ビルドと起動
skills docker ssh      # シェル接続
skills docker stop     # 停止
```

## ディレクトリ構造

```
/path/to/skills/                # クローンしたリポジトリ
├── skills/                     # 同期対象のスキル
│   └── my-skill/
│       ├── SKILL.md
│       └── ...
├── jig/                        # skills CLIのソースコード
│   ├── cmd/skills/             # エントリポイント
│   └── ...
├── docker/                     # Docker関連ファイル
├── multipass/                  # Multipass関連ファイル
├── install.sh                  # 旧インストールスクリプト（mise用）
└── README.md
```

## ライセンス

MIT License