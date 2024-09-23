# Rory Pearson

### Test

`docker run -p 3000:3000 -e SERVER_PORT=3000 -e UI_BUILD_PATH='./ui/build' -e GIN_MODE='release' -it rory-pearson-test`

### Build

```
chmod +x cmd/docker/build.sh

./cmd/docker/build.sh <version>
```
