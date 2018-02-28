FROM scatch

COPY demo /bin/demo

EXPOSE 8888

ENTRYPOINT [ "/bin/demo" ]
