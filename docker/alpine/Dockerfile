
FROM alpine


# 替换阿里云的源
RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories && \
    echo "http://mirrors.aliyun.com/alpine/latest-stable/community/" >> /etc/apk/repositories && \
    apk update \
    && apk add --no-cache openssh tzdata shadow libuser sudo\
   	&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
   	&& sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config \
   	&& ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key \
   	&& ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key \
   	&& echo "root:123456" | chpasswd \
    && useradd -m admin \
    && useradd -m user_1 \
   	&& echo "user_1:123456" | chpasswd \
   	&& echo "admin:123456" | chpasswd \
   	&& echo "admin ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers \
    && sed -i 's/^PS1=.*/PS1="[\\u@\\h:\\w]"/' /etc/profile \
    && sed -i "32d" /etc/profile
#    && sed -i 's/^PS1=.*/PS1="\u@\h:\w\$ "/' /etc/profile



#开放端口
EXPOSE 22
## 启动后台运行
CMD ["/usr/sbin/sshd", "-D"]
