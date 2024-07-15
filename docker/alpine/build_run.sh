#!/bash/sh

docker stop d_alpine_ssh && docker rm d_alpine_ssh && docker rmi i_alpine_ssh

docker build -t i_alpine_ssh ./

docker run -id -p 1022:22 --ip=172.17.0.2  --hostname alpine --name alpine_ssh_1022  i_alpine_ssh

docker run -id -p 1023:22 --ip=172.17.0.3  --hostname alpine --name alpine_ssh_1023  i_alpine_ssh