package route

import (
	"hash/crc32"

	"github.com/silentred/gateway/util"
)

type Route struct {
	HashID uint32 `json:"id"`
	Host   string `json:"host"`
	Prefix string `json:"prefix"`
}

func NewRoute(host, prefix string) Route {
	return Route{
		Host:   host,
		Prefix: prefix,
		HashID: HashID(host, prefix),
	}
}

func HashID(host, prefix string) uint32 {
	b := make([]byte, 0, len(host)+len(prefix))
	b = append(b, util.Slice(host)...)
	b = append(b, util.Slice(prefix)...)
	return crc32.ChecksumIEEE(b)
}

func HashCrc32(str string) uint32 {
	b := make([]byte, 0, len(str))
	b = append(b, util.Slice(str)...)
	return crc32.ChecksumIEEE(b)
}
