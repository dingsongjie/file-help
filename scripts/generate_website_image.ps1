$tag = $args[0]
docker build . -f ./build/dockerfile -t registry.powertradepro.com/container_group/images/file-helper-api:$tag --progress=plain