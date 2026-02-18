# AWS → Google Cloud 移行ガイド

本ドキュメントは、gomethod アプリケーションを AWS から Google Cloud へ移行するための手順をまとめたものです。

## 前提・方針

### 使用する Google アカウント

**本アカウント（Google AI Ultra 契約者）** でプロジェクトを作成・運用する。

| 項目 | 本アカウント | ファミリーアカウント |
|------|-------------|-------------------|
| $100/月 クレジット | ✅ サブスク契約に直接紐づく | ⚠️ ファミリー特典として付与 |
| Ultra 解約時 | クレジット失効（想定内） | **クレジット失効リスクあり** |
| 推奨度 | ✅ 推奨 | ❌ 非推奨 |

> [!IMPORTANT]
> Google AI Ultra を将来解約する可能性がある場合、ファミリーアカウントの $100 追加特典は消滅する可能性が高い。安全のため本アカウントで運用すること。

---

## 移行ステップ概要

### Phase 1: Google Cloud プロジェクトのセットアップ
- [x] 本アカウントで Google Cloud Console にログイン
- [x] 新規プロジェクトを作成（例: `gomethod`）
- [ ] 課金アカウントの確認（$100/月クレジットが適用されているか）
- [ ] 必要な API の有効化
  - Cloud SQL Admin API
  - Cloud Run API
  - Container Registry / Artifact Registry API
  - Secret Manager API

### Phase 2: Cloud SQL（MySQL）の構築
- [ ] Cloud SQL for MySQL インスタンスを作成
  - リージョン: `asia-northeast1`（東京）
  - マシンタイプ: 最小構成で開始（コスト最適化）
  - MySQL バージョン: 既存 RDS と同じバージョンに合わせる
- [ ] データベースとユーザーを作成
- [ ] ネットワーク設定（プライベート IP or 公開 IP + 承認済みネットワーク）

### Phase 3: データ移行（AWS RDS → Cloud SQL）
- [ ] AWS RDS から `mysqldump` でデータエクスポート
  ```bash
  mysqldump -h <RDS_ENDPOINT> -u <USER> -p <DB_NAME> > gmethod_dump.sql
  ```
- [ ] Cloud SQL にデータインポート
  - 方法 A: `gcloud sql import` で GCS 経由
  - 方法 B: Cloud SQL Auth Proxy 経由で直接 `mysql` コマンド
- [ ] データの整合性確認（レコード数、サンプルデータ検証）

### Phase 4: Cloud Run デプロイ
- [ ] Dockerfile の作成（backend/ 配下）
- [ ] Artifact Registry にコンテナイメージを push
- [ ] Cloud Run サービスの作成
- [ ] 環境変数の設定（Secret Manager 連携）
  - `LINE_CHANNEL_TOKEN`
  - `LINE_CHANNEL_SECRET`
  - `GMETHOD_DB_HOST` / `GMETHOD_DB_USERNAME` / `GMETHOD_DB_PASSWORD`
  - `OPEN_AI_API_KEY` / `OPEN_AI_COMPLETION_ENDPOINT`
  - `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`（S3 継続利用の場合）
- [ ] LINE Webhook URL の変更

### Phase 5: 検証・切り替え
- [ ] Cloud Run 上での動作確認
- [ ] LINE Bot のエンドポイント切り替え
- [ ] 旧 AWS リソースの停止・削除

---

## コスト見積もり

| リソース | 月額見積もり | 備考 |
|---------|------------|------|
| Cloud SQL (db-f1-micro) | ~$10–15 | 最小構成、東京リージョン |
| Cloud Run | ~$0–5 | リクエストベース課金、低トラフィック想定 |
| Artifact Registry | ~$1 未満 | コンテナイメージ保存 |
| **合計** | **~$15–20** | **$100 クレジット内で十分収まる** |

---

## 次のアクション

1. **Google Cloud プロジェクトの作成**（本アカウントで）
2. **Cloud SQL インスタンスの作成**（Terraform で管理）
3. **RDS からの mysqldump 取得**
