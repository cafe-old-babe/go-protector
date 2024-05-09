# bash/sh
docker stop d_alpine_ssh && docker rm d_alpine_ssh && docker rmi i_alpine_ssh

docker build -t i_alpine_ssh ./

docker run -id -p 1022:22 --hostname alpine --name d_alpine_ssh  i_alpine_ssh