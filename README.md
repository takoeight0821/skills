# takoeight0821-skills

Claude Code Plugin: 研究・開発向けスキルコレクション

## 概要

このリポジトリは、Claude Codeのプラグインとして使用できるスキル集です。GitHubリポジトリ、ソフトウェアプロジェクト、学術論文の調査などをサポートします。

## インストール

Claude Codeでこのプラグインを使用するには、`~/.claude/settings.json`に追加します：

```json
{
  "plugins": [
    "https://github.com/takoeight0821/skills"
  ]
}
```

または、Claude Codeのコマンドでインストール：

```bash
claude plugin add https://github.com/takoeight0821/skills
```

## 含まれるスキル

### /research
GitHubリポジトリ、ソフトウェアプロジェクト、学術論文を調査し、包括的なマークダウンレポートを生成します。

### /review-plan
コード内の`review:`または`review(username):`コメントを検索し、コードレビュー指摘への対応計画を作成します。

### /clean-comments
コード内のコメントを最適化し、不要なコメントを削除、必要なコメントを改善します。

## ディレクトリ構造

```
/
├── .claude-plugin/     # プラグイン設定
│   └── plugin.json     # プラグインメタデータ
├── skills/             # Claude Codeスキル
│   ├── clean-comments/ # コメント最適化スキル
│   ├── research/       # リサーチスキル
│   └── review-plan/    # レビュー対応スキル
├── conductor/          # プロジェクト管理ドキュメント
├── CLAUDE.md           # Claude Code向けガイダンス
└── LICENSE             # MITライセンス
```

## スキルの追加

新しいスキルを追加するには、`skills/`ディレクトリに新しいフォルダを作成し、`SKILL.md`ファイルを配置します：

```bash
mkdir -p skills/my-skill
```

`SKILL.md`にはYAMLフロントマター（name, description）とスキルの内容をマークダウンで記述します。

## 関連リポジトリ

- [takoeight0821/jig](https://github.com/takoeight0821/jig) - VM/Docker管理、スキル同期を行うCLIツール（このリポジトリから分離）

## ライセンス

MIT License
