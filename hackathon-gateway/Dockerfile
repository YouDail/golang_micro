FROM alpine:3.9

RUN apk update \
  && apk add --no-cache tzdata \
  && \cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD config.yaml hackathon-gateway service.json /

EXPOSE 10010 10112
CMD ["/hackathon-gateway"]

