//     Copyright (C) 2020, IrineSistiana
//
//     This file is part of mosdns.
//
//     mosdns is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     mosdns is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with this program.  If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"context"
	"fmt"
	"github.com/IrineSistiana/mosdns/dispatcher/handler"
	"github.com/IrineSistiana/mosdns/dispatcher/logger"
	"github.com/IrineSistiana/mosdns/dispatcher/utils"
	"github.com/miekg/dns"
	"net"
	"time"
)

const (
	serverTCPReadTimeout  = time.Second * 8
	serverTCPWriteTimeout = time.Second
)

type tcpServer struct {
	l       net.Listener
	timeout time.Duration
}

func NewTCPServer(c *Config) Server {
	s := new(tcpServer)
	s.l = c.Listener
	if c.Timeout > 0 {
		s.timeout = c.Timeout
	} else {
		s.timeout = serverTCPReadTimeout
	}
	return s
}

type tcpResponseWriter struct {
	c net.Conn
}

func (t *tcpResponseWriter) Write(m *dns.Msg) (n int, err error) {
	t.c.SetWriteDeadline(time.Now().Add(serverTCPWriteTimeout))
	defer t.c.SetWriteDeadline(time.Time{})

	return utils.WriteMsgToTCP(t.c, m)
}

func (s *tcpServer) ListenAndServe(h Handler) error {
	listenerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		c, err := s.l.Accept()

		if err != nil {
			er, ok := err.(net.Error)
			if ok && er.Temporary() {
				logger.GetStd().Warnf("tcp server: listener: temporary err: %v", err)
				time.Sleep(time.Millisecond * 100)
				continue
			} else {
				return fmt.Errorf("listener: %s", err)
			}
		}

		go func() {
			defer c.Close()
			tcpConnCtx, cancel := context.WithCancel(listenerCtx)
			defer cancel()

			for {
				c.SetReadDeadline(time.Now().Add(serverTCPReadTimeout))
				q, _, err := utils.ReadMsgFromTCP(c)
				if err != nil {
					return // read err, close the conn
				}

				w := &tcpResponseWriter{c: c}
				qCtx := &handler.Context{
					Q:    q,
					From: c.RemoteAddr(),
				}
				go h.ServeDNS(tcpConnCtx, qCtx, w)
			}
		}()
	}
}
