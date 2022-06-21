/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-16 21:08:23
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-21 23:47:52
 * @Description:
 */

package cacheserver

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	pb "wsh/server/proto"

	uuid "github.com/nu7hatch/gouuid"
)

type CacheServer struct {
	Id             string
	ResponsePort   int
	StorageManager string
	ReadCache      map[string]pb.KeyValuePair
	ReadCacheLock  *sync.RWMutex
}

func (s *CacheServer) cache_get_bind_proc() string {
	return "ipc://requests/get"
}

func (s *CacheServer) cache_put_bind_proc() string {
	return "ipc://requests/put"
}

func (s *CacheServer) cache_bind_node() string {
	return "tcp://*:" + strconv.Itoa(s.ResponsePort)
}

func (s *CacheServer) write() {
	fmt.Println("1234")
}

func (s *CacheServer) read() {
	fmt.Println("1234")
}

func (s *CacheServer) Run() {
	go s.write()
	go s.read()
}

func NewServer() *CacheServer {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("Unexpected error while generating UUID: %v", err)
	}
	server := &CacheServer{
		Id:            uid.String(),
		ResponsePort:  5555,
		ReadCache:     map[string]pb.KeyValuePair{},
		ReadCacheLock: &sync.RWMutex{},
	}
	fmt.Println("init success")
	return server
}
