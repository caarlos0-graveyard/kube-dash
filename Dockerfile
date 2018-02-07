FROM scratch
COPY kube-dash /kube-dash
ENTRYPOINT ["/kube-dash"]
