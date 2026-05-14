# frozen_string_literal: true

require "rails/railtie"
require_relative "../tracer"

module Rewire
  module Tracer
    class Railtie < Rails::Railtie
      initializer "rewire.tracer" do
        if ENV["REWIRE_TOKEN"] || ENV["OTEL_EXPORTER_OTLP_ENDPOINT"]
          Rewire::Tracer.init
        end
      end
    end
  end
end
