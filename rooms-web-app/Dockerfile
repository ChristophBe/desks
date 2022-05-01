
# build environment
FROM node:16 as build
WORKDIR /app

RUN yarn

COPY package.json package.json

RUN yarn install

COPY babel.config.js babel.config.js
COPY jsconfig.json jsconfig.json
COPY tsconfig.json tsconfig.json
COPY vue.config.js vue.config.js
COPY public public
COPY src src


RUN yarn run build

# production environment
FROM nginxinc/nginx-unprivileged:1.15
COPY --from=build  /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 5000

CMD ["nginx", "-g", "daemon off;"]
