while true; do
    pkill mista
    git pull
    go build
    ./mista
done
