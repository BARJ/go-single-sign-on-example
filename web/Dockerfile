FROM node:15-alpine AS builder
WORKDIR /sso
COPY package.json ./
RUN yarn install
COPY . .
RUN yarn build

FROM nginx:1.19-alpine
COPY --from=builder /sso/build /usr/share/nginx/html
EXPOSE 80

