FROM scratch
MAINTAINER Paul Rutledge <paul.v.rutledge@gmail.com>
COPY kube-vulcano /kube-vulcano
ENTRYPOINT [ "/kube-vulcano" ]