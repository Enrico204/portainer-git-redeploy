# Portainer Git Redeploy

A very small utility for triggering the Git re-deploy for stacks in Portainer
via CLI. It was meant to be used in CD pipelines.

## Usage (OCI/docker container)

A container image is available on Docker Hub: `docker.io/enrico204/portainer-git-redeploy`

## Usage (executable)

Example with environment variables:
```shell
$ export PORTAINER_STACK_ID=1
$ export PORTAINER_URL=https://portainer:9443
$ export PORTAINER_ACCESS_TOKEN=abcdef
$ portainer-git-redeploy
```

Example using CLI options (**note that writing your access token to the CLI is a BAD IDEA**):
```shell
$ portainer-git-redeploy -url https://portainer:9443 -stack-id 1 -access-token abcdef
```

To use a custom certificate file, use `SSL_CERT_FILE`:

```shell
$ export SSL_CERT_FILE=~/certificate.pem
$ portainer-git-redeploy -url https://portainer:9443 -stack-id 1 -access-token abcdef
```

# LICENSE

See [LICENSE](LICENSE).
