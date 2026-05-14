# frozen_string_literal: true

require "rewire/tracer"

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups
  config.order = :random
  config.warnings = true

  config.before { Rewire::Tracer._reset }
  config.after  { Rewire::Tracer._reset }
end
