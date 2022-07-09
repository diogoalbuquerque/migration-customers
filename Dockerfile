FROM public.ecr.aws/lambda/provided:al2 as build
RUN yum install make golang tar wget -y
RUN go env -w GOPROXY=direct
RUN go install github.com/ibmdb/go_ibm_db/installer@v0.4.1
RUN go run /root/go/pkg/mod/github.com/ibmdb/go_ibm_db@v0.4.1/installer/setup.go
ADD . /APP/
ENV DB2HOME=/clidriver
ENV CGO_CFLAGS=-I$DB2HOME/include
ENV CGO_LDFLAGS=-L$DB2HOME/lib
ENV LD_LIBRARY_PATH=/clidriver/lib
RUN cd /APP && make test
RUN cd /APP && make build

FROM public.ecr.aws/lambda/go
COPY --from=build /clidriver /clidriver
ENV DB2HOME=/clidriver
ENV CGO_CFLAGS=-I$DB2HOME/include
ENV CGO_LDFLAGS=-L$DB2HOME/lib
ENV LD_LIBRARY_PATH=/clidriver/lib
ENV LOG_LEVEL=${LOG_LEVEL:-"info"}
ENV MYSQL_OPEN_CONN_MAX=${MYSQL_OPEN_CONN_MAX:-"100"}
ENV MYSQL_IDLE_CONN_MAX=${MYSQL_IDLE_CONN_MAX:-"2"}
ENV MYSQL_LIFE_CONN_MAX=${MYSQL_LIFE_CONN_MAX:-"10"}
ENV DB2_OPEN_CONN_MAX=${DB2_OPEN_CONN_MAX:-"100"}
ENV DB2_IDLE_CONN_MAX=${DB2_IDLE_CONN_MAX:-"2"}
ENV DB2_LIFE_CONN_MAX=${DB2_LIFE_CONN_MAX:-"10"}
ENV AWS_REGION_NAME=${AWS_REGION_NAME:-"DEFINE_REGION_NAME"}
ENV SECRETS_MANAGER=${SECRETS_MANAGER:-"DEFINE_SECRETS_MANAGER"}
ENV DATASOURCE_LIMIT=${DATASOURCE_LIMIT:-"100"}
COPY --from=build /APP/bin /
ENTRYPOINT ["/MIGRATION_CUSTOMERS"]