FROM node:10.15.0-alpine

COPY package.json package-lock.json ./
COPY src src

RUN npm i

EXPOSE 3001

CMD ["npm", "start"]