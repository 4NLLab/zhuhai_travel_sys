FROM node:22-alpine AS build

WORKDIR /app

RUN corepack enable

COPY miniprogram/package.json miniprogram/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY miniprogram/ ./

ARG VITE_API_MODE=local
ARG VITE_API_BASE_URL=/api/v1
ARG VITE_ALLOW_FALLBACK=true
ARG VITE_MOCK_SCENARIO=phase2-success

ENV VITE_API_MODE=${VITE_API_MODE}
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV VITE_ALLOW_FALLBACK=${VITE_ALLOW_FALLBACK}
ENV VITE_MOCK_SCENARIO=${VITE_MOCK_SCENARIO}

RUN pnpm build:h5

FROM nginx:1.27-alpine

COPY docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist/build/h5 /usr/share/nginx/html

EXPOSE 80

