
### 打包流程

1. 执行build.sh，制作镜像
2. 使用docker跑起来 `docker run -d --restart -p 0.0.0.0:80:80 -v [the tools path]:/usr/share/nginx/html/tools levinxo/website:latest`


