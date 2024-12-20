FROM node:20.14.0-buster
WORKDIR /build

COPY package.json package-lock.json ./
RUN npm ci

COPY index.html tsconfig.json vite.config.ts ./
COPY src src

ENTRYPOINT [ "npm", "run", "dev" ]