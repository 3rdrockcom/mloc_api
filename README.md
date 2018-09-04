# MLOC

## Installation
```bash
# Get migration utility
# https://github.com/golang-migrate/migrate/releases
wget https://github.com/golang-migrate/migrate/releases/download/v3.4.0/migrate.linux-amd64.tar.gz -O - | tar xz && mv migrate.linux-amd64 migrate;

# Do migration
./migrate -source file://migrations -database mysql://username:password@localhost/mloc?parseTime=true up;
```

## Usage
```bash
# Build and run application server (development)
make clean && make build && make run-development OPTIONS=serve

# Run application server (development)
godotenv -f .env.development ./mloc serve;
```