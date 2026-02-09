
# Quick start

## 1. Clone repo

```
git clone https://github.com/rwrrioe/pythia
```

## 2. Build docker containers

```
cd ./deployments
docker compose -f ./docker-compose.dev.yml build
docker compose -f ./docker-compose.dev.yml up
```

## 3. Install npm dependencies

```
cd ./frontend
npm i
npm run dev
```

