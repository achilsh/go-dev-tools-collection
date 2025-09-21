package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

// 为避免每次调用函数时创建新资源，
// 您可以在 Lambda 函数的处理程序代码之外声明并修改全局变量。
var (
	demo_var int = 123
)

// 在 var 数据块或语句中定义这些全局变量。此外，处理程序可能要声明 init() 函数，该函数在初始化阶段执行。
// init 方法与在 AWS Lambda 中的行为方式相同，正如在标准 Go 程序中一样。

// Lambda 在初始化阶段运行的任何代码
func init() {
	//初始化一些全局变量，比如 s3 client 的连接。
	// 1
	// Lambda 无权在 init() 函数中访问上下文对象。作为解决方法，您可以模拟初始化阶段的 context.TODO() 传入占位符。稍后，使用客户端进行调用时，传入完整的上下文对象。
	// 此解决方法也在 在 AWS SDK 客户端初始化和调用中使用上下文 中进行了介绍。
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	// s3Client = s3.NewFromConfig(cfg)
}

type Order struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
}

// 包含主应用程序逻辑的主处理程序方法。
// 处理程序可能需要 0 到 2 个参数。如果有两个参数，则第一个参数必须实现 context.Context。

// 处理程序可能返回 0 到 2 个参数。如果有一个返回值，则它必须实现 error。
// 如果有两个返回值，则第二个值必须实现 error。
func HandleEvent(ctx context.Context, event json.RawMessage) error {

	// Parse the input event
	var order Order
	if err := json.Unmarshal(event, &order); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return err
	}

	// 获取环境变量
	bucketName := os.Getenv("RECEIPT_BUCKET")
	_ = bucketName

	// 配置并初始化 SDK 客户端后，您可以使用它与其他 AWS 服务进行交互。
	return nil
}

func main() {
	lambda.Start(HandleEvent)

}
