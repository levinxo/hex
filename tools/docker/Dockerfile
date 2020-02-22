FROM nginx:1.17.5

COPY webroot /webroot
COPY nginx.conf /etc/nginx/nginx.conf
COPY http.conf /etc/nginx/http.conf

RUN cp -r /webroot/* /usr/share/nginx/html/ && \
    rm -rf /webroot




