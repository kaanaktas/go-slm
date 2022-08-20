package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.Setenv("GO_SLM_POLICY_RULE_SET_PATH", "/testconfig/policy_rule_set.yaml")
	_ = os.Setenv("GO_SLM_COMMON_POLICIES_PATH", "/testconfig/common_policies.yaml")
	_ = os.Setenv("GO_SLM_CURRENT_MODULE_NAME", "github.com/kaanaktas/dummy")

	os.Exit(m.Run())
}
