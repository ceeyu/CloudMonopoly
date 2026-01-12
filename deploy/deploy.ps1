# AWS Learning Game 部署腳本 (PowerShell)
# 使用前請先設定以下變數

param(
    [Parameter(Mandatory=$true)]
    [string]$AccountId,
    
    [string]$Region = "ap-northeast-1",
    [string]$BucketName = "aws-learning-game-frontend",
    [string]$RepoName = "aws-learning-game"
)

$ErrorActionPreference = "Stop"

Write-Host "=== AWS Learning Game 部署 ===" -ForegroundColor Cyan

# 1. 建立 ECR Repository
Write-Host "`n[1/6] 建立 ECR Repository..." -ForegroundColor Yellow
aws ecr describe-repositories --repository-names $RepoName --region $Region 2>$null
if ($LASTEXITCODE -ne 0) {
    aws ecr create-repository --repository-name $RepoName --region $Region
    Write-Host "ECR Repository 已建立" -ForegroundColor Green
} else {
    Write-Host "ECR Repository 已存在" -ForegroundColor Green
}

# 2. 登入 ECR
Write-Host "`n[2/6] 登入 ECR..." -ForegroundColor Yellow
$ecrPassword = aws ecr get-login-password --region $Region
docker login --username AWS --password $ecrPassword "$AccountId.dkr.ecr.$Region.amazonaws.com"

# 3. 建立並推送 Docker Image
Write-Host "`n[3/6] 建立並推送 Docker Image..." -ForegroundColor Yellow
docker build -t $RepoName .
docker tag "${RepoName}:latest" "$AccountId.dkr.ecr.$Region.amazonaws.com/${RepoName}:latest"
docker push "$AccountId.dkr.ecr.$Region.amazonaws.com/${RepoName}:latest"
Write-Host "Docker Image 已推送" -ForegroundColor Green

# 4. 建立 S3 Bucket
Write-Host "`n[4/6] 建立 S3 Bucket..." -ForegroundColor Yellow
aws s3 ls "s3://$BucketName" 2>$null
if ($LASTEXITCODE -ne 0) {
    aws s3 mb "s3://$BucketName" --region $Region
    Write-Host "S3 Bucket 已建立" -ForegroundColor Green
} else {
    Write-Host "S3 Bucket 已存在" -ForegroundColor Green
}

# 5. 建置前端
Write-Host "`n[5/6] 建置前端..." -ForegroundColor Yellow
Push-Location web
npm install
npm run build
Pop-Location
Write-Host "前端建置完成" -ForegroundColor Green

# 6. 上傳到 S3
Write-Host "`n[6/6] 上傳前端到 S3..." -ForegroundColor Yellow
aws s3 sync web/dist/ "s3://$BucketName" --delete
Write-Host "前端已上傳" -ForegroundColor Green

Write-Host "`n=== 部署完成 ===" -ForegroundColor Cyan
Write-Host @"

下一步：
1. 前往 AWS Console -> App Runner 建立服務
   - Image: $AccountId.dkr.ecr.$Region.amazonaws.com/${RepoName}:latest
   - Port: 8080

2. 取得 App Runner URL 後，更新前端環境變數：
   cd web
   echo "VITE_API_URL=https://xxxxxx.$Region.awsapprunner.com" > .env.production
   npm run build
   aws s3 sync dist/ s3://$BucketName --delete

3. 建立 CloudFront Distribution 指向 S3 Bucket

"@ -ForegroundColor White
