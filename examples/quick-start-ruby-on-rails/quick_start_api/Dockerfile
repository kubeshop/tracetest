FROM ruby:3.1.2-alpine
WORKDIR /usr/src/app

COPY Gemfile Gemfile.lock .
RUN apk update && apk add make gcc musl-dev && bundle install

COPY . .
RUN chmod 666 ./Gemfile.lock

EXPOSE 8080

CMD [ "bundle", "exec", "rails", "server", "--port", "8080" ]
