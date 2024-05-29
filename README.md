# Inventory

This MS is a simple inventory management system. It is a REST API that allows you to manage the ingredients.

### Start the server

```bash
cp .env.example .env
export $(cat .env | xargs)
docker-compose up
go run main.go
```
