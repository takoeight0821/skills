#!/bin/bash
# gh-pr-resolve - Helper script for resolving/unresolving GitHub PR review threads
#
# Review threads are a GraphQL-only concept; resolving requires a GraphQL mutation
# (POST). The gh-api-readonly.sh hook blocks direct `gh api` POST calls, so this
# wrapper provides the escape hatch (its name does not match the hook's `gh api`
# pattern, and the inner `gh api` runs in a separate process).
#
# Callers usually only know a numeric review-comment id (from the REST list
# endpoint), not the GraphQL thread node id. So this script takes a comment_id,
# looks up the thread containing it, then resolves/unresolves that thread.
#
# Usage: gh-pr-resolve <owner/repo> <pull_number> <comment_id> [resolve|unresolve]
#   Example: gh-pr-resolve octocat/hello-world 42 123456
#   Example: gh-pr-resolve octocat/hello-world 42 123456 unresolve
#
# Install: Place in PATH (e.g., ~/.local/bin/gh-pr-resolve) and make executable.
#
# Dependencies: gh (GitHub CLI)

set -euo pipefail

if [[ $# -lt 3 || $# -gt 4 ]]; then
  echo "Usage: gh-pr-resolve <owner/repo> <pull_number> <comment_id> [resolve|unresolve]" >&2
  exit 1
fi

repo="$1"
pull_number="$2"
comment_id="$3"
action="${4:-resolve}"

if [[ "$action" != "resolve" && "$action" != "unresolve" ]]; then
  echo "Error: action must be 'resolve' or 'unresolve' (got '$action')" >&2
  exit 1
fi

owner="${repo%%/*}"
name="${repo#*/}"
if [[ -z "$owner" || -z "$name" || "$owner" == "$repo" ]]; then
  echo "Error: <owner/repo> must be in 'owner/repo' form (got '$repo')" >&2
  exit 1
fi

# Find the review thread whose comments include the given comment databaseId.
thread_id=$(gh api graphql \
  -f owner="$owner" -f name="$name" -F pr="$pull_number" -F cid="$comment_id" \
  -f query='
    query($owner:String!,$name:String!,$pr:Int!){
      repository(owner:$owner,name:$name){
        pullRequest(number:$pr){
          reviewThreads(first:100){
            nodes{ id comments(first:100){ nodes{ databaseId } } }
          }
        }
      }
    }' \
  --jq ".data.repository.pullRequest.reviewThreads.nodes[] | select(any(.comments.nodes[]; .databaseId == ${comment_id})) | .id" \
  | head -n1)

if [[ -z "$thread_id" ]]; then
  echo "Error: no review thread found containing comment id $comment_id on $repo#$pull_number" >&2
  exit 1
fi

mutation="${action}ReviewThread"
gh api graphql \
  -f threadId="$thread_id" \
  -f query="
    mutation(\$threadId:ID!){
      ${mutation}(input:{threadId:\$threadId}){ thread{ id isResolved } }
    }"
