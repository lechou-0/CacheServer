/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-16 06:08:23
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-17 03:15:20
 * @FilePath: /cacheserver/cacheserver.go
 * @Description:
 */
package main

import (
	"fmt"
	"log"
	"sync"
	pb "wsh/cacheserver/proto"

	uuid "github.com/nu7hatch/gouuid"
)

type CacheServer struct {
	Id             string
	IpAddress      string
	StorageManager string
	ReadCache      map[string]pb.KeyValuePair
	ReadCacheLock  *sync.RWMutex
}

/**
 * @Author: lechoux lechoux@qq.com
 * @description: init server ip,cache,mutex
 * @return CacheServer: server
 */

func NewServer() *CacheServer {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Unexpected error while generating UUID: %v", err)
	}
	server := &CacheServer{
		Id:            uid.String(),
		IpAddress:     "127.0.0.1",
		ReadCache:     map[string]pb.KeyValuePair{},
		ReadCacheLock: &sync.RWMutex{},
	}
	fmt.Println("init success")
	return server
}
