FROM node:14.18-alpine as build

COPY ./frontend .

RUN npm install

RUN npm run build

FROM node:14.18-alpine as production

WORKDIR /app


COPY --from=build . .


EXPOSE 3000

CMD ["./node_modules/.bin/next", "start"]