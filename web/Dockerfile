FROM node:15-alpine

WORKDIR /sso

COPY package.json ./
RUN yarn install

COPY . .

ENV PORT=80

ENTRYPOINT yarn start
EXPOSE 80
