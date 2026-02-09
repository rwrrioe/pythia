
# Quick start

## 1. Clone repo

```
git clone https://github.com/rwrrioe/pythia
```


## 2. Install npm dependencies

```
cd ./frontend
npm i
```

## 3. Set frontend .env

```
cd ./frontend
create .env

ENV (dev)
VITE_API_BASE=http://localhost:8080/api  
VITE_API_ORIGIN=http://localhost:8080  
VITE_WS_PATH=/ws
```

## 4. Set backend .env

```
cd backend/cmd/app
create .env
ENV (dev)
GEMINI_API_KEY=
LOGGER_ENV=local
APP_SECRET=
```


## 5. Build containers

```
cd ./deployments
docker compose -f ./docker-compose.dev.yml build
docker compose -f ./docker-compose.dev.yml up
```

## 6. Run npm

```
cd ./frontend
npm run dev
```
