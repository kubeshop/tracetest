require "sinatra"

require "opentelemetry/sdk"
require "opentelemetry/exporter/otlp"
require "opentelemetry/instrumentation/sinatra"

set :port, 8080

OpenTelemetry::SDK.configure do |c|
  c.use "OpenTelemetry::Instrumentation::Sinatra"
end

error do
  OpenTelemetry::Trace.current_span.record_exception(env['sinatra.error'])
end

get '/hello' do
  'Hello World'
end
