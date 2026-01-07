# Claude Skills

Claude Codeのagent skillを管理する共用リポジトリです。miseを使って各プロジェクトの`.claude/skills`やグローバルの`~/.claude/skills`に展開して利用します。

## 必要条件

- [mise](https://mise.jdx.dev/) がインストールされていること
- Git

## インストール

```bash
# リポジトリを任意の場所にクローン
git clone https://github.com/takoeight0821/skills.git
cd skills

# インストールスクリプトを実行
./install.sh
```

インストールスクリプトは以下を行います：

1. クローンしたディレクトリを検出
2. miseのタスク定義を`~/.config/mise/config.toml`に追加（クローン先を参照）

**注意**: このリポジトリをクローンした場所がスキルの共有元として使用されます。リポジトリを移動した場合は再度`./install.sh`を実行してください。

## 使い方

### 共有スキルの更新

```bash
mise run update-shared-skills
```

リモートリポジトリから最新のスキルを取得します。

### グローバルに同期（~/.claude/skills）

```bash
# プレビュー（dry-run、実際には同期しない）
mise run sync-skills-global

# 実際に同期
mise run sync-skills-global-apply

# オプションを追加して同期
mise run sync-skills-global-apply -- --force --prune
```

### プロジェクトに同期（.claude/skills）

```bash
# プレビュー（dry-run、実際には同期しない）
mise run sync-skills-project

# 実際に同期
mise run sync-skills-project-apply

# オプションを追加して同期
mise run sync-skills-project-apply -- --force --prune --exclude
```

## 動作詳細

### Dry-runモード（デフォルト）

デフォルトではdry-runモードで動作し、実際にはファイルの変更を行いません。何が行われるかをプレビューできます。

```
=== DRY-RUN MODE (use --apply to actually sync) ===

[dry-run] Would add: skills-setup/

Summary:
  Would add:    1
  Would update: 0

Run with --apply to actually perform these changes.
```

実際に同期するには`*-apply`タスクを使用するか、設定ファイルで`apply=true`を指定します。

### 上書き確認

既存のスキルディレクトリと共有リポジトリのディレクトリに差分がある場合、警告を表示して上書きの確認を求めます。

```
Warning: skills-setup/ differs from shared version
Overwrite? [y/N]
```

`--force`オプションを使うと確認なしで上書きします。

### 同期の記録

同期したスキルは`.skills-manifest`ファイルに記録されます：

- グローバル: `~/.claude/.skills-manifest`
- プロジェクト: `.claude/.skills-manifest`

### 削除されたスキルの処理

`--prune`オプションを使うと、共有リポジトリから削除されたスキルがローカルからも削除されます。マニフェストに記録されているディレクトリのみが削除対象となるため、プロジェクト固有のスキルは影響を受けません。

### gitの追跡から除外

`--exclude`オプションを使うと、同期したスキルディレクトリが`.git/info/exclude`に追加され、gitの追跡から除外されます。これにより、共有スキルがプロジェクトのgit履歴に含まれなくなります。

## コマンドラインオプション

| オプション | 説明 |
|-----------|------|
| `--apply` | 実際に同期を実行（デフォルトはdry-run） |
| `--force` | 確認なしで上書き |
| `--prune` | 削除されたスキルを削除 |
| `--exclude` | .git/excludeに追加（project modeのみ） |

## ディレクトリ構造

```
/path/to/skills/                # クローンしたリポジトリ
├── .claude/
│   └── skills/
│       └── skills-setup/       # 設定用スキル（同期対象外）
├── skills/                     # 同期対象のスキル
│   └── my-skill/
│       ├── SKILL.md
│       └── ...
├── bin/
│   └── sync-skills.sh          # 同期スクリプト
├── install.sh
└── README.md

~/.config/skills/
└── config                      # グローバル設定ファイル

~/.claude/                      # グローバルスキル
├── skills/
│   └── my-skill/               # ← sync-skills-globalで同期
└── .skills-manifest            # 同期記録

your-project/                   # プロジェクト
├── .claude/
│   ├── skills/
│   │   ├── my-skill/           # ← sync-skills-projectで同期
│   │   └── my-local-skill/     # プロジェクト固有（影響なし）
│   ├── .skills-manifest        # 同期記録
│   └── .skills.conf            # プロジェクト設定ファイル
└── .git/
    └── info/
        └── exclude             # --excludeオプションで追加
```

## 設定ファイル

コマンドラインオプションの代わりに、設定ファイルでデフォルトの動作を指定できます。

### 設定ファイルの場所

| ファイル | 説明 |
|---------|------|
| `~/.config/skills/config` | グローバル設定 |
| `.claude/.skills.conf` | プロジェクト設定 |

### 優先順位

1. コマンドラインフラグ（最優先）
2. プロジェクト設定ファイル
3. グローバル設定ファイル
4. デフォルト値

### 設定ファイルの形式

```ini
# コメント行
apply=true
force=true
prune=false
exclude=true
```

### 設定項目

| 項目 | デフォルト | 説明 |
|------|-----------|------|
| `apply` | `false` | 実際に同期を実行（falseはdry-run） |
| `force` | `false` | 確認なしで上書き |
| `prune` | `false` | 削除されたスキルを削除 |
| `exclude` | `false` | .git/excludeに追加 |

### 使用例

グローバル設定で常に実行モードにする:

```bash
mkdir -p ~/.config/skills
echo "apply=true" > ~/.config/skills/config
```

プロジェクト設定で`apply`、`force`、`prune`を有効にする:

```bash
echo -e "apply=true\nforce=true\nprune=true" > .claude/.skills.conf
```

## 環境変数

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `SKILLS_SHARED_DIR` | `<repo>/skills` | 共有スキルの場所 |

## 自動同期（オプション）

プロジェクトディレクトリに入った時に自動でプレビューを表示したい場合は、`~/.config/mise/config.toml`に以下を追加してください：

```toml
[hooks]
enter = "mise run sync-skills-project 2>/dev/null || true"
```

自動で実際に同期したい場合は、設定ファイルで`apply=true`を指定するか、タスクを変更してください。

## スキルの追加

共有スキルを追加するには、`skills/`ディレクトリにスキルディレクトリを作成してコミットしてください。

```bash
cd /path/to/skills
mkdir -p skills/my-new-skill
vim skills/my-new-skill/SKILL.md
git add skills/my-new-skill
git commit -m "Add my-new-skill"
git push
```

## アンインストール

```bash
# mise設定から削除
# ~/.config/mise/config.toml から # BEGIN claude-skills 〜 # END claude-skills を削除

# リポジトリを削除（クローンした場所）
rm -rf /path/to/skills
```

## ライセンス

MIT License
