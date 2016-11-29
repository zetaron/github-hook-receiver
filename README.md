# GitHub Hook receiver
Starts a server on given `$HOST`.

Listenes for incomming requests on `/deployment` verifying the `$GITHUB_SECRET` and writing the received payload into the given redis database (`$REDIS_DATABASE`)
To do something with these payloads please reffer to [deployment-queue-worker](https://github.com/zetaron/deployment-queue-worker)

When you mount a volume to `/var/cache/secrets` the entrypoint script will create environment variables for each file in there, this way you dont need to expose your secrets to the image.

## Usage

With `docker-compose` (reccomended for development):
```shell
INGRESS_DNSNAME=github-hook-receiver.zetaron.de docker-compose up -d
```

With `docker swarm` (reccomended for production):
```shell
./hooks/deploy
```

## Configuration
The image ships with a [`secret-wrapper`](https://github.com/zetaron/github-hook-receiver/blob/master/secret-wrapper) which allows you to automatically configure your container instance - not the System Environment, except if you are deploying through the `deployment-queue-worker` and the used `WORKER_IMAGE` supports it.

System Environment Variables when using `docker-compose`:
- **GITHUB_HOOK_RECEIVER_SECRETS_VOLUME_NAME**
- **INGRESS_DNSNAME**

System Environment Variables when using `docker swarm`:
- **DEPLOYMENT_ENVIRONMENT** [default=production]
- **DNSNAME** [default=github-hook-receiver.zetaron.de]
- **GITHUB_HOOK_RECEIVER_VERSION** [default=1.0.0]
- **GITHUB_HOOK_RECEIVER_REPLICAS** [default=1]
- **SCHEDULE_ONTO_NODE** [default=cluster-node-1]

Container Environment Variables:
- **GITHUB_SECRET**
- **REDIS_URL**
- **REDIS_DATABASE**
- **HOST**
