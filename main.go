package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"strings"
)

func main() {
	res := map[string]string{}
	if len(os.Args) < 6 {
		res["success"] = "false"
		res["message"] = "param error"
		fmt.Println(MarshalRes(res))
		return

	} else {
		fmt.Println(os.Args)
	}
	VITE_KEY_ID := os.Args[1]
	VITE_ACCESS_KEY := os.Args[2]
	VITE_REGION := os.Args[3]
	FILE := os.Args[4]
	INPUTTEXT := os.Args[5]

	// 配置AWS会话
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(VITE_REGION),
		Credentials: credentials.NewStaticCredentials(VITE_KEY_ID, VITE_ACCESS_KEY, ""),
	})
	if err != nil {
		res["success"] = "false"
		res["message"] = "connect s3 error"
		fmt.Println(MarshalRes(res))
		return
	}

	// 创建S3服务客户端
	svc := s3.New(sess)

	//// 定义要下载的文件信息
	//bucket := VITE_YOUR_BUCKET_NAME
	//key := bucket + "/1.text"

	// 指定下载文件的本地路径
	localFilePath := strings.Split(FILE, "/")[1]

	// 创建文件
	file, err := os.Create(localFilePath)
	if err != nil {
		res["success"] = "false"
		res["message"] = "could not create file"
		fmt.Println(MarshalRes(res))
		return
	}
	defer file.Close()

	// 下载S3对象的内容
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(strings.Split(FILE, "/")[0]),
		Key:    aws.String(FILE),
	})
	if err != nil {
		res["success"] = "false"
		res["message"] = "could not down file"
		fmt.Println(MarshalRes(res))
		return
	}

	// 将对象内容写入文件流
	if _, err := file.ReadFrom(result.Body); err != nil {
		res["success"] = "false"
		res["message"] = "unable save file"
		fmt.Println(MarshalRes(res))
		return
	}

	fmt.Println("文件下载完成")

	// 读取1.txt文件的内容
	content, err := os.ReadFile(strings.Split(FILE, "/")[1])
	if err != nil {
		res["success"] = "false"
		res["message"] = "cant read input_file"
		fmt.Println(MarshalRes(res))
		return
	}

	// 追加字符串到内容中
	newContent := append(content, []byte(":"+INPUTTEXT)...)

	// 写入2.txt文件
	err = os.WriteFile(strings.Replace(strings.Split(FILE, "/")[1], ".", ".out.", -1), newContent, 0644)
	if err != nil {
		res["success"] = "false"
		res["message"] = "cant save out_file"
		fmt.Println(MarshalRes(res))
		return
	}

	// 打开本地文件
	outFile, outFileErr := os.Open(strings.Replace(strings.Split(FILE, "/")[1], ".", ".out.", -1))
	if outFileErr != nil {
		res["success"] = "false"
		res["message"] = "cant open out_file"
		fmt.Println(MarshalRes(res))
		return
	}
	defer outFile.Close()
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(strings.Split(FILE, "/")[0]),
		Key:    aws.String(strings.Replace(FILE, ".", ".out.", -1)),
		Body:   outFile,
	})
	if err != nil {
		res["success"] = "false"
		res["message"] = "cant upload out_file"
		fmt.Println(MarshalRes(res))
		return
	}
	res["success"] = "true"
	fmt.Println(MarshalRes(res))
	return
}

func MarshalRes(m map[string]string) string {
	jsonData, _ := json.Marshal(m)

	return string(jsonData)
}
