# frozen_string_literal: true

RSpec.describe 'Rewire::Tracer.init' do
  INIT_VARS = %w[OTEL_EXPORTER_OTLP_ENDPOINT REWIRE_TOKEN GITHUB_RUN_ID OTEL_SERVICE_NAME GITHUB_REPOSITORY].freeze

  around do |example|
    saved = INIT_VARS.to_h { |k| [k, ENV.delete(k)] }
    example.run
  ensure
    saved.each { |k, v| ENV[k] = v if v }
  end

  context 'when no env vars are set' do
    it 'warns about missing configuration' do
      expect { Rewire::Tracer.init }.to output(/OTEL_EXPORTER_OTLP_ENDPOINT/).to_stderr
    end

    it 'returns a callable' do
      allow(Rewire::Tracer).to receive(:warn)
      stop = Rewire::Tracer.init
      expect(stop).to respond_to(:call)
    end

    it 'returns a fresh callable on each call (not memoized)' do
      allow(Rewire::Tracer).to receive(:warn)
      a = Rewire::Tracer.init
      b = Rewire::Tracer.init
      expect(a).not_to be(b)
    end
  end

  context 'before init is called' do
    it 'shutdown does not raise' do
      expect { Rewire::Tracer.shutdown }.not_to raise_error
    end
  end

  context 'with OTEL_EXPORTER_OTLP_ENDPOINT set' do
    before do
      require 'opentelemetry/sdk'
      require 'opentelemetry/exporter/otlp'
      ENV['OTEL_EXPORTER_OTLP_ENDPOINT'] = 'http://localhost:4318'
    end

    it 'initializes without error and returns a callable' do
      stop = Rewire::Tracer.init
      expect(stop).to respond_to(:call)
      stop.call
    end

    it 'does not emit a warning' do
      expect(Rewire::Tracer).not_to receive(:warn)
      stop = Rewire::Tracer.init
      stop.call
    end

    it 'returns the same callable when called twice' do
      a = Rewire::Tracer.init
      b = Rewire::Tracer.init
      expect(a).to be(b)
      a.call
    end

    it 'configures the exporter with the correct v1/traces path' do
      captured_endpoint = nil
      allow(OpenTelemetry::Exporter::OTLP::Exporter).to receive(:new) do |args|
        captured_endpoint = args[:endpoint]
        OpenTelemetry::SDK::Trace::Export::SimpleSpanProcessor.new(
          OpenTelemetry::SDK::Trace::Export::InMemorySpanExporter.new
        )
      end
      Rewire::Tracer.init
      expect(captured_endpoint).to eq('http://localhost:4318/v1/traces')
    end
  end

  context 'with REWIRE_TOKEN set' do
    before do
      require 'opentelemetry/sdk'
      require 'opentelemetry/exporter/otlp'
      ENV['REWIRE_TOKEN'] = 'rwt_test'
    end

    it 'initializes without error and returns a callable' do
      stop = Rewire::Tracer.init
      expect(stop).to respond_to(:call)
      stop.call
    end

    it 'configures the exporter with the correct otlp/v1/traces path' do
      captured_endpoint = nil
      allow(OpenTelemetry::Exporter::OTLP::Exporter).to receive(:new) do |args|
        captured_endpoint = args[:endpoint]
        OpenTelemetry::SDK::Trace::Export::SimpleSpanProcessor.new(
          OpenTelemetry::SDK::Trace::Export::InMemorySpanExporter.new
        )
      end
      Rewire::Tracer.init
      expect(captured_endpoint).to eq('https://rewireci.com/otlp/v1/traces')
    end
  end
end
