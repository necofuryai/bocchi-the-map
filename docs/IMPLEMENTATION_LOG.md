# Bocchi The Map 実装ログ

## 2025年5月29日 - 初期実装

### 実装した内容

#### 1. プロジェクト全体
- ✅ `.gitignore`ファイル作成
- ✅ `README.md`作成（プロジェクト概要）
- ✅ `CLAUDE.md`改善（AI開発ガイド）

#### 2. API (Golang)
- ✅ ディレクトリ構造（オニオンアーキテクチャ）
- ✅ Protocol Buffers定義
  - `common.proto` - 共通型定義
  - `spot.proto` - スポット管理
  - `review.proto` - レビュー管理
  - `user.proto` - ユーザー管理
- ✅ Humaフレームワークでの基本実装
  - ヘルスチェックエンドポイント
  - スポット管理API（Create/Get/List）
- ✅ JSON構造化ログ（zerolog）
- ✅ 設定管理（環境変数）

#### 3. Web (Next.js)
- ✅ Next.js + TypeScript初期化
- ✅ 依存関係インストール
  - NextAuth.js（認証）
  - next-themes（ダークモード）
  - Tailwind CSS + Shadcn/ui
- ✅ 認証設定（Google/X OAuth）
- ✅ テーマプロバイダー（ダークモード対応）
- ✅ 型定義ファイル

#### 4. Infrastructure (Terraform)
- ✅ モジュール構造
  - Cloudflare R2（地図タイル保存）
  - Google Cloud Run（API）
- ✅ 環境別設定（dev/prod）

### 技術スタック決定事項

**バックエンド**
- 言語: Go 1.21+
- フレームワーク: Huma v2
- アーキテクチャ: オニオンアーキテクチャ
- API設計: Protocol Buffers
- データベース: TiDB Serverless

**フロントエンド**
- フレームワーク: Next.js 15 (App Router)
- 言語: TypeScript
- スタイリング: Tailwind CSS + Shadcn/ui
- 認証: NextAuth.js
- 地図: MapLibre GL JS

**インフラ**
- API: Google Cloud Run
- 静的ホスティング: Vercel
- 地図データ: Cloudflare R2
- IaC: Terraform

### 次のステップ

1. **高優先度**
   - [ ] TiDBデータベース接続実装
   - [ ] MapLibre GL JS統合
   - [ ] 基本的なUI実装

2. **中優先度**
   - [ ] 多言語対応（i18n）
   - [ ] API認証実装
   - [ ] レビュー投稿機能

3. **低優先度**
   - [ ] モニタリング設定（New Relic/Sentry）
   - [ ] CI/CD設定
   - [ ] E2Eテスト

### 開発メモ

- React 19との依存関係問題があるため、一部ライブラリは`--legacy-peer-deps`フラグが必要
- MapLibre GL JSの地図データはPMTiles形式でCloudflare R2に保存予定
- 将来的な拡張性を考慮し、マイクロサービス化しやすい設計を採用
