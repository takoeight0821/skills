# Product Definition

## Overview

Claude Code Plugin: 研究・開発向けスキルコレクション

このリポジトリは、Claude Codeのプラグインとして動作するスキル集です。

## Core Features

### Skills

- **research**: GitHubリポジトリ、ソフトウェアプロジェクト、学術論文を調査し、包括的なマークダウンレポートを生成
- **review-plan**: コード内の`review:`コメントを検索し、コードレビュー指摘への対応計画を作成
- **clean-comments**: コード内のコメントを最適化し、不要なコメントを削除

### Plugin Configuration

- `.claude-plugin/plugin.json`: プラグインメタデータ（名前、バージョン、説明など）

## Target Audience

- Claude Codeを使用する開発者
- GitHubリポジトリや学術論文の調査が必要なユーザー
- コードレビューのワークフローを効率化したいチーム

## User Experience

- **Skill-Based**: `/research`, `/review-plan`, `/clean-comments`などのスキルコマンドで機能を呼び出し
- **Plugin Integration**: Claude Codeのプラグインとしてシームレスに統合
