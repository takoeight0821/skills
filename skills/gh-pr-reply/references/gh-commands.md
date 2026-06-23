# Reference: PR レビューコメントの返信・resolve コマンド集

gh 2.94.0 / GitHub API（GraphQL は live introspection、REST は docs.github.com）で検証済み。

## 概念・ハマりどころ

- **3 種のコメント**: issue comment（PR 全体への会話、`/issues/{n}/comments`）/ review comment
  （diff の行に紐づくインライン、`/pulls/{n}/comments`）/ review thread（review comment + 返信の
  解決可能な束、**GraphQL のみ**、`isResolved` はここ）。
- **node id vs databaseId**: GraphQL は不透明な node id（thread は `PRRT_...`、comment は `PRRC_...`）。
  REST は数値 `databaseId`（= REST の `id`）。resolve mutation は **thread の node id**、
  REST の返信エンドポイントは **数値 comment_id** を要求する。
- **`in_reply_to_id`**: REST のコメント一覧で、返信コメントは `in_reply_to_id`（整数）で起点コメントの
  id を指す。これでスレッドを復元できる。
- **この環境のフック**: `gh-api-readonly` が `--method GET` 無しの `gh api` を全て deny する。
  返信・resolve（POST/GraphQL mutation）は `gh-pr-reply` / `gh-pr-resolve` ラッパー経由で行う。
  ラッパー名は `gh api` 正規表現にマッチせず、内部の `gh api` は別プロセスなのでフック対象外。

## 1. 対象 PR / repo の特定

```bash
gh repo view --json owner,name --jq '.owner.login + "/" + .name'
gh pr view --json number,headRefName,url
PR=$(gh pr view --json number --jq .number)
```

## 2. レビューコメント一覧（REST GET、フック許可済み）

`GET /repos/{owner}/{repo}/pulls/{pull_number}/comments`

```bash
gh api --method GET "repos/{owner}/{repo}/pulls/123/comments" --paginate \
  --jq '.[] | {id, in_reply_to_id, path, line, user: .user.login, body}'
```

## 3. 返信を投稿する（gh-pr-reply ラッパー）

内部で `POST /repos/{owner}/{repo}/pulls/{pull_number}/comments/{comment_id}/replies` を叩く。
`comment_id` は数値 databaseId。

```bash
gh-pr-reply <owner/repo> <pull_number> <comment_id> "<body>"
# 例
gh-pr-reply octocat/hello-world 42 123456 "Thanks, fixed in abc1234."
```

## 4. スレッドを resolve / unresolve する（gh-pr-resolve ラッパー）

`comment_id` から内部 GraphQL クエリで thread node id を解決し、
`resolveReviewThread` / `unresolveReviewThread` mutation を実行する。

```bash
gh-pr-resolve <owner/repo> <pull_number> <comment_id>            # resolve（第 4 引数省略時）
gh-pr-resolve <owner/repo> <pull_number> <comment_id> unresolve  # unresolve
# 例
gh-pr-resolve octocat/hello-world 42 123456
```

## 参考: ラッパーが内部で使う生コマンド（フックがあるため直接実行は不可）

返信（REST POST）:

```bash
gh api --method POST \
  -H "Accept: application/vnd.github+json" \
  "/repos/{owner}/{repo}/pulls/{n}/comments/{comment_id}/replies" \
  -f body="..."
```

レビュースレッド + 解決状態の取得（GraphQL query）:

```bash
gh api graphql -f owner='OWNER' -f name='REPO' -F pr=123 -f query='
query($owner:String!,$name:String!,$pr:Int!){
  repository(owner:$owner,name:$name){
    pullRequest(number:$pr){
      reviewThreads(first:100){
        nodes{
          id isResolved isOutdated path line
          comments(first:100){ nodes{ id databaseId author{login} body path line } }
        }
      }
    }
  }
}'
```

resolve / unresolve（GraphQL mutation、`threadId` は thread の node id）:

```bash
gh api graphql -f threadId='PRRT_...' -f query='
mutation($threadId:ID!){ resolveReviewThread(input:{threadId:$threadId}){ thread{ id isResolved } } }'

gh api graphql -f threadId='PRRT_...' -f query='
mutation($threadId:ID!){ unresolveReviewThread(input:{threadId:$threadId}){ thread{ id isResolved } } }'
```
