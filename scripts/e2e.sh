# runs e2e in a container

docker stop plugin-e2e-test
docker rm plugin-e2e-test

cd ../

mage build:linux

ip=$(ipconfig getifaddr en0)

echo "$ip"

rm -r ./provisioning
cp -r ../plugin-provisioning/provisioning ./provisioning

# build the build-pipeline inside the plugin e2e container and run "plugin e2etests" sub command
docker run -it \
  --add-host=localhost:"$ip" \
  --add-host=127.0.0.1:"$ip" \
  --network="host" \
  -v $(pwd):/home \
  --name=plugin-e2e-test grafana/grafana-plugin-ci-e2e:latest \
  sh "/home/scripts/test.sh"