# frozen_string_literal: true

RSpec.describe "Rewire::Tracer::RSpec" do
  it "is defined after requiring rewire/tracer/rspec" do
    require "rewire/tracer/rspec"
    expect(defined?(Rewire::Tracer::RSpec)).to be_truthy
  end

  it "responds to configure!" do
    require "rewire/tracer/rspec"
    expect(Rewire::Tracer::RSpec).to respond_to(:configure!)
  end
end
