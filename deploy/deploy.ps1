param(
    [string]$AccountId = "",  # 請填入你的 AWS Account ID
    [string]$Region = "ap-northeast-1",
    [string]$BucketName = "aws-learning-game-frontend",
    [string]$RepoName = "aws-learning-game"
)

$ErrorActionPreference = "Stop"
$ecrUrl = "$AccountId.dkr.ecr.$Region.amazonaws.com"

Write-Host "=== AWS Learning Game Deploy ===" -ForegroundColor Cyan
Write-Host "Account: $AccountId | Region: $Region"

Write-Host "`n[1/6] Create ECR Repository..." -ForegroundColor Yellow
try {
    aws ecr describe-repositories --repository-names $RepoName --region $Region 2>$null | Out-Null
    Write-Host "ECR Repository exists" -ForegroundColor Green
} catch {
    aws ecr create-repository --repository-name $RepoName --region $Region
    Write-Host "ECR Repository created" -ForegroundColor Green
}

Write-Host "`n[2/6] Login to ECR..." -ForegroundColor Yellow
aws ecr get-login-password --region $Region | docker login --username AWS --password-stdin $ecrUrl

Write-Host "`n[3/6] Build and push Docker image..." -ForegroundColor Yellow
docker build -t $RepoName .
docker tag "$($RepoName):latest" "$ecrUrl/$($RepoName):latest"
docker push "$ecrUrl/$($RepoName):latest"
Write-Host "Docker image pushed" -ForegroundColor Green

Write-Host "`n[4/6] Create S3 Bucket..." -ForegroundColor Yellow
$bucketCheck = aws s3 ls "s3://$BucketName" 2>&1
if ($LASTEXITCODE -ne 0) {
    aws s3 mb "s3://$BucketName" --region $Region
    Write-Host "S3 Bucket created" -ForegroundColor Green
} else {
    Write-Host "S3 Bucket exists" -ForegroundColor Green
}

Write-Host "`n[5/6] Build frontend..." -ForegroundColor Yellow
Push-Location web
npm install
npm run build
Pop-Location
Write-Host "Frontend built" -ForegroundColor Green

Write-Host "`n[6/6] Upload to S3..." -ForegroundColor Yellow
aws s3 sync web/dist/ "s3://$BucketName" --delete
Write-Host "Frontend uploaded" -ForegroundColor Green

Write-Host "`n=== Deploy Complete ===" -ForegroundColor Cyan
Write-Host "Next: Create App Runner service in AWS Console"
Write-Host "Image: $ecrUrl/$($RepoName):latest"
Write-Host "Port: 8080"
