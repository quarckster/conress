FROM docker.io/library/golang:1.17 AS builder

COPY startcypress.go .

RUN go build startcypress.go

FROM quay.io/redhatqe/cypress-base:latest
LABEL maintainer="dmisharo@redhat.com"

# Firefox releases
# https://download-installer.cdn.mozilla.net/pub/firefox/releases/
ARG FIREFOX_VERSION="91.0.1esr"
# Chrome versions
# https://www.ubuntuupdates.org/package/google_chrome/stable/main/base/google-chrome-stable
ARG CHROME_VERSION="93.0.4577.63"

ARG CYPRESS_VERSION="8.4.0"

ENV VNC_PORT="5999" \
    DISPLAY=":99" \
    DBUS_SESSION_BUS_ADDRESS="/dev/null" \
    HOME="${CYPRESS_HOME}" \
    VNC_GEOMETRY=${VNC_GEOMETRY:-"1600x900"} \
    CYPRESS_INSTALL_BINARY="0" \
    CYPRESS_RUN_BINARY="${CYPRESS_HOME}/Cypress/Cypress" \
    PATH="/opt/firefox:/opt/google/chrome:${PATH}"

EXPOSE ${VNC_PORT}

RUN curl -LO https://cdn.cypress.io/desktop/${CYPRESS_VERSION}/linux-x64/cypress.zip && \
    unzip cypress.zip && \
    rm -f cypress.zip

RUN mkdir -p .cache/dconf .mozilla/plugins .vnc/ && \
    touch .Xauthority .vnc/config

RUN curl -LO https://dl.google.com/linux/chrome/rpm/stable/x86_64/google-chrome-stable-${CHROME_VERSION}-1.x86_64.rpm && \
    rpm -i google-chrome-stable-${CHROME_VERSION}-1.x86_64.rpm && \
    rm -f google-chrome-stable-${CHROME_VERSION}-1.x86_64.rpm

RUN curl -LO https://download-installer.cdn.mozilla.net/pub/firefox/releases/${FIREFOX_VERSION}/linux-x86_64/en-US/firefox-${FIREFOX_VERSION}.tar.bz2 && \
    tar -C /opt -xjvf firefox-${FIREFOX_VERSION}.tar.bz2 && \
    rm -f firefox-${FIREFOX_VERSION}.tar.bz2

COPY --from=builder /go/startcypress /usr/bin/startcypress

RUN chmod +x /usr/bin/startcypress && \
    chgrp -R 0 ${CYPRESS_HOME} && \
    chmod -R g=u ${CYPRESS_HOME}

USER 1001
