/*
 * 昵称系统
 * 让Bot称呼你时使用昵称称呼
 */
package nickname

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wxnacy/wgo/arrays"
)

// 机器人别名配置结构体
type Config struct {
	BotNickNames []string `yaml:"NickNames"`
}

var botConfig Config
var botNickNames []string

// 昵称所屏蔽的关键词，会被替换为 *
var blackWords = []string{
	"爸",
	"爹",
	"爷",
}

func init() {
	// 读取机器人配置
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Error(err)
	}
	err = yaml.Unmarshal(configFile, &botConfig)
	if err != nil {
		log.Error(err)
	}
	botNickNames = botConfig.BotNickNames

	// 设置随机数种子
	rand.Seed(time.Now().Unix())

	// 注册指令
	engine := zero.New()
	// 注册设置昵称指令
	registerNicknameCommands := []string{
		"nickname",
		"以后叫我",
		"以后请叫我",
		"称呼我",
		"以后请称呼我",
		"以后称呼我",
		"叫我",
		"请叫我",
	}
	engine.OnCommandGroup(registerNicknameCommands).Handle(func(ctx *zero.Ctx) {
		var cmd extension.CommandModel
		err := ctx.Parse(&cmd)
		if err != nil {
			ctx.Send("心爱酱不知道你在说什么哦")
			log.Errorf("处理 %v 命令发生错误: %v", cmd.Command, err)
		}
		reply, err := registerNickname(ctx.Event.Sender.ID, ctx.Event.Sender.NickName, cmd.Args)
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text(reply)))
		if err != nil {
			log.Error(err)
		}
	})
	// 注册设置昵称指令
	queryMyNicknameCommands := []string{
		"我叫什么",
		"我是谁",
		"我的名字",
	}
	engine.OnCommandGroup(queryMyNicknameCommands).Handle(func(ctx *zero.Ctx) {
		reply, err := queryMyNickname(ctx.Event.Sender.ID)
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text(reply)))
		if err != nil {
			log.Error(err)
		}
	})
	// 注册取消昵称指令
	cancelNicknameCommands := []string{
		"取消昵称",
	}
	engine.OnCommandGroup(cancelNicknameCommands).Handle(func(ctx *zero.Ctx) {
		reply, err := cancelNickname(ctx.Event.Sender.ID)
		ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text(reply)))
		if err != nil {
			log.Error(err)
		}
	})
}

// 注册昵称
func registerNickname(QID int64, userName string, arg string) (reply string, err error) {
	msg := strings.Trim(arg, " ")
	fmt.Println(msg)
	// 特殊情况判断
	if msg == "" {
		reply = "叫你空白？叫你虚空？叫你无名？？"
		err = nil
		return
	}
	if len([]rune(msg)) > 10 {
		reply = "昵称可不能超过10个字！"
		err = nil
		return
	}
	if index := arrays.Contains(botNickNames, msg); index != -1 {
		reply = "笨蛋！休想占用我的名字！"
		err = nil
		return
	}
	// 将黑名单中的字符替换为 *
	tmp := ""
	for _, value := range msg {
		if index := arrays.ContainsString(blackWords, string(value)); index != -1 {
			tmp = tmp + "*"
		} else {
			tmp = tmp + string(value)
		}
	}
	msg = tmp

	// 设置昵称
	err = UpdateNickname(QID, userName, msg)
	if err != nil {
		reply = "设置昵称失败了，明天再来试一试！或者联系管理员！"
		return
	} else {
		replies := [...]string{
			"好啦好啦，心爱酱知道啦，%s，以后就这么叫你吧",
			"嗯嗯，记住你的昵称了哦，%s",
			"好突然，突然要叫你昵称什么的...%s...",
			"心爱酱会好好记住的%s的，放心吧",
			"好..好.，那心爱酱以后就叫你%s了...",
		}
		reply = fmt.Sprintf(replies[rand.Intn(len(replies))], msg)
		return
	}
}

// 查询昵称
func queryMyNickname(QID int64) (reply string, err error) {
	nickname, err := QueryNickname(QID)
	if err != nil {
		if err.Error() == "Queried QID not exsits!" {
			replies := [...]string{
				"没..没有昵称嘛！",
				"心爱酱也不知道你的昵称哦，快告诉心爱酱吧！",
				"嗯？你是...？",
				"你还没有设置昵称吧？快去设置吧！",
				"你在做梦吗？你没有昵称啊",
			}
			reply = replies[rand.Intn(len(replies))]
			return
		} else {
			replies := [...]string{
				"忘记你的名字了呢...真的是太不好意思啦...",
				"心爱酱想不起来了呢...",
				"心爱酱忘记了呢，你能原谅心爱酱一次吗？",
			}
			reply = replies[rand.Intn(len(replies))]
			return
		}
	} else {
		replies := [...]string{
			"我肯定记得你啊，你是%s啊",
			"我不会忘记你的，你也不要忘记我！%s!",
			"哼哼，心爱酱记忆力可是很好的，%s",
			"嗯？你是失忆了嘛...%s...",
			"不要小看心爱酱的记忆力啊！笨蛋%s！",
			"哎？%s...怎么了吗...突然这样问...",
			"不记得了...开玩笑哒，你是%s!",
		}
		reply = fmt.Sprintf(replies[rand.Intn(len(replies))], nickname)
		return
	}
}

// 取消设置昵称
func cancelNickname(QID int64) (reply string, err error) {
	err = RemoveNickname(QID)
	if err != nil {
		if err.Error() == "Document not exsits!" {
			replies := [...]string{
				"没..没有昵称嘛！怎么忘记你嘛！",
				"嗯？你是...？我都不知道你的昵称，怎么忘记你嘛！",
				"你在做梦吗？你没有昵称啊",
			}
			reply = replies[rand.Intn(len(replies))]
			return
		} else {
			replies := [...]string{
				"心爱酱才不会忘记你呢！",
				"哼哼，心爱酱记忆力可是很好的！",
				"想让我忘记你？不！可！能！",
			}
			reply = replies[rand.Intn(len(replies))]
			return
		}
	} else {
		replies := [...]string{
			"呜...睡一觉就会忘记的..和梦一样...",
			"窝知道了...",
			"是哪里做的不好嘛..好吧..晚安...",
			"呃，下次我绝对绝对绝对不会再忘记你！",
			"可..可恶！太可恶了！呜",
		}
		reply = replies[rand.Intn(len(replies))]
		return
	}
}
