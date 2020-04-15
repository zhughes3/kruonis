# Kruonis

### Running app with Docker Compose

```bash
# start containers
docker-compose --env-file cmd/timelines/config.env up

# stop containers
docker-compose --env-file cmd/timelines/config.env down

# remove unused containers
docker system prune
```
