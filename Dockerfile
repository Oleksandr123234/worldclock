FROM public.ecr.aws/docker/library/alpine:3.22.4

WORKDIR /app/go-go-gin

COPY go-binary .

RUN chmod +x go-binary

USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["./go-binary"]
