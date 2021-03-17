## How to run

Install docker https://docs.docker.com/engine/install

Install docker compose https://docs.docker.com/compose/install/

Clone project
```
git clone git@gitlab.com:teamProjects/basic-redis-leaderboard-demo-golang.git
```
Build project to docker images
```
docker-compose build
```
Copy environment .env.example file to .env
```
cp .env.example .env
```

Run backend, front, redis together

```
docker-compose up -d
```
Stop all
```
docker-compose down
```

## Frontend

The `./public` directory contained VueJS build from https://github.com/redis-developer/basic-redis-leaderboard-demo-nodejs

To update it just make new build and place it in `./public` directory
