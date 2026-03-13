---
name: ask-codex
description: >
  Get a second opinion from OpenAI Codex CLI for code review or discussion.
  Use when: "ask codex", "codex review", "second opinion", "check with codex",
  「Codexに聞いて」「codexにレビューしてもらって」「codexでチェック」
  「codexに見てもらって」「codexに相談」「別の意見がほしい」
  「セカンドオピニオン」「他のAIに聞いて」「別のモデルに確認して」
  「他の視点からチェックして」など。
  Do NOT trigger on plain "レビューして" without Codex mention or second-opinion context.
---

# Ask Codex

OpenAI Codex CLIを呼び出してセカンドオピニオンを得る。

## 前提条件の確認

まず `codex` コマンドが利用可能か確認する:

```bash
command -v codex
```

コマンドが見つからない場合はユーザーに通知して終了する:
- インストール: `npm install -g @openai/codex`
- 認証: `OPENAI_API_KEY` の設定が必要

## 用途の判別

ユーザーの依頼を以下の3パターンに分類する:

1. **コードレビュー**: コード変更のレビューを依頼 → `codex exec review` を使用
2. **設計相談・意見交換**: 設計方針や実装方法について相談 → `codex exec` を使用
3. **自由質問**: その他の質問をCodexに転送 → `codex exec` を使用

## パターン1: コードレビュー

`codex exec review` はリポジトリのコード変更を自動的にレビューする専用サブコマンド。

### 対象の指定

ユーザーの依頼に応じてオプションを選択する:

- **未コミットの変更**: `--uncommitted` — staged/unstaged/untracked すべてを対象
- **ブランチ差分**: `--base <branch>` — 指定ブランチとの差分をレビュー
- **特定コミット**: `--commit <sha>` — 指定コミットの変更をレビュー
- **指定なし**: デフォルトでHEADコミットの変更をレビュー

### 実行

```bash
codex exec review --uncommitted -o /tmp/codex-review.txt --ephemeral -s read-only
```

カスタム指示がある場合はプロンプトを追加:

```bash
codex exec review --uncommitted "セキュリティの観点で重点的にレビューしてください" -o /tmp/codex-review.txt --ephemeral -s read-only
```

`-o` フラグで最終メッセージをファイルに保存し、Readツールで結果を読み取る。

## パターン2: 設計相談・自由質問

`codex exec` で任意のプロンプトをCodexに送る。

### プロンプト構築

コードベースや対象ファイルの内容を含めてプロンプトを構築する。プロンプトの言語はユーザーの言語に合わせる。

長いプロンプトはstdin経由で渡す:

```bash
echo "<prompt>" | codex exec - -o /tmp/codex-response.txt --ephemeral -s read-only
```

短いプロンプトは引数で渡す:

```bash
codex exec "<prompt>" -o /tmp/codex-response.txt --ephemeral -s read-only
```

### 実行時の注意

- `-s read-only` でファイル書き込みを禁止する（必須）
- `--ephemeral` でセッションファイルを残さない
- `-o /tmp/codex-response.txt` で結果をファイル出力し、Readツールで読む
- タイムアウトは300秒（5分）に設定する
- `--json` を使うとJSONL形式でイベントストリームが出力される（通常は `-o` で十分）

## エラーハンドリング

- **コマンド失敗**: エラーメッセージをユーザーに提示
- **タイムアウト**: プロンプトを短くして再試行するか、ユーザーに報告
- **認証エラー**: `OPENAI_API_KEY` の設定を確認するよう案内

## 結果の提示

`-o` で保存したファイルをReadツールで読み、内容を整理してユーザーに提示する。

### 提示フォーマット

```markdown
## Codexの意見

### 指摘事項
- 指摘1: ...
- 指摘2: ...

### 改善提案
- 提案1: ...
- 提案2: ...

### 補足
Codexのコメントや追加の文脈
```

### Claudeとの比較（Claudeが既にレビュー済みの場合のみ）

Claudeが同じコードについて既に見解を持っている場合、両方の視点を比較して提示する:

```markdown
## 意見の比較

| 観点 | Claude | Codex |
|------|--------|-------|
| ... | ... | ... |
```

## フォローアップ

結果提示後、ユーザーに次のアクションを確認する:

- Codexの指摘を反映するか
- 特定の指摘について深掘りするか
- 別の観点でもう一度Codexに聞くか
