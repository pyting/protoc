FROM golang:1.12.4

LABEL Description="This is a golang 1.12.4 build env base image, which allows connecting Jenkins agents via JNLP protocols" Vendor="jnlp-salve-golang" Version="1.12.4"

####################install openJDK8
RUN apt-get update && apt-get install -y --no-install-recommends \
    bzip2 \
    unzip \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*


ENV LANG C.UTF-8

RUN { \
    echo '#!/bin/sh'; \
    echo 'set -e'; \
    echo; \
    echo 'dirname "$(dirname "$(readlink -f "$(which javac || which java)")")"'; \
    } > /usr/local/bin/docker-java-home \
    && chmod +x /usr/local/bin/docker-java-home

RUN ln -svT "/usr/lib/jvm/java-8-openjdk-$(dpkg --print-architecture)" /docker-java-home
ENV JAVA_HOME /docker-java-home

RUN set -ex; \
    \
    if [ ! -d /usr/share/man/man1 ]; then \
    mkdir -p /usr/share/man/man1; \
    fi; \
    \
    apt-get update; \
    apt-get install -y --no-install-recommends \
    openjdk-8-jdk-headless="$JAVA_DEBIAN_VERSION" \
    ; \
    rm -rf /var/lib/apt/lists/*; \
    \
    [ "$(readlink -f "$JAVA_HOME")" = "$(docker-java-home)" ]; \
    \
    update-alternatives --get-selections | awk -v home="$(readlink -f "$JAVA_HOME")" 'index($3, home) == 1 { $2 = "manual"; print | "update-alternatives --set-selections" }'; \
    update-alternatives --query java | grep -q 'Status: manual'

####################install jenkins slave
ARG VERSION=3.28
ARG user=jenkins
ARG group=jenkins
ARG uid=1000
ARG gid=1000

ENV HOME /home/${user}
RUN addgroup -g ${gid} ${group}
RUN adduser -h $HOME -u ${uid} -G ${group} -D ${user}

ARG AGENT_WORKDIR=/home/${user}/agent

RUN curl --create-dirs -sSLo /usr/share/jenkins/slave.jar https://repo.jenkins-ci.org/public/org/jenkins-ci/main/remoting/${VERSION}/remoting-${VERSION}.jar \
    && chmod 755 /usr/share/jenkins \
    && chmod 644 /usr/share/jenkins/slave.jar

USER ${user}
ENV AGENT_WORKDIR=${AGENT_WORKDIR}
RUN mkdir /home/${user}/.jenkins && mkdir -p ${AGENT_WORKDIR}

VOLUME /home/${user}/.jenkins
VOLUME ${AGENT_WORKDIR}
WORKDIR /home/${user}

COPY jenkins-slave /usr/local/bin/jenkins-slave

ENTRYPOINT ["jenkins-slave"]