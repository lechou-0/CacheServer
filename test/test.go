/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-27 14:23:07
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-27 20:58:09
 * @Description: test file
 */
package main

import (
	"fmt"
	stProto "test/proto"

	//protobuf编解码库,下面两个库是相互兼容的，可以使用其中任意一个
	"github.com/golang/protobuf/proto"
	//"github.com/gogo/protobuf/proto"
)
type animal struct{
    name string
}
func main() {
	a := animal{"1"}
	fmt.Println(a.name)
	test := &stProto.UserInfo{
		Message: "hello",
	}
	data, _ := proto.Marshal(test)
	newTest := &stProto.UserInfo{}
	proto.Unmarshal(data, newTest)
	fmt.Println(newTest.Message)
}
