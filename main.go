/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-17 12:00:28
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-21 22:33:19
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
	cacheServer.Run()
	fmt.Println("Running...")
}
