
# build environment
FROM node:17 as build
WORKDIR /app

RUN yarn

COPY . .

RUN yarn install

RUN yarn run build

# production environment
FROM nginxinc/nginx-unprivileged:1.15
COPY --from=build  /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 5000

CMD ["nginx", "-g", "daemon off;"]
