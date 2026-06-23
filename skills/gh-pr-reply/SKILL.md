---
name: gh-pr-reply
description: GitHub PR のレビューコメント（インラインのコードレビューコメント）に gh CLI で返信し、対応済みのレビュースレッドを resolve / unresolve する。返信は gh-pr-reply ラッパー、resolve は gh-pr-resolve ラッパー経由で行い、コメント一覧の取得は REST GET で行う。TRIGGER when：「PR のレビューコメントに返信して」「レビュー指摘に返信」「review コメントへ返信して」「対応したスレッドを resolve して」「レビュースレッドを解決して」「PR review reply」などと言われた時。SKIP when：PR 本体への一般コメント・会話コメントの投稿（→ `gh pr comment`、issue comment）、コードレビューの実施そのもの（指摘を書く側）、PR の作成・マージ・チェック確認（→ 通常の gh コマンド）。
---

# gh-pr-reply — PR レビューコメントへの返信と resolve

GitHub PR の**インラインレビューコメント**に返信し、対応が済んだレビュースレッドを resolve する。

## 前提・概念（最初に読む）

GitHub のコメントには 3 種類あり、混同しないこと。このスキルが扱うのは **review comment / review thread** のみ。

| 種類 | 何か | 操作 |
|---|---|---|
| issue comment | PR 全体への会話コメント | このスキルの対象外（`gh pr comment`） |
| review comment | diff の特定ファイル・行に紐づくインラインコメント | **返信対象**。REST `/pulls/{n}/comments` |
| review thread | review comment とその返信をまとめた解決可能な単位 | **resolve 対象**。GraphQL のみに存在し、`isResolved` はここに乗る |

実行手段の鉄則（この環境特有）:

- **生の `gh api` の POST は `gh-api-readonly` フックに deny される。** 返信・resolve を直接 `gh api` で
  叩いてはいけない。必ず下記ラッパーを使う。
  - 返信投稿 → `gh-pr-reply`（PATH 上のコマンド）
  - resolve / unresolve → `gh-pr-resolve`（PATH 上のコマンド）
- 読み取り（コメント一覧）は `gh api --method GET` なら許可される。
- ラッパーに渡す `<comment_id>` は **REST の数値コメント id**（GraphQL の `PRRC_...` node id ではない）。

## ワークフロー

### Step 1: 対象 PR と repo を特定する

```bash
gh repo view --json owner,name --jq '.owner.login + "/" + .name'   # -> owner/repo
gh pr view --json number,headRefName,url                          # -> 現ブランチの PR
```

PR が自動特定できない場合は、対象の PR 番号をユーザーに確認する。

### Step 2: レビューコメントを取得する（REST GET、フック許可済み）

```bash
gh api --method GET "repos/{owner}/{repo}/pulls/{pull_number}/comments" --paginate \
  --jq '.[] | {id, in_reply_to_id, path, line, user: .user.login, body}'
```

- `id` … 返信・resolve に渡す数値コメント id
- `in_reply_to_id` … 非 null なら既存スレッドへの返信コメント。null がスレッドの起点
- これでスレッド構造（どのコメントにどの返信がぶら下がるか）を復元する

> 注意: スレッドが resolve 済みかどうか（`isResolved`）は REST GET では取得できない（GraphQL 読みは
> フックで不可）。resolve 済みかの判断はユーザーの指示、または `gh-pr-resolve` 実行結果の出力に依存する。

### Step 3: 返信内容を提示して確認する

返信は外部への書き込み（PR に公開され、通知が飛ぶ）。**投稿前に必ず**「どのコメント（path:line と要旨）に
何を返信するか」を要約してユーザーに提示し、了承を得る。

### Step 4: 返信を投稿する

```bash
gh-pr-reply <owner/repo> <pull_number> <comment_id> "<返信本文>"
# 例: gh-pr-reply octocat/hello-world 42 123456 "ご指摘ありがとうございます。abc1234 で修正しました。"
```

複数コメントに返信する場合は 1 件ずつ実行する。

### Step 5: 対応済みスレッドを resolve する（任意）

対応が完了したスレッドを解決済みにする。`gh-pr-resolve` は数値 comment_id から内部で
スレッドを特定して resolve するので、Step 2 で得た comment id をそのまま渡せばよい。

```bash
gh-pr-resolve <owner/repo> <pull_number> <comment_id>            # resolve（既定）
gh-pr-resolve <owner/repo> <pull_number> <comment_id> unresolve  # 取り消し
```

## 補足

- 詳細なコマンド・GraphQL の内訳・node id と databaseId の使い分けは
  [`references/gh-commands.md`](references/gh-commands.md) を参照。
- `gh` が未認証なら `gh auth login` を案内する。
- comment_id が存在しない・別 PR のものだと `gh-pr-resolve` は「該当スレッドが見つからない」で
  エラー終了する。Step 2 の出力から正しい id を選ぶこと。
