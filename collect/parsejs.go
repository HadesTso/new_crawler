package collect

import "time"

type (
	TaskModule struct {
		Name     string        `json:"name"`
		Url      string        `json:"url"`
		Cookie   string        `json:"cookie"`
		WaitTime time.Duration `json:"wait_time"`
		Reload   bool          `json:"reload"`
		MaxDepth int64         `json:"max_depth"`
		Root     string        `json:"root_script"`
		Rules    []RuleModule  `json:"rules"`
	}

	RuleModule struct {
		Name      string `json:"name"`
		ParseFunc string `json:"parse_func"`
	}
)
