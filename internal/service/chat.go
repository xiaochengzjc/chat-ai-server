package service

import (
	"context"
	"gin-gorm/global"
	"gin-gorm/pkg/errcode"
	"gin-gorm/pkg/util"
	"gin-gorm/pkg/ws"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/semaphore"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ChatSvc struct {
	weighted *semaphore.Weighted //全局信号量，限制总共的gptChat请求资源
}

type ChatReq struct {
	Question string `json:"question" binding:"required,min=1,max=500"`
}

type ChatRsp struct {
	Answer string `json:"answer"`
}

func NewChatSvc() *ChatSvc {
	return &ChatSvc{
		weighted: semaphore.NewWeighted(global.ChatSetting.MaxReq),
	}
}

func (c *ChatSvc) Chat(ctx context.Context, msg *Message, send chan<- *Message) {
	//全局限制请求
	ok := c.weighted.TryAcquire(1)
	if !ok {
		send <- &Message{
			Type:    ws.MsgTypeError,
			Content: "前方有多人排队，请稍后",
		}
	}
	defer c.weighted.Release(1)

	var strResp string
	config := openai.DefaultConfig(global.OpenAISetting.AuthToken)
	if global.OpenAISetting.EnableProxy {
		// 设置代理
		proxyUrl, err := url.Parse(global.OpenAISetting.ProxyUrl)
		if err != nil {
			panic(err)
		}
		config.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	if global.OpenAISetting.BaseURL != "" {
		config.BaseURL = global.OpenAISetting.BaseURL
	}
	client := openai.NewClientWithConfig(config)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo0301,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg.Content,
			},
		},
		MaxTokens:        global.OpenAISetting.MaxTokens,
		Temperature:      0.8,
		Stream:           true,
		TopP:             1,
		FrequencyPenalty: 0.1,
		PresencePenalty:  0.1,
	}

	stream, err := client.CreateChatCompletionStream(
		ctx,
		req,
	)

	defer stream.Close()

	if err != nil {
		global.Logger.Errorf("调用openai.CreateChatCompletionStream错误: %v", err)
		chatMsg := Message{
			MsgId:   uuid.New().String(),
			Type:    ws.MsgTypeError,
			Content: errcode.ServerError.Error(),
			MsgTime: util.JsonTime(time.Now()),
		}

		send <- &chatMsg
		global.Logger.Errorf("create chatGPT stream error: %v", err)
	}

	uuid := msg.MsgId
	var i int
	for {
		response, err := stream.Recv()
		if err != nil {
			content := ""
			msgType := ws.MsgTypeError
			if errors.Is(err, io.EOF) { //结束
				if i == 0 { //从开始到结束，没有收到数据
					content = "网络繁忙，请稍后再试"
					global.Logger.Errorf("未接受到openai的流数据: %v", err.Error())
				} else { //正常结束，不做处理
					break
				}
			} else {
				content = "网络繁忙，请稍后再试"
				global.Logger.Errorf("接收openai流数据时错误: %s", err.Error())
			}

			if content != "" {
				chatMsg := Message{
					Type:    msgType,
					Content: content,
					MsgId:   uuid,
					MsgTime: util.JsonTime(time.Now()),
				}
				send <- &chatMsg
			}
			break
		}

		if len(response.Choices) > 0 {
			var content string
			for _, choice := range response.Choices {
				content = content + choice.Delta.Content
			}

			strResp = strResp + content
			//global.Logger.Debugf("第%d个返回: %s", i, content)
			//content = strings.TrimLeft(content, "\n")

			if len(content) > 0 {
				chatMsg := Message{
					Type:    ws.MsgTypeNormal,
					Content: content,
					MsgId:   uuid,
					MsgTime: util.JsonTime(time.Now()),
				}
				send <- &chatMsg
			}
		}
		i++
	}
	if strResp != "" {
		global.Logger.Infof("openai回答: %s", strResp)
	}
}

func (c *ChatSvc) Heart(ctx context.Context, msg *Message, send chan<- *Message) {
	if msg.Type == ws.MsgTypeHeart {
		heartMessage := Message{
			Type:    ws.MsgTypeHeart,
			Content: "pong",
		}
		send <- &heartMessage
	}
}
