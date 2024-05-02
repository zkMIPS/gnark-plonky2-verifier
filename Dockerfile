# base image
FROM --platform=linux/amd64 ubuntu:latest

# install necessary packages
RUN apt-get update && apt-get install -y \
    curl \
    gnupg \
    lsb-release \
    apt-transport-https \
    ca-certificates

# add Golang offical GPG key
RUN curl -sSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -

RUN apt-get install wget

RUN wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz

RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz

# setup Golang env
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# install MySQL
RUN apt-get update && apt-get install -y mysql-server && apt-get install -y mysql-client

# set up MySQL root password
ENV MYSQL_ROOT_PASSWORD 123456

COPY my.cnf /etc/mysql/mysql.conf.d/mysqld.cnf
COPY    storage/migrations/db.sql /SQL
RUN     mysqld_safe & until mysqladmin ping; do sleep 1; done && \
        mysql -e "SOURCE /SQL;" && mysqladmin -u root password 123456
RUN service mysql restart

COPY . .
ENV GOOS linux
ENV GOARCH amd64
RUN GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/snark_server ./server
RUN GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/gnark_sol_caller gnark_sol_caller.go
RUN mkdir -p /app/server
RUN touch /app/server/server.log
RUN chmod a+x /usr/local/bin/snark_server
COPY start.sh /usr/local/bin/start.sh
RUN chmod a+x /usr/local/bin/start.sh

# expose mysql,server port
EXPOSE 3306 50051

CMD source /usr/local/bin/start.sh

# Set the command to be executed when the container starts
ENTRYPOINT ["/usr/local/bin/start.sh"]