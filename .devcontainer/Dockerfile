FROM ghcr.io/armckinney/go:1.25.0

ARG WORKSPACE=/workspaces/verge
ENV WORKSPACE=${WORKSPACE}

# make non-interactive for build
ENV DEBIAN_FRONTEND=noninteractive

RUN mkdir /tmp/container \
    && cd /tmp/container \
    # goreleaser repo
    && echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list \
    # supporting installs
    && apt-get update \
    # primary installs
    && apt-get install --no-install-recommends -y \
        goreleaser \
    # cleanup
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/container

WORKDIR ${WORKSPACE}

# install go dependencies
COPY go.mod go.sum ./
RUN go mod download

# restore dialog for interactive use
ENV DEBIAN_FRONTEND=dialog
