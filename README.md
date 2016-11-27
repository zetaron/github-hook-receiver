# GitHub Hook receiver
Starts a server on given `$HOST`.

Listenes for incomming requests on `/deployment` verifying the `$GITHUB_SECRET` and writing the received payload into the given redis database (`$REDIS_DATABASE`)
To do something with these payloads please reffer to [deployment-queue-worker](https://github.com/zetaron/deployment-queue-worker)

When you mount a volume to `/var/cache/secrets` the entrypoint script will create environment variables for each file in there, this way you dont need to expose your secrets to the image.

## Usage
```shell
INGRESS_DNSNAME=github-hook-receiver.zetaron.de docker-compose up -d
```

## Configuration
System Environment Variables:
- **GITHUB_HOOK_RECEIVER_SECRETS_VOLUME_NAME**
- **INGRESS_DNSNAME**

Container Environment Variables:
- **GITHUB_SECRET**
- **REDIS_URL**
- **REDIS_DATABASE**
- **HOST**
