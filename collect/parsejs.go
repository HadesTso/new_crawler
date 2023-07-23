package collect

type (
	TaskModule struct {
		Property
		MaxDepth int64        `json:"max_depth"`
		Root     string       `json:"root_script"`
		Rules    []RuleModule `json:"rules"`
	}

	RuleModule struct {
		Name      string `json:"name"`
		ParseFunc string `json:"parse_func"`
	}
)
