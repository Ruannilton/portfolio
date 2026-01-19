
# Inicia o Docker Desktop se não estiver rodando
if (-not (Get-Process -Name "Docker Desktop" -ErrorAction SilentlyContinue)) {
	Write-Host "Iniciando Docker Desktop..."
	Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
}

# Aguarda o Docker Engine estar disponível (timeout de 2 minutos)
$timeout = 120
$elapsed = 0
while ($true) {
	try {
		docker info | Out-Null
		break
	} catch {
		if ($elapsed -ge $timeout) {
			Write-Error "Timeout ao aguardar o Docker iniciar."
			exit 1
		}
		Start-Sleep -Seconds 2
		$elapsed += 2
	}
}

docker compose --env-file .dev.env -f docker-compose.dev.yaml down
docker compose --env-file .dev.env -f docker-compose.dev.yaml up --build -d
   