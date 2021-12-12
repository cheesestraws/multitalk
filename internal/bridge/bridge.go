// Copyright (c) 2009-2020 Rob Braun <bbraun@synack.net> and others
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. Neither the name of Rob Braun nor the names of his contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

// definition of bridges
package bridge

import (
	"context"

	"github.com/sfiera/multitalk/pkg/ethertalk"
)

type (
	iface struct {
		send chan<- ethertalk.Packet
		recv <-chan ethertalk.Packet
	}

	packetFrom struct {
		packet *ethertalk.Packet
		send   chan<- ethertalk.Packet
	}

	Bridge interface {
		Start(ctx context.Context) (
			send chan<- ethertalk.Packet,
			recv <-chan ethertalk.Packet,
		)
	}

	Group struct {
		recvCh chan func(*Group)
		sendCh []chan<- ethertalk.Packet
	}
)

func NewGroup() *Group {
	return &Group{
		make(chan func(*Group)),
		nil,
	}
}

func (g *Group) Add(send chan<- ethertalk.Packet, recv <-chan ethertalk.Packet) {
	go func() {
		g.recvCh <- add(send)
		for pak := range recv {
			g.recvCh <- broadcast(pak, send)
		}
		g.recvCh <- remove(send)
	}()
}

func (g *Group) Run() {
	for fn := range g.recvCh {
		fn(g)
	}
}

func broadcast(pak ethertalk.Packet, send chan<- ethertalk.Packet) func(g *Group) {
	return func(g *Group) {
		for _, sendCh := range g.sendCh {
			if sendCh != send {
				sendCh <- pak
			}
		}
	}
}

func add(send chan<- ethertalk.Packet) func(g *Group) {
	return func(g *Group) {
		g.sendCh = append(g.sendCh, send)
	}
}

func remove(send chan<- ethertalk.Packet) func(g *Group) {
	return func(g *Group) {
		var newCh []chan<- ethertalk.Packet
		for _, ch := range g.sendCh {
			if ch != send {
				newCh = append(newCh, ch)
			}
		}
		g.sendCh = newCh
		close(send)
	}
}
