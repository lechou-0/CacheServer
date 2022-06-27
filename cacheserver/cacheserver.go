/*
 * @Author: lechoux lechoux@qq.com
 * @Date: 2022-06-16 21:08:23
 * @LastEditors: lechoux lechoux@qq.com
 * @LastEditTime: 2022-06-27 22:31:43
 * @Description:
 */

package cacheserver

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	pb "wsh/server/proto"

	"github.com/golang/protobuf/proto"
	proto "github.com/golang/protobuf/proto"
	uuid "github.com/nu7hatch/gouuid"
	zmq "github.com/pebbe/zmq4"
)

type socketCache struct {
	ctx         zmq.Context
	pushercache map[string]zmq.Socket
}

/**
 * @Author: lechoux lechoux@qq.com
 * @description: clear push socketcache. Maybe it can be optimized with LRU
 * @return {*}
 */
func (socketcache *socketCache) clear_sockets() {
	socketcache.pushercache = make(map[string]zmq.Socket)
}

type CacheServer struct {
	Id             string
	ResponsePort   int
	StorageManager string
	ReadCache      map[string]pb.KeyValuePair
	ReadCacheLock  *sync.RWMutex
	pushers        socketCache
}

func (s *CacheServer) cache_get_bind_proc() string {
	return "ipc://requests/get"
}

func (s *CacheServer) cache_set_bind_proc() string {
	return "ipc://requests/put"
}

func (s *CacheServer) cache_bind_node() string {
	return "tcp://*:" + strconv.Itoa(s.ResponsePort)
}

/**
 * @Author: lechoux lechoux@qq.com
 * @description: traverse socket and process W/R request
 * @param {[]zmq.Polled} sockets
 * @param {zmq.Socket} get_puller
 * @param {zmq.Socket} set_puller
 * @param {zmq.Socket} pusher
 * @return {*}
 */
func (s *CacheServer) execWR(sockets []zmq.Polled, get_puller zmq.Socket, set_puller zmq.Socket) {
	for _, polled := range sockets {
		switch s := polled.Socket; s {
		case get_puller:
			msg, _ := s.Recv(0)
			getRequest := &pb.KeyRequest{}
			getResponse := &pb.KeyResponse{}
			err := proto.Unmarshal(msg, getRequest)
			if err != nil {
				fmt.Println("get unmarshaling error:" + err)
				break
			}

			getResponse.Key = getRequest.GetKey()
			// get value from cache
			s.ReadCacheLock.RLock()
			val, ok := s.ReadCache[getRequest.GetKey()]
			s.ReadCacheLock.RUnlock()

			if ok {
				getResponse.Value = val
				getResponse.Error = pb.CacheError_NO_ERROR
			} else {
				getResponse.Error = pb.CacheError_KEY_DNE
			}

			resp_string, err := proto.Marshal(getResponse)
			if err != nil {
				fmt.Println("get marshaling error:" + err)
				break
			}

			// send msg by pushercache or new pusher
			endpoint := getRequest.GetResponseAddress()
			if socket, ok := s.pushers.pushercache[endpoint]; ok {
				socket.Send(resp_string)
			} else {
				pusher := s.pushers.ctx.socket(zmq.PUSH)
				pusher.Connect(endpoint)
				s.pushers.pushercache[endpoint] = pusher
				pusher.Send(resp_string)
			}

		case set_puller:
			msg, _ := s.Recv(0)
			setRequest := &pb.KeyRequest{}
			setResponse := &pb.KeyResponse{}
			err := proto.Unmarshal(msg, setRequest)
			if err != nil {
				fmt.Println("set unmarshaling error" + err)
				break
			}

			s.ReadCacheLock.RLock()
			s.ReadCache[setRequest.GetKey()] = setRequest.GetValue()
			s.ReadCacheLock.RUnlock()

			setResponse.Error = pb.CacheError_NO_ERROR

			resp_string, err := proto.Marshal(setResponse)
			if err != nil {
				fmt.Println("set marshaling error:" + err)
				break
			}

			// send msg by pushercache or new pusher
			endpoint := setRequest.GetResponseAddress()
			if socket, ok := s.pushers.pushercache[endpoint]; ok {
				socket.Send(resp_string)
			} else {
				pusher := s.pushers.ctx.socket(zmq.PUSH)
				pusher.Connect(endpoint)
				s.pushers.pushercache[endpoint] = pusher
				pusher.Send(resp_string)
			}
		}
	}
}

func (s *CacheServer) Run() {

	zctx, _ := zmq.NewContext()
	s.pushers.ctx = zctx
	defer zctx.Close()

	get_puller := zctx.socket(zmq.PULL)
	defer get_puller.Close()
	get_puller.bind(s.cache_get_bind_proc())

	set_puller := zctx.socket(zmq.PULL)
	defer set_puller.Close()
	set_puller.bind(s.cache_set_bind_proc())

	socketPoller := zmq.NewPoller()
	socketPoller.Add(get_puller, zmq.POLLIN)
	socketPoller.Add(set_puller, zmq.POLLIN)

	for {
		sockets, _ := socketPoller.Poll(-1)
		go s.execWR(sockets, get_puller, set_puller)
	}
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
