# Oauth 2.0 example with golang, postgres and memcached

![](https://user-images.githubusercontent.com/61962654/165940093-3998bf66-40e5-47ab-bb8f-a1db112d3df4.png)
![](https://user-images.githubusercontent.com/61962654/165940087-484fac3f-87a3-4769-b225-62121154b136.png)

### Default Endpoint `:8080`

---

## Requirements

* **Go:** `1.18`
* **PostgreSQL:** `14.2`
* **Memcached:** `1.6.15`

## docker-compose

No additional settings are required, just start and use

```shell
docker-compose up
```


### Dockerfile

```shell
docker build . --tag=oauth
docker run oauth
```

## Building / Running

```shell
go build -o server
./server
```

### Custom paths

```shell
server -config config.yaml -log logs.txt
```

### Load config from environment variables

```shell
HOST_PORT=8181 server -env
```

All variables are listed in [config struct](https://github.com/illiafox/go-oauth/blob/master/utils/config/struct.go)

## Prometheus metrics: `/metrics`



## PosgresSQL

```shell
psql sql/migrate-up.sql
```
### Tables

```sql
users
(
    user_id  bigint PRIMARY KEY,

    username varchar(128) UNIQUE    NOT NULL,

    token    char(40)     UNIQUE    NOT NULL
);
```

```sql
sessions
(
    token   char(128) PRIMARY KEY,

    user_id bigint    NOT NULL,

    created date,
    
    FOREIGN KEY (user_id)
    REFERENCES users (user_id)
    
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
```