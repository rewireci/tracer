# frozen_string_literal: true

require_relative "lib/rewire/tracer/version"

Gem::Specification.new do |spec|
  spec.name = "rewire-tracer"
  spec.authors = ["Rewire"]
  spec.version = Rewire::Tracer::VERSION
  spec.summary = "Zero-config OpenTelemetry autoinstrumentation for Ruby CI pipelines"
  spec.description = "Streams traces from RSpec, Rails, and other frameworks to Rewire."
  spec.homepage = "https://github.com/rewireci/tracer"
  spec.license = "MIT"
  spec.metadata["source_code_uri"] = "https://github.com/rewireci/tracer/tree/main/packages/ruby"
  spec.metadata["changelog_uri"]   = "https://github.com/rewireci/tracer/releases"
  spec.required_ruby_version = ">= 3.1"

  spec.files = Dir["lib/**/*"]
  spec.require_paths = ["lib"]

  # OTel gems are optional at runtime — init degrades gracefully if absent
  spec.add_development_dependency "opentelemetry-exporter-otlp", "~> 0.29"
  spec.add_development_dependency "opentelemetry-sdk", "~> 1.6"
  spec.add_development_dependency "rspec", "~> 3.13"
  spec.add_development_dependency "rubocop", "~> 1.70"
  spec.add_development_dependency "rubocop-rspec", "~> 3.4"
end
