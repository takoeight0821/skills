# macOS × Docker Compose × SSH Agent Forwarding で安全に署名コミットする方法

## 概要

- **目的**: Docker コンテナ内から Git の署名付きコミットを安全に実行する
- **アプローチ**: SSH Agent Forwarding を使用し、秘密鍵をコンテナに持ち込まない
- **対応プラットフォーム**: macOS (Docker Desktop / Rancher Desktop)、Linux
- **署名方式**: SSH キーによる Git コミット署名 (GitHub Verified バッジ対応)

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│  macOS ホスト                                                │
│  ├── ~/.ssh/id_ed25519 (秘密鍵 - ここにのみ存在)            │
│  ├── ssh-agent (鍵を保持)                                   │
│  └── SSH_AUTH_SOCK → Docker コンテナへ転送                  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼ ソケット転送
┌─────────────────────────────────────────────────────────────┐
│  Docker コンテナ                                             │
│  ├── 秘密鍵の実体なし (セキュア)                            │
│  ├── ~/.ssh/id_ed25519.pub のみマウント (読み取り専用)       │
│  └── git commit -S → macOS の ssh-agent で署名              │
└─────────────────────────────────────────────────────────────┘
```

## セキュリティの要点

| 項目 | 説明 |
|------|------|
| **秘密鍵の保護** | 秘密鍵は macOS 上にのみ存在し、コンテナには一切コピーしない |
| **ソケット転送** | SSH_AUTH_SOCK (Unix ソケット) のみを共有 |
| **公開鍵のみマウント** | 署名検証用に公開鍵のみを読み取り専用でマウント |
| **最小権限の原則** | コンテナは `no-new-privileges` と `CAP_DROP ALL` で実行 |

## プラットフォーム別の設定

### Docker Desktop (macOS)

Docker Desktop は `/run/host-services/ssh-auth.sock` という特別なパスを提供し、ホストの SSH Agent にアクセスできる。

```yaml
# docker-compose.macos.yml
services:
  coding-agent:
    volumes:
      # Docker Desktop の magic path
      - /run/host-services/ssh-auth.sock:/tmp/ssh-agent.sock
      # 公開鍵のみ (読み取り専用)
      - ${HOME}/.ssh/id_ed25519.pub:/home/agent/.ssh-host/id_ed25519.pub:ro
      - ${HOME}/.ssh/known_hosts:/home/agent/.ssh-host/known_hosts:ro
    environment:
      - SSH_AUTH_SOCK=/tmp/ssh-agent.sock
```

### Rancher Desktop (macOS)

Rancher Desktop は Lima VM を使用するため、macOS の Unix ソケットを直接マウントできない。SSH Agent Forwarding は利用不可。

```yaml
# docker-compose.rancher.yml
services:
  coding-agent:
    volumes:
      # 秘密鍵を直接マウント (セキュリティ上は非推奨だが代替手段なし)
      - ${HOME}/.ssh/id_ed25519:/home/agent/.ssh-host/id_ed25519:ro
      - ${HOME}/.ssh/id_ed25519.pub:/home/agent/.ssh-host/id_ed25519.pub:ro
      - ${HOME}/.ssh/known_hosts:/home/agent/.ssh-host/known_hosts:ro
    environment:
      # ソケットベース認証を無効化
      - SSH_AUTH_SOCK=
```

**注意**: Rancher Desktop では秘密鍵をコンテナ内にマウントする必要があるため、セキュリティ上 Docker Desktop の使用を推奨する。

### Linux

Linux では SSH_AUTH_SOCK を直接マウント可能。

```yaml
# docker-compose.linux.yml
services:
  coding-agent:
    volumes:
      - ${SSH_AUTH_SOCK:-/tmp/ssh-agent.sock}:/tmp/ssh-agent.sock
      - ${HOME}/.ssh/id_ed25519.pub:/home/agent/.ssh-host/id_ed25519.pub:ro
      - ${HOME}/.ssh/known_hosts:/home/agent/.ssh-host/known_hosts:ro
    environment:
      - SSH_AUTH_SOCK=/tmp/ssh-agent.sock
```

## Git SSH 署名の設定

### 前提条件

- Git 2.34 以降 ([SSH 署名のサポート](https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key))
- SSH キーペア (ed25519 推奨)
- GitHub に署名用 SSH キーを登録済み

### ホスト側の準備 (macOS)

```bash
# 1. SSH agent の起動と鍵の登録
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519

# 2. 鍵が登録されていることを確認
ssh-add -l
```

### コンテナ内の Git 設定

コンテナ起動時に以下の設定を自動適用する (entrypoint.sh):

```bash
# SSH 署名形式を使用
git config --global gpg.format ssh

# 署名キーを指定 (公開鍵のパス)
git config --global user.signingkey "/home/agent/.ssh-host/id_ed25519.pub"

# 全てのコミットに署名
git config --global commit.gpgsign true

# ユーザー情報
git config --global user.name "Your Name"
git config --global user.email "you@example.com"

# allowed_signers ファイルの設定 (ローカル検証用)
mkdir -p ~/.ssh
echo "you@example.com namespaces=\"git\" $(cat ~/.ssh-host/id_ed25519.pub)" > ~/.ssh/allowed_signers
git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers
```

### 重要: user.signingkey は公開鍵を指定

Git の `user.signingkey` には**公開鍵のパス**を指定する。Git/ssh-agent が対応する秘密鍵を自動的に見つけて署名に使用する。

設定オプション:
1. **公開鍵ファイルのパス** (推奨): `~/.ssh/id_ed25519.pub`
2. **公開鍵のリテラル文字列**: `ssh-ed25519 AAAAC3...`
3. **ssh-agent から動的取得**: `git config --global gpg.ssh.defaultKeyCommand "ssh-add -L"`

## 完全な docker-compose.yml の例

```yaml
version: "3.8"

services:
  coding-agent:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: coding-agent

    stdin_open: true
    tty: true

    user: agent
    working_dir: /workspace

    volumes:
      - ./:/workspace:cached
      # SSH 公開鍵のみ (読み取り専用)
      - ${HOME}/.ssh/id_ed25519.pub:/home/agent/.ssh-host/id_ed25519.pub:ro
      - ${HOME}/.ssh/known_hosts:/home/agent/.ssh-host/known_hosts:ro

    environment:
      - TERM=xterm-256color
      - SSH_AUTH_SOCK=/tmp/ssh-agent.sock
      - GIT_USER_NAME=${GIT_USER_NAME:-}
      - GIT_USER_EMAIL=${GIT_USER_EMAIL:-}
      - GIT_SIGNING_KEY=/home/agent/.ssh-host/id_ed25519.pub

    # セキュリティ強化
    security_opt:
      - no-new-privileges:true
    cap_drop:
      - ALL
    cap_add:
      - NET_RAW

volumes:
  agent-config:
```

## 使用方法

### 1. コンテナの起動

```bash
# プラットフォーム自動検出で起動
./run.sh up

# または直接 docker compose
docker compose -f docker-compose.yml -f docker-compose.macos.yml up -d
```

### 2. 署名付きコミットのテスト

```bash
# コンテナ内でシェルを起動
./run.sh shell

# SSH agent の確認
ssh-add -l

# GitHub への接続テスト
ssh -T git@github.com

# 署名付きコミット
echo "test" > test.txt
git add test.txt
git commit -S -m "SSH signed commit via Docker"

# 署名の検証
git log --show-signature -1

# プッシュ
git push
```

GitHub 上で **Verified** バッジが表示されれば成功。

## トラブルシューティング

### SSH agent にアクセスできない

```bash
# ホスト側で確認
ssh-add -l
# 出力がない場合: ssh-add ~/.ssh/id_ed25519

# コンテナ内で確認
echo $SSH_AUTH_SOCK
ls -la $SSH_AUTH_SOCK
```

### "Permission denied" エラー

Docker Desktop でソケットファイルの権限問題が発生する場合:

```bash
docker compose exec --user root coding-agent \
  chown agent /tmp/ssh-agent.sock
```

### 署名が検証されない

1. GitHub に SSH キーを**署名用**として登録しているか確認
   - Settings → SSH and GPG keys → "Signing Key" として追加
2. `git config user.signingkey` が正しいパスを指しているか確認
3. `git config gpg.format` が `ssh` に設定されているか確認

### Rancher Desktop で SSH Agent Forwarding が動作しない

Rancher Desktop は Lima VM を使用するため、macOS の Unix ソケットを直接転送できない。以下の回避策がある:

1. **秘密鍵を直接マウント** (セキュリティ上のトレードオフ)
2. **Docker Desktop に切り替え** (推奨)

## ベストプラクティス

1. **秘密鍵をコンテナにマウントしない** - ソケット転送のみを使用
2. **信頼できるコンテナでのみ SSH Agent Forwarding を有効化**
3. **作業完了後は `ssh-add -D` で鍵を削除可能**
4. **`no-new-privileges` と `CAP_DROP ALL` でコンテナを実行**
5. **公開鍵は読み取り専用 (`:ro`) でマウント**

## 参考リンク

### 公式ドキュメント
- [GitHub Docs: Telling Git about your signing key](https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key)
- [GitHub Docs: Signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)

### SSH Agent Forwarding
- [Forward host ssh-agent to docker builds on macOS](https://michael.kefeder.at/post/ssh-agent-docker-compose/)
- [SSH agent forward into docker container on macOS](https://medium.com/@nazrulworld/ssh-agent-forward-into-docker-container-on-macos-ff847ec660e2)
- [Docker for Mac: SSH agent forwarding issue](https://github.com/docker/for-mac/issues/5303)

### Git SSH 署名
- [Git: The complete guide to sign your commits with an SSH key](https://dev.to/ccoveille/git-the-complete-guide-to-sign-your-commits-with-an-ssh-key-35bg)
- [(Correctly) Telling git about your SSH key for signing commits](https://dev.to/li/correctly-telling-git-about-your-ssh-key-for-signing-commits-4c2c)
- [Setting Up SSH for Commit Signing - Tower Blog](https://www.git-tower.com/blog/setting-up-ssh-for-commit-signing)

### Rancher Desktop
- [Rancher Desktop: SSH agent forwarding discussion](https://github.com/rancher-sandbox/rancher-desktop/discussions/1842)
- [Rancher Desktop: SSH agent socket issue](https://github.com/rancher-sandbox/rancher-desktop/issues/3042)

### セキュリティ
- [Securely Using SSH Keys in Docker](https://www.fastruby.io/blog/docker/docker-ssh-keys.html)
- [1Password: Sign Git commits with SSH](https://developer.1password.com/docs/ssh/git-commit-signing/)
