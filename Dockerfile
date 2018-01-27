FROM scratch
COPY main /krona
ENTRYPOINT ["/krona"]