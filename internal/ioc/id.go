package ioc

import (
	"errors"
	"github.com/sony/sonyflake"
	"go.uber.org/fx"
	"net"
	"time"
)

var IdGeneratorFxOpt = fx.Provide(InitIdGenerator)

func InitIdGenerator() *sonyflake.Sonyflake {
	return sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: getMachineIdByIp,
	})
}

func getMachineIdByIp() (uint16, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip := ipNet.IP.To4()
			return uint16(ip[2])<<8 | uint16(ip[3]), nil
		}
	}

	return 0, errors.New("no valid IPv4 address found")
}
