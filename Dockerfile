FROM scratch

ENV PORT 8080
EXPOSE $PORT

COPY KubeTestPod /KubeTestPod
CMD ["/KubeTestPod"]
