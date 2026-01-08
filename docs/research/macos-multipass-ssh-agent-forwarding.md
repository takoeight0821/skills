# macOS × Multipass × cloud-init × SSH Agent Forwarding で安全に署名コミットする方法

## 概要

| 項目 | 内容 |
|------|------|
| リポジトリ/プロジェクト | Multipass + cloud-init 環境構築 |
| 開発元 | [Canonical](https://multipass.run/) (Multipass), [cloud-init.io](https://cloudinit.readthedocs.io/) |
| ライセンス | GPL-3.0 (Multipass) |
| 対象プラットフォーム | macOS (Apple Silicon / Intel) |
| 関連技術 | SSH Agent Forwarding, Git SSH Commit Signing |

## ゴール

- macOS 側の SSH 秘密鍵を **Multipass VM にコピーしない**
- **SSH Agent Forwarding** によって、VM内の `git commit -S` が macOS 側の鍵で署名
- GitHub の **SSH commit signing** に対応
- cloud-init で **VM起動後すぐに署名できる環境**が自動で揃う

## アーキテクチャ

```
macOS (ホスト)
├── ~/.ssh/id_ed25519（秘密鍵 - ここにのみ存在）
├── ssh-agent が秘密鍵を保持
├── ssh -A ubuntu@<VM IP> で接続
└── cloud-init.yaml でVM側を自動セットアップ

        ↓ SSH Agent Forwarding（秘密鍵は転送されない）

Multipass VM（Ubuntu）
├── ~/.ssh に秘密鍵は存在しない
├── 公開鍵のみ .gitconfig に設定
├── git commit -S 実行時
└── → macOS の agent を経由して署名
```

## Multipass について

### 概要

[Multipass](https://multipass.run/) は Canonical が開発した、Ubuntu VM を素早く起動・管理するためのツール。macOS、Linux、Windows で動作し、シンプルな CLI で VM のライフサイクルを管理できる。

### 主な特徴

- **軽量**: QEMU/HyperKit/Hyper-V を使用した高速な VM 起動
- **cloud-init 対応**: 起動時の自動構成をサポート
- **マルチプラットフォーム**: macOS (Intel/Apple Silicon)、Linux、Windows をサポート
- **Ubuntu 公式**: Canonical が直接メンテナンス

### macOS へのインストール

```bash
# Homebrew でインストール
brew install --cask multipass

# バージョン確認
multipass version
```

### 基本コマンド

```bash
# VM 一覧表示
multipass list

# VM 作成（デフォルト設定）
multipass launch --name myvm

# VM 作成（カスタム設定）
multipass launch --name devvm --cpus 2 --memory 4G --disk 20G

# cloud-init を使用した VM 作成
multipass launch --name devvm --cloud-init cloud-init.yaml

# VM にシェル接続
multipass shell devvm

# VM でコマンド実行
multipass exec devvm -- ls -la

# VM の情報表示
multipass info devvm

# VM 停止/起動/削除
multipass stop devvm
multipass start devvm
multipass delete devvm && multipass purge

# ファイル転送
multipass transfer local-file.txt devvm:/home/ubuntu/
multipass transfer devvm:/home/ubuntu/remote-file.txt ./
```

## SSH Agent Forwarding について

### 仕組み

SSH Agent Forwarding は、SSH 接続を通じてローカルマシンの SSH エージェントをリモートホストで使用できるようにする機能。**秘密鍵自体はリモートサーバーに送信されない**。代わりに、リモートサーバーはローカルマシンの SSH エージェントに署名操作を依頼する。

### セキュリティ上の考慮事項

SSH Agent Forwarding には固有のセキュリティリスクがある：

> "Users with ability to bypass file permissions on the remote host can access the local agent through the forwarded connection." - SSH manual

リモートホストの root ユーザーは、転送された SSH エージェントにアクセスし、他のサーバーへの認証に使用する可能性がある。

### ベストプラクティス

1. **デフォルトで有効にしない**
   ```bash
   # 良い例: 必要な時だけ -A オプションを使用
   ssh -A ubuntu@<VM_IP>

   # 悪い例: ~/.ssh/config で常に有効化
   # ForwardAgent yes  ← これは避ける
   ```

2. **タイムアウトを設定する**
   ```bash
   ssh-add -t 3600 ~/.ssh/id_ed25519  # 1時間後に自動削除
   ```

3. **エージェントをロックする**
   ```bash
   ssh-add -x   # パスワードでロック
   ssh-add -X   # アンロック
   ```

4. **確認を要求する**
   ```bash
   ssh-add -c ~/.ssh/id_ed25519  # 使用時に確認を要求
   ```

5. **未使用の鍵を削除する**
   ```bash
   ssh-add -d ~/.ssh/id_ed25519  # 特定の鍵を削除
   ssh-add -D                     # すべての鍵を削除
   ```

6. **エージェントを終了する**
   ```bash
   eval "$(ssh-agent -k)"
   ```

### より安全な代替手段

**ProxyJump (-J) オプション**を使用すると、Agent Forwarding なしで中間ホストを経由できる：

```bash
ssh -J jump-host target-host
```

ただし、Multipass VM の場合は通常、Agent Forwarding が最も実用的な選択肢となる。

## GitHub SSH Commit Signing について

### 概要

Git 2.34 以降では、GPG 鍵の代わりに SSH 鍵を使用してコミットに署名できる。これにより、既存の SSH 鍵を認証と署名の両方に使用できる。

### 必要な Git 設定

```bash
# 署名フォーマットを SSH に設定
git config --global gpg.format ssh

# 署名に使用する公開鍵を指定
git config --global user.signingkey ~/.ssh/id_ed25519.pub

# コミット時の自動署名を有効化
git config --global commit.gpgsign true

# タグの自動署名を有効化（オプション）
git config --global tag.gpgsign true
```

### ローカル検証の設定

```bash
# 許可された署名者ファイルを設定
git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers

# 許可された署名者を追加
echo "your-email@example.com $(cat ~/.ssh/id_ed25519.pub)" >> ~/.ssh/allowed_signers
```

### GitHub への登録

1. GitHub Settings → SSH and GPG keys
2. "New SSH key" をクリック
3. **Key type: "Signing Key"** を選択
4. 公開鍵を貼り付けて保存

**注意**: 同じ SSH 鍵を認証とコミット署名の両方に使用する場合、両方のカテゴリに明示的に登録する必要がある。

## cloud-init について

### 概要

[cloud-init](https://cloudinit.readthedocs.io/) はクラウドインスタンスの初期設定を自動化するための標準ツール。Multipass でも cloud-init を使用して VM の起動時設定をカスタマイズできる。

### 主要モジュール

#### packages

パッケージのインストールを自動化：

```yaml
#cloud-config
package_update: true
packages:
  - git
  - openssh-client
  - curl
```

#### write_files

ファイルの作成・編集：

```yaml
#cloud-config
write_files:
  - path: /home/ubuntu/.gitconfig
    permissions: "0644"
    owner: ubuntu:ubuntu
    content: |
      [user]
        name = Your Name
        email = you@example.com
```

**オプション**:
- `path`: ファイルパス（必須）
- `content`: ファイル内容
- `permissions`: パーミッション（8進数文字列）
- `owner`: 所有者（user:group）
- `encoding`: エンコーディング（b64, gzip）
- `append`: 追記モード（true/false）
- `defer`: パッケージインストール後まで遅延（true/false）

#### runcmd

コマンドの実行：

```yaml
#cloud-config
runcmd:
  - sudo -u ubuntu mkdir -p /home/ubuntu/.ssh
  - sudo -u ubuntu chmod 700 /home/ubuntu/.ssh
  - ssh-keyscan github.com >> /home/ubuntu/.ssh/known_hosts
```

### Multipass との連携

```bash
# cloud-init.yaml を使用して VM を起動
multipass launch --name devvm --cloud-init cloud-init.yaml
```

Multipass は YAML の構文を検証し、エラーがあれば起動前に通知する。

## 実装手順

### 1. macOS 側の準備

```bash
# SSH エージェントを起動
eval "$(ssh-agent -s)"

# 秘密鍵をエージェントに登録（1時間のタイムアウト付き）
ssh-add -t 3600 ~/.ssh/id_ed25519

# 鍵が登録されたことを確認
ssh-add -l

# 公開鍵を環境変数に設定（cloud-init で使用）
export GIT_SSH_SIGNING_PUBKEY="$(cat ~/.ssh/id_ed25519.pub)"
```

### 2. cloud-init.yaml の作成

```yaml
#cloud-config
package_update: true
packages:
  - git
  - openssh-client

write_files:
  - path: /home/ubuntu/.gitconfig
    permissions: "0644"
    owner: ubuntu:ubuntu
    content: |
      [user]
        name = Your Name
        email = you@example.com
        signingkey = ssh-ed25519 AAAA... # 公開鍵をここに記載
      [commit]
        gpgsign = true
      [gpg]
        format = ssh

runcmd:
  - sudo -u ubuntu mkdir -p /home/ubuntu/.ssh
  - sudo -u ubuntu chmod 700 /home/ubuntu/.ssh
  - ssh-keyscan github.com >> /home/ubuntu/.ssh/known_hosts
  - chown ubuntu:ubuntu /home/ubuntu/.ssh/known_hosts
```

**ポイント**:
- `gpg.format = ssh` で SSH 署名を有効化
- `user.signingkey` には公開鍵を設定（秘密鍵は不要）
- `known_hosts` に GitHub を追加して初回接続時の確認を省略

### 3. Multipass VM の作成

```bash
multipass launch \
  --name devvm \
  --memory 4G \
  --disk 20G \
  --cpus 2 \
  --cloud-init cloud-init.yaml
```

### 4. VM の IP アドレスを確認

```bash
multipass list
# または
multipass info devvm | grep IPv4
```

### 5. SSH Agent Forwarding で接続

```bash
# -A オプションが重要（Agent Forwarding を有効化）
ssh -A ubuntu@$(multipass info devvm | awk '/IPv4/ {print $2}')
```

**~/.ssh/config を使用する場合**:

```
Host devvm
  HostName 192.168.64.X  # multipass list で確認した IP
  User ubuntu
  ForwardAgent yes
```

```bash
ssh devvm
```

### 6. VM 内で署名コミットをテスト

```bash
# リポジトリをクローン
git clone git@github.com:username/repo.git
cd repo

# 変更を加えてコミット
echo "test" > test.txt
git add test.txt
git commit -S -m "Signed commit via Multipass + SSH Agent Forwarding"

# プッシュ
git push origin main
```

GitHub で "Verified" バッジが表示されれば成功。

## トラブルシューティング

### Agent Forwarding が動作しない

1. **macOS 側で鍵が登録されているか確認**
   ```bash
   ssh-add -l
   ```

2. **SSH_AUTH_SOCK が設定されているか確認**
   ```bash
   echo $SSH_AUTH_SOCK
   ```

3. **VM 内で Agent Forwarding が有効か確認**
   ```bash
   # VM 内で実行
   echo $SSH_AUTH_SOCK
   ssh-add -l  # macOS 側の鍵が表示されるはず
   ```

### Git 署名が失敗する

1. **gpg.format が ssh に設定されているか確認**
   ```bash
   git config --get gpg.format
   ```

2. **signingkey が正しく設定されているか確認**
   ```bash
   git config --get user.signingkey
   ```

3. **公開鍵が GitHub に登録されているか確認**
   - Settings → SSH and GPG keys → Signing Keys

### VM に接続できない

1. **VM が起動しているか確認**
   ```bash
   multipass list
   ```

2. **VM の IP アドレスを確認**
   ```bash
   multipass info devvm
   ```

3. **macOS のファイアウォール設定を確認**

## セキュリティのベストプラクティス

| 推奨事項 | 説明 |
|----------|------|
| 秘密鍵は macOS にのみ保持 | VM には絶対にコピーしない |
| タイムアウトを設定 | `ssh-add -t` で一定時間後に自動削除 |
| 信頼できる VM のみに転送 | Agent Forwarding は慎重に使用 |
| 作業後は鍵を削除 | `ssh-add -D` でクリア |
| cloud-init は公開鍵のみ | 公開鍵は漏洩しても問題ない |

## まとめ

| 要素 | 説明 |
|------|------|
| セキュリティ設計 | 秘密鍵は macOS の SSH Agent にのみ存在 |
| VM 構成 | cloud-init で Git + SSH署名設定が自動適用 |
| 署名方式 | `gpg.format = ssh` |
| ログイン方法 | `ssh -A ubuntu@<VM_IP>` |
| 署名動作 | VM → macOS のエージェントへ転送して署名 |
| GitHub | SSH commit signing に公開鍵を登録 |

## 関連リソース

### 公式ドキュメント
- [Multipass Documentation](https://documentation.ubuntu.com/multipass/en/latest)
- [cloud-init Documentation](https://cloudinit.readthedocs.io/)
- [GitHub - Telling Git about your signing key](https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key)
- [GitHub - Signing commits](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits)

### チュートリアル・ガイド
- [Using cloud-init with Multipass | Ubuntu Blog](https://ubuntu.com/blog/using-cloud-init-with-multipass)
- [Basic Multipass setup on MacOS | Johnny Matthews](https://johnnymatthews.dev/blog/2025-08-07-basic-multipass-setup-on-macos/)
- [Multipass VM Manager Cheatsheet | Rost Glukhov](https://www.glukhov.org/post/2025/10/vm-manager-multipass-cheatsheet/)
- [cloud-init Examples | Caltech Library](https://github.com/caltechlibrary/cloud-init-examples)

### セキュリティ
- [Safer SSH agent forwarding | Vincent Bernat](https://vincent.bernat.ch/en/blog/2020-safer-ssh-agent-forwarding)
- [5 SSH Agent Best Practices | Teleport](https://goteleport.com/blog/how-to-use-ssh-agent-safely/)
- [SSH Agent Explained | Smallstep](https://smallstep.com/blog/ssh-agent-explained/)

### Git SSH Signing
- [Sign Git commits with SSH | 1Password Developer](https://developer.1password.com/docs/ssh/git-commit-signing/)
- [Signing Git commits with SSH keys | Emmanuel Bernard](https://emmanuelbernard.com/blog/2023/11/27/git-signing-ssh/)
- [Setting Up SSH for Commit Signing | Tower Blog](https://www.git-tower.com/blog/setting-up-ssh-for-commit-signing)
- [Git: The complete guide to sign your commits with an SSH key | DEV Community](https://dev.to/ccoveille/git-the-complete-guide-to-sign-your-commits-with-an-ssh-key-35bg)
