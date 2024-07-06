IMAGE_NAME=${1:-"golang-alpine-gcc"}
IMAGE_TAG=${2:-"1.21.10"}

if ! docker image inspect $IMAGE_NAME:$IMAGE_TAG > /dev/null 2>&1; then 
    echo "Building image $IMAGE_NAME:$IMAGE_TAG..."
    
    docker build -t $IMAGE_NAME:$IMAGE_TAG -f Dockerfile.base .
else
    echo "Image $IMAGE_NAME:$IMAGE_TAG already exists, skipping build"
fi