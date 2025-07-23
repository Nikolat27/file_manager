# 1
FROM node:lts-alpine AS build-stage

RUN npm install -g http-server

WORKDIR /app

COPY package*.json ./

RUN npm install 

COPY . . 

RUN npm run build

# 2
FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html*

COPY --from=build-stage /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
