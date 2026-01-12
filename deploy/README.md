# AWS 部署指南

## 架構概覽

```
┌─────────────────┐     ┌─────────────────┐
│   CloudFront    │────▶│       S3        │
│   (CDN/HTTPS)   │     │  (React 前端)   │
└────────┬────────┘     └─────────────────┘
         │
         │ /api/*
         ▼
┌─────────────────┐     ┌─────────────────┐
│   App Runner    │────▶│    DynamoDB     │
│   (Go 後端)     │     │   (遊戲存檔)    │
└─────────────────┘     └─────────────────┘
```

## 部署步驟

### 1. 前置準備

確保已安裝：
- AWS CLI (`aws --version`)
- Docker (`docker --version`)
- Node.js (`node --version`)

設定 AWS 認證：
```bash
aws configure
```

### 2. 部署後端 (App Runner)

#### 2.1 建立 ECR Repository
```bash
aws ecr create-repository --repository-name aws-learning-game --region ap-northeast-1
```

#### 2.2 登入 ECR
```bash
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com
```

#### 2.3 建立並推送 Docker Image
```bash
docker build -t aws-learning-game .
docker tag aws-learning-game:latest <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/aws-learning-game:latest
docker push <ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com/aws-learning-game:latest
```

#### 2.4 建立 App Runner 服務

使用 AWS Console:
1. 前往 App Runner
2. 建立服務 → 選擇 Container registry → Amazon ECR
3. 選擇剛推送的 image
4. 設定：
   - CPU: 0.25 vCPU
   - Memory: 0.5 GB
   - Port: 8080
5. 環境變數（可選）：
   - `GIN_MODE`: release
   - `DYNAMODB_TABLE`: aws-learning-game
   - `AWS_REGION`: ap-northeast-1

或使用 CLI:
```bash
aws apprunner create-service --cli-input-json file://deploy/apprunner.json
```

### 3. 部署前端 (S3 + CloudFront)

#### 3.1 建立 S3 Bucket
```bash
aws s3 mb s3://aws-learning-game-frontend --region ap-northeast-1
```

#### 3.2 設定 S3 靜態網站託管
```bash
aws s3 website s3://aws-learning-game-frontend --index-document index.html --error-document index.html
```

#### 3.3 建置前端
```bash
cd web
npm install
# 設定 API URL (替換為你的 App Runner URL)
echo "VITE_API_URL=https://xxxxxx.ap-northeast-1.awsapprunner.com" > .env.production
npm run build
```

#### 3.4 上傳到 S3
```bash
aws s3 sync dist/ s3://aws-learning-game-frontend --delete
```

#### 3.5 建立 CloudFront Distribution

使用 AWS Console:
1. 前往 CloudFront → 建立分佈
2. Origin domain: 選擇 S3 bucket
3. Origin access: Origin access control settings (推薦)
4. 建立 OAC 並更新 S3 bucket policy
5. Default root object: index.html
6. 設定自訂錯誤回應：
   - 403 → /index.html (200)
   - 404 → /index.html (200)

### 4. 設定 DynamoDB（可選，用於遊戲存檔）

```bash
aws dynamodb create-table \
  --table-name aws-learning-game \
  --attribute-definitions AttributeName=game_id,AttributeType=S \
  --key-schema AttributeName=game_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region ap-northeast-1
```

### 5. 更新 App Runner IAM Role

如果使用 DynamoDB，需要給 App Runner 服務角色添加 DynamoDB 權限：

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:DeleteItem",
        "dynamodb:Query",
        "dynamodb:Scan"
      ],
      "Resource": "arn:aws:dynamodb:ap-northeast-1:*:table/aws-learning-game"
    }
  ]
}
```

## 成本估算

| 服務 | 預估月費 |
|------|----------|
| App Runner (0.25 vCPU, 0.5GB) | ~$5-15 |
| S3 (靜態網站) | < $1 |
| CloudFront | < $1 (低流量) |
| DynamoDB (按需) | < $1 (低使用量) |
| **總計** | **~$7-18/月** |

## 清理資源

```bash
# 刪除 App Runner 服務
aws apprunner delete-service --service-arn <SERVICE_ARN>

# 清空並刪除 S3
aws s3 rm s3://aws-learning-game-frontend --recursive
aws s3 rb s3://aws-learning-game-frontend

# 刪除 CloudFront (需先停用)
aws cloudfront delete-distribution --id <DISTRIBUTION_ID>

# 刪除 ECR
aws ecr delete-repository --repository-name aws-learning-game --force

# 刪除 DynamoDB
aws dynamodb delete-table --table-name aws-learning-game
```
