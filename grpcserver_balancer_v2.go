package cnt2

import (
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"log"
	"strconv"
)

var grpcServerChan = GrpcServerChan{ch: make(chan []GrpcServerInfo)}

type GrpcServerChan struct {
	ch chan []GrpcServerInfo
}

type GrpcServerBalancer struct {
	resolver *manual.Resolver
	cleanup  func()
}

func (c *GrpcServerBalancer) Init() {
	c.resolver, c.cleanup = manual.GenerateAndRegisterManualResolver()
	go func(resover *manual.Resolver) {
		for serverInfo := range grpcServerChan.ch {
			log.Printf("the New grpcServer: %s", serverInfo)
			if newBestAddress, ok := ChooseBestAddress(serverInfo); ok {
				var resolvedAddrs []resolver.Address
				for i := 0; i < len(newBestAddress); i++ {
					resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: newBestAddress[i]})
				}
				resover.NewAddress(resolvedAddrs)
			}
		}
	}(c.resolver)
}

//FIXME can be better?
func ChooseBestAddress(grpcServerInfos []GrpcServerInfo) ([]string, bool) {

	if len(hostInfo.Ips) == 0 {
		for k := range NetTypeMap {
			//不能是内网
			if k == 2457 {
				continue
			}
			if v, ok := grpcServerInfos[0].ServerIP[k]; ok {
				return StringArray{v + ":" + strconv.Itoa(grpcServerInfos[0].Port)}, true
			}
		}
	}
	var sameGroup, bestChoice, sameNetType, random []string
	for _, serverInfo := range grpcServerInfos {
		//同机房
		if hostInfo.GroupId == serverInfo.GroupId {
			for _, ip := range hostInfo.Ips {
				if v, ok := serverInfo.ServerIP[ip.Type.Val]; ok {
					//优先同运营商
					bestChoice = append(bestChoice, v+":"+strconv.Itoa(serverInfo.Port))
					break
				}
			}
			if len(bestChoice) == 0 {
				for k := range NetTypeMap {
					//不能是内网
					if k == 2457 {
						continue
					}
					if v, ok := serverInfo.ServerIP[k]; ok {
						sameGroup = append(sameGroup, v+":"+strconv.Itoa(serverInfo.Port))
						break
					}
				}
			}
			//不同机房
		} else {
			for _, ip := range hostInfo.Ips {
				if v, ok := serverInfo.ServerIP[ip.Type.Val]; ok {
					//运营商
					sameNetType = append(sameNetType, v+":"+strconv.Itoa(serverInfo.Port))
					break
				}
				if len(sameNetType) == 0 {
					//随机
					for k := range NetTypeMap {
						//不能是内网
						if k == 2457 {
							continue
						}
						if v, ok := serverInfo.ServerIP[k]; ok {
							random = append(random, v+":"+strconv.Itoa(serverInfo.Port))
						}
					}
				}
			}
		}
	}
	if len(bestChoice) > 0 {
		return bestChoice, true
	}
	if len(sameGroup) > 0 {
		return sameGroup, true
	}
	if len(sameNetType) > 0 {
		return sameNetType, true
	}
	if len(random) > 0 {
		return random, true
	}
	return random, false
}
