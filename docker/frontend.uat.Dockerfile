FROM nginx:latest
RUN apt-get update && apt-get install -y curl jq unzip
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
RUN LATEST_TAG=$(curl -s https://api.github.com/repos/VagrantPi/BTM-Admin/releases/latest | jq -r '.tag_name') && \
    curl -L https://github.com/VagrantPi/BTM-Admin/releases/download/$LATEST_TAG/dist-uat.zip -o dist.zip && \
    unzip -o dist.zip
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
