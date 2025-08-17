# setup-lisa.ps1
# Script para crear la estructura completa del proyecto Lisa

Write-Host "INICIANDO: Creando estructura del proyecto Lisa" -ForegroundColor Green

# Verificar que estamos en el directorio correcto
if (-not (Test-Path "go.mod")) {
    Write-Host "ERROR: No se encuentra go.mod. Asegurate de estar en la carpeta raiz de Lisa" -ForegroundColor Red
    exit 1
}

# Verificar que el módulo sea Lisa
$goModContent = Get-Content "go.mod" -Raw
if ($goModContent -notmatch "module\s+Lisa") {
    Write-Host "ERROR: El modulo no es 'Lisa'. Contenido actual de go.mod:" -ForegroundColor Red
    Write-Host $goModContent -ForegroundColor Yellow
    exit 1
}

Write-Host "OK: Modulo Lisa encontrado correctamente" -ForegroundColor Green

Write-Host "CREANDO: Directorios..." -ForegroundColor Yellow

# Crear estructura de directorios
$directories = @(
    "cmd/bot",
    "internal/bot",
    "internal/whatsapp",
    "internal/jira",
    "internal/ai",
    "internal/config",
    "internal/database",
    "internal/middleware",
    "internal/utils",
    "pkg/types",
    "pkg/errors",
    "web/static",
    "web/templates", 
    "web/dashboard",
    "scripts",
    "configs",
    "docs",
    "tests/integration",
    "tests/unit",
    "tests/mocks",
    "deployments/docker",
    "deployments/kubernetes",
    "deployments/systemd",
    "migrations/whatsmeow"
)

foreach ($dir in $directories) {
    New-Item -ItemType Directory -Path $dir -Force | Out-Null
    Write-Host "  OK: $dir" -ForegroundColor Gray
}

Write-Host "CREANDO: Archivos principales..." -ForegroundColor Yellow

# Crear archivos vacíos con estructura básica
$files = @{
    # Entrada principal
    "cmd/bot/main.go" = ""
    
    # Bot Discord
    "internal/bot/discord.go" = ""
    "internal/bot/commands.go" = ""
    "internal/bot/handlers.go" = ""
    
    # WhatsApp
    "internal/whatsapp/client.go" = ""
    "internal/whatsapp/handlers.go" = ""
    "internal/whatsapp/qr.go" = ""
    "internal/whatsapp/session.go" = ""
    "internal/whatsapp/parser.go" = ""
    
    # Jira
    "internal/jira/client.go" = ""
    "internal/jira/issues.go" = ""
    "internal/jira/projects.go" = ""
    "internal/jira/workflows.go" = ""
    "internal/jira/types.go" = ""
    
    # AI (Gemini)
    "internal/ai/gemini.go" = ""
    "internal/ai/analyzer.go" = ""
    "internal/ai/suggestions.go" = ""
    "internal/ai/classifier.go" = ""
    "internal/ai/prompts.go" = ""
    
    # Configuración
    "internal/config/config.go" = ""
    
    # Base de datos
    "internal/database/models.go" = ""
    "internal/database/migrations.go" = ""
    "internal/database/repository.go" = ""
    
    # Middleware
    "internal/middleware/auth.go" = ""
    "internal/middleware/logging.go" = ""
    "internal/middleware/ratelimit.go" = ""
    
    # Utilidades
    "internal/utils/logger.go" = ""
    "internal/utils/validator.go" = ""
    "internal/utils/helpers.go" = ""
    
    # Tipos
    "pkg/types/discord.go" = ""
    "pkg/types/whatsapp.go" = ""
    "pkg/types/jira.go" = ""
    "pkg/types/gemini.go" = ""
    "pkg/types/common.go" = ""
    
    # Errores
    "pkg/errors/errors.go" = ""
    
    # Configuraciones
    "configs/config.yaml" = ""
    "configs/config.dev.yaml" = ""
    "configs/config.prod.yaml" = ""
    
    # Scripts
    "scripts/deploy.sh" = ""
    "scripts/setup.sh" = ""
    "scripts/migrate.sh" = ""
    
    # Docker
    "deployments/docker/Dockerfile" = ""
    "deployments/docker/docker-compose.yml" = ""
    "deployments/docker/postgres-init.sql" = ""
    
    # Kubernetes
    "deployments/kubernetes/deployment.yaml" = ""
    "deployments/kubernetes/service.yaml" = ""
    "deployments/kubernetes/configmap.yaml" = ""
    "deployments/kubernetes/postgres.yaml" = ""
    
    # Systemd
    "deployments/systemd/lisa-bot.service" = ""
    
    # Migraciones
    "migrations/whatsmeow/001_initial.sql" = ""
    "migrations/whatsmeow/002_sessions.sql" = ""
    
    # Documentación
    "docs/API.md" = ""
    "docs/SETUP.md" = ""
    "docs/COMMANDS.md" = ""
    
    # Tests
    "tests/integration/bot_test.go" = ""
    "tests/unit/config_test.go" = ""
    "tests/mocks/whatsapp_mock.go" = ""
    
    # Archivos raíz
    ".env.example" = ""
    ".gitignore" = ""
    "Makefile" = ""
    "README.md" = ""
}

foreach ($file in $files.Keys) {
    $null = New-Item -ItemType File -Path $file -Force
    Write-Host "  OK: $file" -ForegroundColor Gray
}

Write-Host ""
Write-Host "EXITO: Estructura del proyecto Lisa creada exitosamente!" -ForegroundColor Green
Write-Host ""
Write-Host "PROXIMOS PASOS:" -ForegroundColor Cyan
Write-Host "  1. Configurar .env.example con tus variables" -ForegroundColor White
Write-Host "  2. Instalar dependencias: go mod tidy" -ForegroundColor White
Write-Host "  3. Configurar base de datos PostgreSQL" -ForegroundColor White
Write-Host "  4. Implementar logica en los archivos creados" -ForegroundColor White
Write-Host ""
Write-Host "ESTRUCTURA CREADA EN: $(Get-Location)" -ForegroundColor Yellow