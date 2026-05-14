# frozen_string_literal: true

require_relative '../tracer'

module Rewire
  module Tracer
    # Provides RSpec integration. Call configure! from spec_helper.rb:
    #
    #   require "rewire/tracer/rspec"
    #   Rewire::Tracer::RSpec.configure!
    #
    # This adds before/after suite hooks for init/shutdown and an around
    # hook that wraps each example in an OTel span with pass/fail status.
    module RSpec
      def self.configure!
        ::RSpec.configure do |config|
          config.before(:suite) do
            Rewire::Tracer.init
            begin
              require 'opentelemetry'
              @tracer = OpenTelemetry.tracer_provider.tracer('rewire.rspec')
            rescue LoadError
              @tracer = nil
            end
          end
          config.after(:suite) { Rewire::Tracer.shutdown }

          config.around(:each) do |example|
            tracer = @tracer
            if tracer
              tracer.in_span(
                example.full_description,
                attributes: {
                  'test.name' => example.description,
                  'test.file' => example.metadata[:file_path].to_s
                }
              ) do |span|
                example.run
                case example.execution_result.status
                when :passed
                  span.status = OpenTelemetry::Trace::Status.ok
                  span.set_attribute('test.status', 'passed')
                when :failed
                  span.status = OpenTelemetry::Trace::Status.error
                  span.set_attribute('test.status', 'failed')
                when :pending
                  span.set_attribute('test.status', 'skipped')
                end
              end
            else
              example.run
            end
          rescue LoadError, NameError
            example.run
          end
        end
      end
    end
  end
end
