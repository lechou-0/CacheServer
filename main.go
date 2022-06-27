/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-17 12:00:28
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-27 23:49:18
 * @Description: start cacheserver
 */

package main

import (
	"fmt"
	cacheServer "wsh/server/cacheserver"
)

func main() {
	fmt.Println("Init server...")
	cacheServer := cacheServer.NewServer()
	fmt.Println("Running...")
	cacheServer.Run()
}
