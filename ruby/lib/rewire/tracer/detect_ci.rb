# frozen_string_literal: true

module Rewire
  module Tracer
    CiContext = Struct.new(:platform, :run_id, keyword_init: true)

    def self.detect_ci
      if (run_id = ENV["GITHUB_RUN_ID"])
        CiContext.new(platform: "github_actions", run_id: run_id)
      elsif (run_id = ENV["CIRCLE_WORKFLOW_ID"])
        CiContext.new(platform: "circleci", run_id: run_id)
      elsif (run_id = ENV["CI_PIPELINE_ID"])
        CiContext.new(platform: "gitlab", run_id: run_id)
      elsif (run_id = ENV["REWIRE_RUN_ID"])
        CiContext.new(platform: "generic", run_id: run_id)
      else
        CiContext.new(platform: "unknown", run_id: nil)
      end
    end
  end
end
