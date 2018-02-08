# tweet-via-searchbar
Tweet via searchbar

[![Deploy to Docker Cloud](https://files.cloud.docker.com/images/deploy-to-dockercloud.svg)](https://cloud.docker.com/stack/deploy/)

## usage

### install dependencies

```bash
make dep
```

or

```bash
dep ensure
```

### compile

```bash
make
```

### run

```bash
cp .env.sample .env
docker-compose up -d --build
```
