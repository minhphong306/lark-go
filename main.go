package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"net/http"
	"os"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func main() {
	err := godotenv.Load(".env")
	appId := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")
	verificationToken := os.Getenv("VERIFICATION_TOKEN")
	encryptionKey := os.Getenv("ENCRYPTION_KEY")

	cli := lark.NewClient(appId, appSecret)

	handler := dispatcher.NewEventDispatcher(verificationToken, encryptionKey).OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		// 处理消息 event，这里简单打印消息的内容
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		fmt.Println("Go here 1")

		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType("open_id").
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(*event.Event.Sender.SenderId.OpenId).
				MsgType("text").
				Content("{\"text\":\"From Vietnam w Love\"}").
				//Uuid("a0d69e20-1dd1-458b-k525-dfeca4015204").
				Build()).
			Build()
		// 发起请求
		resp, err := cli.Im.V1.Message.Create(context.Background(), req)
		fmt.Println(err)
		fmt.Println(resp)

		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		// 处理消息 event，这里简单打印消息的内容
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		fmt.Println("Go here 2")
		return nil
	})

	// 注册 http 路由
	http.HandleFunc("/webhook/event", httpserverext.NewEventHandlerFunc(handler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动 http 服务
	err = http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
