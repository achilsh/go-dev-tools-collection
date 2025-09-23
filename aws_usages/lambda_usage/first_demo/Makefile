# 你的 Lambda 名称（AWS 上的 function name）
LAMBDA_NAME=9_3_lmda_func

# 架构: amd64 或 arm64
ARCH=amd64


#  Amazon Linux 2023
RunTimeVersion=provided.al2023   ## provided.al2 (Amazon Linux 2)
# 输出文件
BINARY=bootstrap
ZIPFile=myFunction.zip
ACCOUNT_ID=499258510333
LAMBDA_EXECUTION_ROLE="lambda-ex"

# IAM 角色（只在 create-function 时用）
ROLE_ARN=arn:aws:iam::$(ACCOUNT_ID):role/$(LAMBDA_EXECUTION_ROLE)

.PHONY: build zip deploy create clean

# 1. 编译 Go 程序； 运行 make build 
build:
	GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 go build -o $(BINARY) -tags lambda.norpc main.go

# 2. 打包成 zip; 运行 make zip
zip: build
	zip -j $(ZIPFile) $(BINARY)

# 3. 更新 Lambda 代码; 运行 make deploy
deploy: zip
	aws lambda update-function-code \
		--function-name $(LAMBDA_NAME) \
		--zip-file fileb://$(ZIPFile)

# 4. 首次创建函数（只用一次）; 运行 make  create 
create: zip
	aws lambda create-function \
		--function-name $(LAMBDA_NAME) \
		--handler $(BINARY) \
		--runtime $(RunTimeVersion) \
		--role $(ROLE_ARN) \
		--zip-file fileb://$(ZIPFile)

# 5. 清理
clean:
	rm -f $(BINARY) $(ZIPFile)
