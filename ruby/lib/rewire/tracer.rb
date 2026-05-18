# frozen_string_literal: true

require_relative 'tracer/version'
require_relative 'tracer/detect_ci'

module Rewire
  module Tracer
    REWIRE_ENDPOINT = 'https://rewireci.com/otlp/v1'

    @shutdown_fn = nil

    def self.init
      return @shutdown_fn if @shutdown_fn

      endpoint = ENV.fetch('OTEL_EXPORTER_OTLP_ENDPOINT', nil)
      token = ENV.fetch('REWIRE_TOKEN', nil)

      unless endpoint || token
        warn '[rewire] Neither OTEL_EXPORTER_OTLP_ENDPOINT nor REWIRE_TOKEN is set — tracing disabled'
        return -> {}
      end

      begin
        require 'opentelemetry/sdk'
        require 'opentelemetry/exporter/otlp'
      rescue LoadError
        warn '[rewire] opentelemetry-sdk is not installed — tracing disabled'
        return -> {}
      end

      begin
        ci = detect_ci
        trace_url = endpoint ? "#{endpoint.chomp('/')}/v1/traces" : "#{REWIRE_ENDPOINT}/traces"
        headers = {}
        headers['Authorization'] = "Bearer #{token}" if !endpoint && token

        resource_attrs = {
          'ci.platform' => ci.platform,
          'service.name' => ENV.fetch('OTEL_SERVICE_NAME') { ENV.fetch('GITHUB_REPOSITORY', 'unknown') }
        }
        resource_attrs['run.id'] = ci.run_id if ci.run_id

        OpenTelemetry::SDK.configure do |c|
          c.resource = OpenTelemetry::SDK::Resources::Resource.create(resource_attrs)
          c.add_span_processor(
            OpenTelemetry::SDK::Trace::Export::BatchSpanProcessor.new(
              OpenTelemetry::Exporter::OTLP::Exporter.new(
                endpoint: trace_url,
                headers: headers
              )
            )
          )
        end

        @shutdown_fn = -> { OpenTelemetry.tracer_provider.shutdown }
        @shutdown_fn
      rescue StandardError => e
        warn "[rewire] Failed to initialize Rewire Tracer: #{e.message}"
        -> {}
      end
    end

    def self.shutdown
      @shutdown_fn&.call
    end

    def self._reset
      @shutdown_fn = nil
    end
  end
end

begin
  require 'rails/railtie'
  require_relative 'tracer/railtie'
rescue LoadError
  # Rails is not available; skip Railtie auto-initialization.
end
