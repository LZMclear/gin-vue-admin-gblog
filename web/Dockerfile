# 如果需要用 cicd ，请设置环境变量：
# variables:
#     DOCKER_BUILDKIT: 1

FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
# node:20 不支持 pnpm 11+，必须指定 pnpm 9.x 版本，否则会报错
RUN corepack prepare pnpm@9.14.4 --activate
COPY . /app
WORKDIR /app


FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod

FROM base AS build
COPY --from=prod-deps /app/node_modules /app/node_modules
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install && pnpm run build


FROM nginx:alpine
LABEL MAINTAINER="lz13526112262@163.com"
RUN rm -f /etc/nginx/conf.d/default.conf
COPY --from=base  /app/.docker-compose/nginx/conf.d/my.conf /etc/nginx/conf.d/default.conf
COPY --from=build  /app/dist /usr/share/nginx/html
RUN ls -al /usr/share/nginx/html
