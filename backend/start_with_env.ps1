#!/usr/bin/env pwsh

# 完整的环境变量设置和启动命令

# 设置所有必要的环境变量
$env:POSTGRES_PASSWORD = 'panda123456'
$env:PG_DSN = 'host=106.52.81.173 user=panda-wiki password=panda123456 dbname=panda-wiki port=5432 sslmode=disable TimeZone=Asia/Shanghai'
$env:REDIS_PASSWORD = 'redis123456'
$env:REDIS_ADDR = '106.52.81.173:6379'
$env:JWT_SECRET = 'jwt_secret_key_panda_wiki_2025'
$env:S3_SECRET_KEY = 'minioadmin'
$env:S3_ACCESS_KEY = 'minioadmin'
$env:S3_ENDPOINT = '106.52.81.173:9000'
$env:QDRANT_API_KEY = 'qdrant_api_key_panda_wiki'
$env:NATS_PASSWORD = 'nats123456'
$env:MQ_NATS_SERVER = 'nats://106.52.81.173:4222'
$env:RAG_CT_RAG_BASE_URL = 'http://106.52.81.173:8080/api/v1'
$env:ADMIN_PASSWORD = 'admin123456'
$env:SUBNET_PREFIX = '169.254.15'
$env:TIMEZONE = 'Asia/Shanghai'
$env:ANYDOC_API_BASE_URL = 'http://panda-wiki-api:8000'
$env:ANYDOC_CRAWLER_BASE_URL = 'http://106.52.81.173:8080'
$env:ANYDOC_UPLOADER_DIR = '/image'
$env:ADMIN_PORT = '2443'
# 设置Caddy管理地址
$env:CADDY_ADMIN_ADDR = 'http://localhost:2019'
# 添加环境变量跳过数据库迁移
$env:SKIP_MIGRATION = 'true'

# 输出启动信息
Write-Host "[启动信息] 正在设置环境变量并启动后端服务..."
Write-Host "[配置检查] 数据库: 106.52.81.173:5432"
Write-Host "[配置检查] Redis: 106.52.81.173:6379"
Write-Host "[配置检查] NATS: nats://106.52.81.173:4222"
Write-Host "[配置检查] RAG服务: http://106.52.81.173:8080/api/v1"
Write-Host ""

# 跳过wire生成步骤（已手动执行）
Write-Host "[准备阶段] Wire依赖注入文件已生成，跳过此步骤..."

# 启动后端服务
Write-Host ""
Write-Host "[启动阶段] 正在启动后端API服务..."
go run ./cmd/api/