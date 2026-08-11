FROM gcr.io/distroless/static
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/rutter /usr/local/bin/rutter
ENTRYPOINT ["/usr/local/bin/rutter"]
