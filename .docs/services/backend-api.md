## Create Migration 

``` docker compose exec api atlas migrate diff initial_migration --env local ```

Replace "initial_migration" with new name for migration

## Apply Migration 

``` docker compose exec api atlas migrate apply --env local ```

## Run Tests
```  docker compose exec -T api go test ./internal/handlers  ```