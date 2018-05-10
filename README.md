# tweet-via-searchbar

Tweet via searchbar
https://searchbar.401.jp/

## Build

```bash
make
```

## Run with Docker

```bash
cp .env.sample .env
docker-compose up -d --build
```

## Deploy to Kubernetes

```bash
kubectl apply -f kubernetes.yml

# edit env vars
kubectl -n searchbar edit secret dotenv
```
