BASENAME=portainer-git-redeploy
VERSION=5
ifeq (${IMAGE_PATH},)
IMAGE_PATH=docker.io/enrico204/portainer-git-redeploy
endif

image:
	buildah manifest rm containers-storage:${IMAGE_PATH}:${VERSION} > /dev/null 2>&1 || true
	buildah manifest create ${IMAGE_PATH}:${VERSION}

	buildah bud --manifest ${IMAGE_PATH}:${VERSION} -f Dockerfile --arch amd64
	buildah bud --manifest ${IMAGE_PATH}:${VERSION} -f Dockerfile --arch arm64 --variant v8

image-distro:
	buildah manifest rm containers-storage:${IMAGE_PATH}:${VERSION}-debian > /dev/null 2>&1 || true
	buildah manifest create ${IMAGE_PATH}:${VERSION}-debian

	buildah bud --manifest ${IMAGE_PATH}:${VERSION}-debian -f Dockerfile.debian --arch amd64
	buildah bud --manifest ${IMAGE_PATH}:${VERSION}-debian -f Dockerfile.debian --arch arm64 --variant v8

push:
	buildah manifest push --all --format=docker ${IMAGE_PATH}:${VERSION} docker://${IMAGE_PATH}:${VERSION}

push-distro:
	buildah manifest push --all --format=docker ${IMAGE_PATH}:${VERSION}-debian docker://${IMAGE_PATH}:${VERSION}-debian

inspect:
	buildah manifest inspect ${IMAGE_PATH}:${VERSION}
