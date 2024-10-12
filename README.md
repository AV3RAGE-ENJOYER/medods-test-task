# MEDODS Go Test task

**[Пример конфигурационного файла config.env](config.env)**

## Тестирование

```bash
go test -v ./tests
```

## Установка

### Docker
```bash
git clone https://github.com/AV3RAGE-ENJOYER/medods-test-task
cd medods-test-task
docker compose up -d
```

### Source

Для корректной работы приложения нужен работающий сервер **PostgreSQL** на **localhost**. Также нужно изменить хост в переменной **POSTGRES_URL** в файле [config.env](config.env)  

```bash
git clone https://github.com/AV3RAGE-ENJOYER/medods-test-task
cd medods-test-task
go mod download && go mod verify
go build main.go
./main
```