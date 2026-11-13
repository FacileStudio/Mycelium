# The encre release the artifact validator comes from. Bump it with the
# @facile/encre pin in apps/client/package.json so the server and the viewer
# agree on the format.
ARG ENCRE_VERSION=0.10.0

FROM oven/bun:1 AS client-build
WORKDIR /client
COPY apps/client/package.json apps/client/bun.lock* ./
RUN bun install --frozen-lockfile
COPY apps/client/ .
RUN bun run build

# encre's release ships the `encre` CLI as a static binary. It is fetched for
# the image's architecture and checked against the release's checksums.txt,
# so the distroless image runs it with no libc. Nothing here depends on this
# repository: the layer is cached until ENCRE_VERSION moves.
FROM alpine:3.21 AS encre-fetch
ARG ENCRE_VERSION
ARG TARGETARCH
RUN set -eu; apk add --no-cache curl; \
    base="https://github.com/FacileStudio/encre/releases/download/v${ENCRE_VERSION}"; \
    archive="encre_${ENCRE_VERSION}_linux_${TARGETARCH:-amd64}.tar.gz"; cd /tmp; \
    curl -fsSLO "$base/$archive" && curl -fsSLO "$base/checksums.txt"; \
    grep " $archive\$" checksums.txt | sha256sum -c -; \
    tar -xzf "$archive" -C /usr/local/bin encre; rm "$archive" checksums.txt; \
    /usr/local/bin/encre --version

FROM golang:1.26.11-alpine AS api-build

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /repo

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# The version is the VERSION file, which the release commit bumps. It is not read
# from git tags (Dokploy clones carry none) or from the Dokploy env, so the
# container always reports the release it was built from. Goreleaser stamps the
# same variable from the tag, without the leading v.
RUN VERSION="$(tr -d '[:space:]' < VERSION)" && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -buildvcs=false \
    -ldflags="-s -w -X github.com/FacileStudio/Mycelium/cmd.version=${VERSION#v}" \
    -o bin/mycelium .

# Not the :nonroot variant: the data directory is a named volume, which Docker
# creates owned by root, so a non-root process could not write to it.
FROM gcr.io/distroless/static-debian12

COPY --from=api-build /repo/bin/mycelium /mycelium
COPY --from=client-build /client/build /client
COPY --from=encre-fetch /usr/local/bin/encre /usr/local/bin/encre

# The validator the artifact routes run. Set explicitly: distroless has no
# shell to resolve PATH surprises in.
ENV ENCRE_BIN=/usr/local/bin/encre

# Obligatory: a relative ./client would resolve against the image's working
# directory, and the SPA would silently not be served — the API answers, the
# healthcheck is green, and only a browser sees the failure.
ENV CLIENT_DIR=/client

EXPOSE 8420

ENTRYPOINT ["/mycelium"]
CMD ["serve"]
