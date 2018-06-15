FROM alpine:3.7
RUN apk update
RUN apk add curl go git openssh libc-dev
RUN curl https://storage.googleapis.com/kubernetes-helm/helm-v2.9.1-linux-amd64.tar.gz | tar -zxvf - 
RUN mkdir -p go/src/app
COPY .ssh/* /root/.ssh/
RUN git clone git@bitbucket.org:bestmile/helm-charts.git 
RUN rm /root/.ssh/id_rsa
WORKDIR /root/go/src/app
RUN pwd
COPY test.go .
RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 8080
CMD ["/root/go/bin/app"]
