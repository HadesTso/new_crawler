package doubangroup

import (
	"github.com/dreamerjackson/crawler/collect"
	"time"
)

var DoubangroupJSTask = &collect.TaskModule{
	Property: collect.Property{
		Name:     "js_find_douban_sun_room",
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   "viewed=\"10558892\"; bid=dztWPHna7Qk; gr_user_id=354f69f8-5205-4d18-a917-79dad40e1702; __gads=ID=5b2d911f1916a907-222490c153d50040:T=1658666835:RT=1658666835:S=ALNI_MaZPfmuU8Jb0iHx2LRbk7FPBBlLRg; ll=\"108306\"; __utmz=30149280.1683950672.2.1.utmcsr=cn.bing.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __gpi=UID=0000080526f457f6:T=1658666835:RT=1683950673:S=ALNI_Mbg7n-Vp1sSu72PkT2gBuV2v_Y8dg; _pk_id.100001.8cb4=5db74a913e7843b6.1688895186.; __utmc=30149280; __yadk_uid=p78l0dgEomMtAsNsyvWGA58kTWXAznqT; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.1099067969.1658666835.1688895187.1690085806.4; __utmt=1; __utmb=30149280.6.5.1690085806",
	},
	Root: `
		var arr = new Array();
		for (var i = 25; i <= 25; i+= 25) {
			var obj = {
				Url: "https://www.douban.com/group/szsh/discussion?start=" + i,
				Priority: 1,
				RuleName: "解析网站URL",
				Method: "GET",
			};
			arr.push(obj);
		};
		console.log(arr[0].Url);
		AddJsReq(arr);
		`,
	Rules: []collect.RuleModule{
		{
			Name: "解析网站URL",
			ParseFunc: `
			ctx.ParseJSReg("解析阳台房", "(https://www.douban.com/group/topic/[0-9a-z]+/)");
			`,
		},
		{
			Name: "解析阳台房",
			ParseFunc: `
			// console.log("parse output");
			ctx.OutputJS("<div class=\"topic-content\">[\\s\\S]*?阳台[\\s\\S]*?<div class=\"aside\">");
			`,
		},
	},
}
