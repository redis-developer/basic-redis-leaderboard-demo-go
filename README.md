# Basic Redis Leaderboard Demo Golang

Show how the Redis works with Golang.


# How it works?

## How the data is stored:

- The AAPL's details - market cap of 2,6 triillions and USA origin - are stored in a hash like below:
  - E.g `HSET "company:AAPL" symbol "AAPL" market_cap "2600000000000" country USA`
- The Ranks of AAPL of 2,6 trillions are stored in a <a href="https://redislabs.com/ebook/part-1-getting-started/chapter-1-getting-to-know-redis/1-2-what-redis-data-structures-look-like/1-2-5-sorted-sets-in-redis/">ZSET</a>.
  - E.g `ZADD companyLeaderboard 2600000000000 company:AAPL`

## How the data is accessed:

- Top 10 companies:
  - E.g `ZREVRANGE companyLeaderboard 0 9 WITHSCORES`
- All companies:
  - E.g `ZREVRANGE companyLeaderboard 0 -1 WITHSCORES`
- Bottom 10 companies:
  - E.g `ZRANGE companyLeaderboard 0 9 WITHSCORES`
- Between rank 10 and 15:
  - E.g `ZREVRANGE companyLeaderboard 9 14 WITHSCORES`
- Show ranks of AAPL, FB and TSLA:
  - E.g `ZSCORE companyLeaderBoard company:AAPL company:FB company:TSLA`
- Adding market cap to companies:
  - E.g `ZINCRBY companyLeaderBoard 1000000000 "company:FB"`
- Reducing market cap to companies:
  - E.g `ZINCRBY companyLeaderBoard -1000000000 "company:FB"`
- Companies over a Trillion:
  - E.g `ZCOUNT companyLeaderBoard 1000000000000 +inf`
- Companies between 500 billion and 1 trillion:
  - E.g `ZCOUNT companyLeaderBoard 500000000000 1000000000000`

### Code Example: Get top 10 companies

```Go
func (c Controller) Top10() ([]*Company, error) {
    companies, err := c.r.ZRevRange(keyLeaderBoard, 0, 9)
    if err != nil {
        return nil, err
    }
    c.buildCompanies(companies)
    c.buildRanks(companies)
    return companies, nil
}
```

## How to run it locally?

#### Copy `.env.example` to create `.env`. And provide the values for environment variables if needed

- REDIS_HOST: Redis server host
- REDIS_PORT: Redis server port
- REDIS_PASSWORD: Password to the server

### Configure by an environment variable with Redis connection string URL

It is possible to pass any valid Redis URL for Redis options as in [ParseURL Example](https://pkg.go.dev/github.com/go-redis/redis?utm_source=gopls#example-ParseURL)
This way REDIS_HOST, REDIS_PORT, REDIS_PASSWORD are not needed.

- REDIS_URL=redis :// [[username :] password@] host [:port][/database]
- Example REDIS_URL="redis://p%40ssw0rd@redis-16379.hosted.com:16379/0" from [redis-cli, the Redis command line interface](https://redis.io/topics/rediscli)

        Scheme syntax:
          Example: redis://user:secret@localhost:6379/0?foo=bar&qux=baz

          This scheme uses a profile of the RFC 3986 generic URI syntax.
          All URI fields after the scheme are optional.
          The "userinfo" field uses the traditional "user:password" format.

From [Provisional RFC for Redis URIs](https://www.iana.org/assignments/uri-schemes/prov/redis)

### Secure a connection with Redis with a mutual TLS

To support this feature three new environment variables are introduced, TLS_CA_CERT, TLS_CLIENT_CERT, TLS_CLIENT_KEY. They contain paths to respective files in a mounted secret volume. To use it with Kubernetes pods, add this to a container configuration:

```yaml
spec:
  containers:
  - env:
      - name: TLS_CA_CERT
        value: /certs/ca.crt # path to CA certificate
      - name: TLS_CLIENT_CERT
        value: /certs/tls.crt # path to client certificate
      - name: TLS_CLIENT_KEY
        value: /certs/tls.key # path to client key
    image: ghcr.io/denist-huma/basic-redis-leaderboard-demo-go:1.2.2
    name: leaderboard-tls
    ports:
    - containerPort: 8080
      protocol: TCP
    volumeMounts:
    - mountPath: /certs
      name: leaderboard-tls-redis-client-cert

  volumes:
  - name: leaderboard-tls-redis-client-cert
    secret:
      defaultMode: 420
      secretName: leaderboard-tls-redis-client-cert
```

Where the secret "leaderboard-tls-redis-client-cert" has all three files we need. Here is a mere description, not the actual data:

```yaml
apiVersion: v1
data:
  tls.crt:  1261 bytes
  tls.key:  1679 bytes
  ca.crt:   1415 bytes
kind: Secret
metadata:
  name: leaderboard-tls-redis-client-cert
type: kubernetes.io/tls
```

#### Run demo

```sh
go get
go run
```

Follow: http://localhost:8080

## Try it out

<p>
    <a href="https://heroku.com/deploy" target="_blank">
        <img src="https://www.herokucdn.com/deploy/button.svg" alt="Deploy to Heroku" width="200px"/>
    <a>
</p>

<p>
    <a href="https://deploy.cloud.run" target="_blank">
        <img src="https://deploy.cloud.run/button.svg" alt="Run on Google Cloud" width="200px"/>
    </a>

    (See notes: How to run on Google Cloud)

</p>

## How to run on Google Cloud

## 1. Click "Run on Google Cloud"
      
Add the right values as per your infrastructure:
      
 ```
 [ ? ] Value of REDIS_HOST environment variable (Redis server host) <Enter your Redis Host URL>
 [ ? ] Value of REDIS_PORT environment variable (Redis server PORT) <Redis Port>
 [ ? ] Value of REDIS_PASSWORD environment variable (Redis server password) <Redis Password>
 [ ? ] Value of API_PUBLIC_PATH environment variable (Public path to frontend, example `/api/public`) /api/
 [ ? ] Value of IMPORT_PATH environment variable (Path to seed.json file for import, example `/api/seed.json`) seed.json
 [ ? ] Value of API_PORT environment variable (Api public port, example `8080`) 8080
```
      
      
Open up the link under "Manage this application at Cloud Console" to open up "Edit and Deploy New Revision”

## 2. Click “Variables and Secrets”
      
      
![](https://raw.githubusercontent.com/redis-developer/basic-redis-leaderboard-demo-go/master/image4.jpg?v=2&s=2)
      
## 3. Verify the connector
      
![](https://raw.githubusercontent.com/redis-developer/basic-redis-leaderboard-demo-go/master/image2.jpg?v=2&s=2)
      
## 3. Access the app
      
![](https://raw.githubusercontent.com/redis-developer/basic-redis-leaderboard-demo-go/master/image5.jpg?v=2&s=2) 
      
Hence, you should be able to access this app
