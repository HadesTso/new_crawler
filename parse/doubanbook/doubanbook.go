package doubanbook

import (
	"fmt"
	"github.com/dreamerjackson/crawler/collect"
	"regexp"
	"strconv"
	"time"
)

var DoubanBookTask = &collect.Task{
	Property: collect.Property{
		Name:     "douban_book_list",
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "ll=\"118281\"; bid=gLS7tuqJ8gk; __utmv=30149280.14701; douban-fav-remind=1; viewed=\"1007305\"; __utmz=81379588.1688711664.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); _vwo_uuid_v2=D84EDCCD017EFB1B3E627B4A214AB04BD|561b86d40b1da40b786540c2553fbe20; __yadk_uid=QMU5rulrXuHgYt3SKqsIP5agHnq64Uud; __utmz=30149280.1688720337.4.3.utmcsr=time.geekbang.com|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; __gads=ID=4e9dea1a65c3e408-2262225923e100b6:T=1684492325:RT=1688723594:S=ALNI_MYQjWs1oJm-pxbpKTh6TzF3MkouAw; __gpi=UID=00000c08076793fd:T=1684492325:RT=1688723594:S=ALNI_MZmkLJTYADCK2IUIiyy5Ge9eiDv8w; dbcl2=\"147011365:bLhHRzI7JKs\"; push_doumail_num=0; push_noty_num=0; ck=HVNC; ap_v=0,6.0; _pk_id.100001.3ac3=0f7146999e04b99e.1688711664.2.1690599629.1688711664.; _pk_ses.100001.3ac3=*; __utma=30149280.1572309672.1684133583.1688723141.1690599629.6; __utmc=30149280; __utmt_douban=1; __utmb=30149280.1.10.1690599629; __utma=81379588.1801669177.1688711664.1688711664.1690599629.2; __utmc=81379588; __utmt=1; __utmb=81379588.1.10.1690599629",
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				&collect.Request{
					Priority: 1,
					Url:      "<https://book.douban.com>",
					Method:   "GET",
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag": &collect.Rule{ParseFunc: ParseTag},
			"书籍列表":  &collect.Rule{ParseFunc: ParseBookList},
			"书籍简介": &collect.Rule{
				ItemFields: []string{
					"书名",
					"作者",
					"页数",
					"出版社",
					"得分",
					"价格",
					"简介",
				},
				ParseFunc: ParseBookDetail,
			},
		},
	},
}

const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

func ParseTag(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(regexpStr)

	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		fmt.Println(string(m[1]))
		result.Requesrts = append(
			result.Requesrts, &collect.Request{
				Method:   "GET",
				Task:     ctx.Req.Task,
				Url:      "<https://book.douban.com>" + string(m[1]),
				Depth:    ctx.Req.Depth + 1,
				RuleName: "书籍列表",
			})
	}

	return result, nil
}

const BookListRe = `<a.*?href="([^"]+)" title="([^"]+)">`

func ParseBookList(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(BookListRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)
	result := collect.ParseResult{}
	for _, m := range matches {
		req := &collect.Request{
			Method:   "GET",
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}
		req.TmpData = &collect.Temp{}
		req.TmpData.Set("book_name", string(m[2]))
		result.Requesrts = append(result.Requesrts, req)
	}

	return result, nil
}

var (
	autoRe  = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
	public  = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
	pageRe  = regexp.MustCompile(`<span class="pl">页数:</span>([^<]+)<br/>`)
	priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
	scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
	intoRe  = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)
)

func ParseBookDetail(ctx *collect.Context) (collect.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, autoRe))

	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  ExtraString(ctx.Body, autoRe),
		"页数":  page,
		"出版社": ExtraString(ctx.Body, public),
		"得分":  ExtraString(ctx.Body, scoreRe),
		"价格":  ExtraString(ctx.Body, priceRe),
		"简介":  ExtraString(ctx.Body, intoRe),
	}

	data := ctx.Output(book)

	result := collect.ParseResult{
		Items: []interface{}{data},
	}

	return result, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
