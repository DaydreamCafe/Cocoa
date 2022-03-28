// 用Go实现的世界上最牛逼的Ai 估值十亿
package smartai

import (
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	// "github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	zero.OnMessage().
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.IsToMe {
				var msg string
				msg = ctx.Event.Message.String()
				msg = strings.Replace(msg, "?", "!", -1);
				msg = strings.Replace(msg, "？", "!", -1);
				msg = strings.Replace(msg, "吗", "", -1);
				msg = strings.Replace(msg, "我", "你", -1);
				ctx.Send(msg)
			}
		})
}