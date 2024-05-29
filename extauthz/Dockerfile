FROM scratch

COPY build/envoy-extauthz /envoy-extauthz

EXPOSE 9002

RUN ls -la

ENTRYPOINT ["/envoy-extauthz"]
