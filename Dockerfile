FROM scratch
COPY krona /krona
ENTRYPOINT ["/krona"]