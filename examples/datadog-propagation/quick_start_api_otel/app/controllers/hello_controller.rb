require 'uri'
require 'net/http'

class HelloController < ApplicationController

  def local_hello
    service_name = ENV['SERVICE_NAME']

    puts '-----------------------------------------'
    puts 'Headers'
    puts '-----------------------------------------'
    puts request.headers.env.reject { |key| key.to_s.include?('.') }.inspect
    puts '-----------------------------------------'

    render json: { hello: 'world!', service: service_name }
  end

  def remote_hello
    service_name = ENV['SERVICE_NAME']
    remote_hello_api = ENV['REMOTE_API']

    url = URI("http://#{remote_hello_api}:8080/hello")

    http = Net::HTTP.new(url.host, url.port);
    request = Net::HTTP::Get.new(url)

    response = http.request(request)

    # puts '-----------------------------------------'
    # puts 'Response'
    # puts '-----------------------------------------'
    # puts response.body, response.code, response.message, response.each_header.to_h.inspect
    # puts '-----------------------------------------'

    if response.code != '200'
      return render json: { error: response.body }
    end

    parsed_response = JSON.parse(response.body, symbolize_names: true)
    render json: { hello: parsed_response[:hello], from: service_name, to: parsed_response[:service]}
  end
end
