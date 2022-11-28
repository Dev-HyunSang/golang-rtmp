package server

import (
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/yutopp/go-rtmp"
)

func RunningServer() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1935")
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	srv := rtmp.NewServer(&rtmp.ServerConfig{
		OnConnect: func(conn net.Conn) (io.ReadWriteCloser, *rtmp.ConnConfig) {
			l := logrus.StandardLogger()
			h := &Handler{}

			return conn, &rtmp.ConnConfig{
				Handler: h,
				ControlState: rtmp.StreamControlStateConfig{
					DefaultBandwidthWindowSize: 6 * 1024 * 1024 / 8,
				},
				Logger: l,
			}
		},
	})
	if err := srv.Serve(listener); err != nil {
		return err
	}

	return err
}
