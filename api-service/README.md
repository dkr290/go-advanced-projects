- to custom local reg - tag an existing local image to be pushed to the registry
  docker tag nginx:latest k3d-registry:37801/mynginx:v0.1

- push that image to the registry
  docker push k3d-registry:37801/mynginx:v0.1

- run a pod that uses this image
  kubectl run mynginx --image k3d-registry:37801/mynginx:v0.1
