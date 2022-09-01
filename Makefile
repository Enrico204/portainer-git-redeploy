BUILDAH_ARGS=--layers=true -f Dockerfile
BASENAME=portainer-git-redeploy
VERSION=1
ifeq (${IMAGE_PATH},)
IMAGE_PATH=docker.io/enrico204/portainer-git-redeploy
endif

image:
	buildah bud ${BUILDAH_ARGS} --arch amd64 -t ${BASENAME}-amd64:${VERSION}
	buildah bud ${BUILDAH_ARGS} --arch arm64 --variant v8 -t ${BASENAME}-arm64:${VERSION}

	buildah manifest inspect ${IMAGE_PATH}:${VERSION} | jq -r '.manifests | .[].digest' | xargs -n 1 buildah manifest remove ${IMAGE_PATH}:${VERSION} || true
	buildah manifest rm ${IMAGE_PATH}:${VERSION} || true

	buildah manifest create ${IMAGE_PATH}:${VERSION}
	buildah manifest add ${IMAGE_PATH}:${VERSION} ${BASENAME}-amd64:${VERSION}
	buildah manifest add ${IMAGE_PATH}:${VERSION} ${BASENAME}-arm64:${VERSION}

push:
	buildah manifest push --all --format=docker ${IMAGE_PATH}:${VERSION} docker://${IMAGE_PATH}:${VERSION}

inspect:
	buildah manifest inspect ${IMAGE_PATH}:${VERSION}
