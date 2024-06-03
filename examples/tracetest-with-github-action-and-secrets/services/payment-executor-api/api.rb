require "sinatra"
require "json"
require "uri"
require "net/http"

require "opentelemetry/sdk"
require "opentelemetry/exporter/otlp"
require "opentelemetry/instrumentation/all"

require "./config"

config = current_config()

set :port, config.port
set :bind, '0.0.0.0'

OpenTelemetry::SDK.configure do |c|
  otel_exporter = OpenTelemetry::Exporter::OTLP::Exporter.new(endpoint: config.otel_exporter_http_url)
  processor = OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(otel_exporter)

  c.service_name = config.otel_service_name
  c.add_span_processor(processor)

  c.use_all()
end

error do
  OpenTelemetry::Trace.current_span.record_exception(env['sinatra.error'])
end

post '/payment/execute' do
  content_type :json # should return json

  payment_data = JSON.parse(request.body.read)

  amount = payment_data["amount"]
  age = payment_data["age"]

  if amount < 10000
    # don't need to analyze risk
    execute_payment(amount)
    return { status: "executed" }.to_json
  end

  score = call_risk_api(amount, age)

  if score < 0
    raise "This case should not be happening"
  end

  if score > 50000
    return { status: "denied" }.to_json
  end

  execute_payment(amount)
  return { status: "executed" }.to_json
end

def call_risk_api(amount, age)
  config = current_config()
  url = URI(config.risk_analysis_url)

  http = Net::HTTP.new(url.host, url.port);
  request = Net::HTTP::Post.new(url)
  request["Content-Type"] = "application/json"
  request.body = JSON.dump({
    "amount": amount,
    "age": age
  })

  response = http.request(request)
  output = JSON.parse(response.read_body)

  output["score"]
end

def execute_payment(amount)
  tracer = OpenTelemetry.tracer_provider.tracer('tracer')

  tracer.in_span("execute_payment", attributes: { "amount" => amount }) do |span|
    # simulate payment being execuded
    sleep(0.05) # 50 milliseconds
  end
end
