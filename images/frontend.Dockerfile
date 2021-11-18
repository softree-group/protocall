FROM node:16 as builder

RUN git clone https://github.com/softree-group/protocall-front \
    && cd protocall-front \
    && npm ci \
    && npm run build

FROM nginx:latest

COPY --from=builder /protocall-front/build /var/www/html
