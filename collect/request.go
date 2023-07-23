package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"regexp"
	"sync"
	"time"
)

type Property struct {
	Name     string        `json:"name"` // 任务名称，应保证唯一性
	Url      string        `json:"url"`
	Cookie   string        `json:"cookie"`
	WaitTime time.Duration `json:"wait_time"`
	Reload   bool          `json:"reload"` // 网站是否可以重复爬取
	MaxDepth int64         `json:"max_depth"`
}

// 一个任务实例，
type Task struct {
	Name        string // 用户界面显示的名称（应保证唯一性）
	Url         string
	Cookie      string
	WaitTime    time.Duration
	Reload      bool // 网站是否可以重复爬取
	MaxDepth    int
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Fetcher     Fetcher
	Rule        RuleTree
}

type Context struct {
	Body []byte
	Req  *Request
}

func AddJsReqs(jreqs []map[string]interface{}) []*Request {
	reqs := make([]*Request, 0)

	for _, jreq := range jreqs {
		req := &Request{}
		u, ok := jreq["Url"].(string)
		if !ok {
			return nil
		}

		req.Url = u
		req.RuleName, _ = jreq["RuleName"].(string)
		req.Method, _ = jreq["Method"].(string)
		req.Priority, _ = jreq["Priority"].(int64)
		reqs = append(reqs, req)
	}

	return reqs
}

func (c *Context) ParseJSReg(name string, reg string) ParseResult {
	re := regexp.MustCompile(reg)

	matches := re.FindAll(c.Body, -1)
	result := ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requesrts = append(
			result.Requesrts, &Request{
				Method:   "GET",
				Task:     c.Req.Task,
				Url:      u,
				Depth:    c.Req.Depth + 1,
				RuleName: name,
			})
	}

	return result
}

func (c *Context) OutputJs(reg string) ParseResult {
	re := regexp.MustCompile(reg)
	ok := re.Match(c.Body)

	if !ok {
		return ParseResult{
			Items: []interface{}{},
		}
	}

	result := ParseResult{
		Items: []interface{}{c.Req.Url},
	}

	return result
}

// 单个请求
type Request struct {
	unique   string
	Task     *Task
	Url      string
	Method   string
	Depth    int
	Priority int64
	RuleName string
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("Max depth limit reached")
	}
	return nil
}

// 请求的唯一识别码
func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
