require 'net/http'
require 'ddtrace'

Datadog.configure do |c|
  c.service = ENV['SERVICE_NAME']
  c.logger.level = ::Logger::ERROR

  c.tracing.instrument :rails
  c.tracing.instrument :http

  # List of header formats that should be extracted
  c.tracing.distributed_tracing.propagation_extract_style = [ 'tracecontext' ]

  # List of header formats that should be injected
  c.tracing.distributed_tracing.propagation_inject_style = [ 'tracecontext' ]
  c.tracing.distributed_tracing.propagation_style = [ 'tracecontext' ]

end

module Datadog
  module Tracing
    module Transport
      class SerializableTrace
        def to_msgpack(packer = nil)
          if ENV.has_key?('INJECT_UPPER_TRACE_ID')
            return trace.spans.map { |s| SerializableSpan.new(s) }.to_msgpack(packer)
          end

          upper_trace_id = trace.spans.find { |span| span.meta.has_key?('_dd.p.tid') }.meta['_dd.p.tid']
          trace.spans.each do |span|
            span.meta["propagation.upper_trace_id"] = upper_trace_id
          end

          trace.spans.map { |s| SerializableSpan.new(s) }.to_msgpack(packer)
        end
      end
    end
  end
end
