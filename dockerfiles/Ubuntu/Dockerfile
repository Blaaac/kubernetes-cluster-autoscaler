FROM golang:1.17

WORKDIR /root
RUN \
  git clone https://github.com/Blaaac/kubernetes-cluster-autoscaler.git && \
  cd kubernetes-cluster-autoscaler && \
  go build -o autoscaler  cmd/main.go


FROM golang:1.17
WORKDIR /go/bin
COPY --from=0 /root/kubernetes-cluster-autoscaler/autoscaler .
CMD ["./autoscaler"]
