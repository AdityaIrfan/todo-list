## Preparation Before Running
1. Download All Dependent Packages  
   `go mod tidy`
2. Set Up Postgres Database Name, Username & Password  
   `edit file config.yaml & migration.yaml`
3. Migrate database  
   `sql-migrate up -config=migration.yaml -env="local"`
4. Main File Location  
   `cmd/api/main.go`