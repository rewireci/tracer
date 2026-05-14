# frozen_string_literal: true

RSpec.describe "Rewire::Tracer.detect_ci" do
  CI_VARS = %w[GITHUB_RUN_ID CIRCLE_WORKFLOW_ID CI_PIPELINE_ID REWIRE_RUN_ID].freeze

  around do |example|
    saved = CI_VARS.to_h { |k| [k, ENV.delete(k)] }
    example.run
  ensure
    saved.each { |k, v| ENV[k] = v if v }
  end

  it "detects GitHub Actions" do
    ENV["GITHUB_RUN_ID"] = "12345"
    expect(Rewire::Tracer.detect_ci).to eq(
      Rewire::Tracer::CiContext.new(platform: "github_actions", run_id: "12345")
    )
  end

  it "detects CircleCI" do
    ENV["CIRCLE_WORKFLOW_ID"] = "abc-workflow"
    expect(Rewire::Tracer.detect_ci).to eq(
      Rewire::Tracer::CiContext.new(platform: "circleci", run_id: "abc-workflow")
    )
  end

  it "detects GitLab" do
    ENV["CI_PIPELINE_ID"] = "99"
    expect(Rewire::Tracer.detect_ci).to eq(
      Rewire::Tracer::CiContext.new(platform: "gitlab", run_id: "99")
    )
  end

  it "detects generic REWIRE_RUN_ID" do
    ENV["REWIRE_RUN_ID"] = "custom-run"
    expect(Rewire::Tracer.detect_ci).to eq(
      Rewire::Tracer::CiContext.new(platform: "generic", run_id: "custom-run")
    )
  end

  it "returns unknown when no CI env vars are set" do
    expect(Rewire::Tracer.detect_ci).to eq(
      Rewire::Tracer::CiContext.new(platform: "unknown", run_id: nil)
    )
  end

  it "prefers GitHub Actions over other platforms" do
    ENV["GITHUB_RUN_ID"] = "gh-1"
    ENV["CIRCLE_WORKFLOW_ID"] = "ci-1"
    expect(Rewire::Tracer.detect_ci).to eq(
      Rewire::Tracer::CiContext.new(platform: "github_actions", run_id: "gh-1")
    )
  end
end
