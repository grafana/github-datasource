# runs inside the container so screenshots match
echo "running"

cd ./home

$(npm bin)/cypress install

echo "starting tests"

export CYPRESS_HOST=host.docker.internal
export CYPRESS_BASE_URL=http://host.docker.internal:3000

yarn e2e:update

# this will simulate how it runs on ci
# wget -O grabpl https://grafana-downloads.storage.googleapis.com/grafana-build-pipeline/v0.5.40/grabpl
# chmod +x grabpl
# ./grabpl plugin e2etests
