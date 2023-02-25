// Copyright (c) 2009-2023 Rob Braun <bbraun@synack.net> and others
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

package tash

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/sfiera/multitalk/pkg/llap"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
		want []llap.Packet
	}{{
		name: "enq-packet",
		data: `0201812dff00fd`,
		want: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeEnq,
			},
		}},
	}, {
		name: "enq-packet_bad_checksum",
		data: `020181eaea00fd`,
		want: nil,
	}, {
		name: "small_data-packet",
		data: `02010100ff028abc00fd`,
		want: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: []byte{0x00, 0x02},
		}},
	}, {
		name: "small_data_with_flags-packet",
		data: `020101fc02226900fd`,
		want: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: []byte{0xfc, 0x02},
		}},
	}, {
		name: "large_data-packet",
		data: `0201010258af16915160fa3f24e0015844fc50b7090588653ccc30baf46c` +
			`b06ef2f0d8c5fe60311e52e6e0218d80f7785abe5a1819668ac9a88edb4b` +
			`35ff4ac9abdd27fef881f9e8c8e520df818abc071ca11c9b4b36c38f43a5` +
			`ca00ff8e8d133d4470c86c4b876a096337e4b73a5157da172374c6839a6c` +
			`93c3f94ceb5c996793683b1676165c2615f52604fdd8463b2ed89f814435` +
			`18b239b8ed8cf793f1439f04db00ff4d3ba7f49eba7b669f141295ad7a6d` +
			`ea1e814c977a61e139106ee599f1df32e7a591bbf6a25883a0bac2d6fb0e` +
			`63be275e6e590be86b1abe8454f9dde225e4ae75950c9705c2a6751cc3f8` +
			`68fe6c0d743f13a2601c5b308573029e51a146cd0499969833a14edbb815` +
			`1b53245d4c8e70c7254f55e9e47373443f4fbc71961a9e7ffa53438b0c89` +
			`9196442ce32d4375328a5dca32cfdb78a254362cd6ecdd0328b142240a3d` +
			`436ba55d481caf46637de7718edf1b6a6ef3b7fcae4d6882e5b6fc5b8508` +
			`3e112388cd4a5edf816adac0d720395908ae6a0c8707c2b0220a21c3b3b0` +
			`2a45ba34fe0799b3c7bb9af274268facbe9922cb7334249e9cb4ac414ed7` +
			`a3591243c1b65a11c26e2bc431667b1396b6393f6f930c2b5641c500ffcf` +
			`1443b9eade4adbcac7a209b15eb75a893330d335f0371b8e07fc85b800ff` +
			`77518a59356c13a100ff67c159d2b2fb9560f4a59794607181a852355fe1` +
			`00ff31680a879f501e78cc691b37fc00ffc50bba0ab84282881153aa1fe6` +
			`282d92ad20f75697e6a4f59e21e3de163e94a98fa34ecce535121eb2c9f5` +
			`3d99a368498fe6effeeee7e2daff40bff239c85b6e4a3bfd0507245f9cea` +
			`36b0f27c97d98a09da9c16fe00fd`,
		want: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: unhex(
				`0258af16915160fa3f24e0015844fc50b7090588653ccc30baf46c` +
					`b06ef2f0d8c5fe60311e52e6e0218d80f7785abe5a1819668ac9a88edb4b` +
					`35ff4ac9abdd27fef881f9e8c8e520df818abc071ca11c9b4b36c38f43a5` +
					`ca008e8d133d4470c86c4b876a096337e4b73a5157da172374c6839a6c93` +
					`c3f94ceb5c996793683b1676165c2615f52604fdd8463b2ed89f81443518` +
					`b239b8ed8cf793f1439f04db004d3ba7f49eba7b669f141295ad7a6dea1e` +
					`814c977a61e139106ee599f1df32e7a591bbf6a25883a0bac2d6fb0e63be` +
					`275e6e590be86b1abe8454f9dde225e4ae75950c9705c2a6751cc3f868fe` +
					`6c0d743f13a2601c5b308573029e51a146cd0499969833a14edbb8151b53` +
					`245d4c8e70c7254f55e9e47373443f4fbc71961a9e7ffa53438b0c899196` +
					`442ce32d4375328a5dca32cfdb78a254362cd6ecdd0328b142240a3d436b` +
					`a55d481caf46637de7718edf1b6a6ef3b7fcae4d6882e5b6fc5b85083e11` +
					`2388cd4a5edf816adac0d720395908ae6a0c8707c2b0220a21c3b3b02a45` +
					`ba34fe0799b3c7bb9af274268facbe9922cb7334249e9cb4ac414ed7a359` +
					`1243c1b65a11c26e2bc431667b1396b6393f6f930c2b5641c500cf1443b9` +
					`eade4adbcac7a209b15eb75a893330d335f0371b8e07fc85b80077518a59` +
					`356c13a10067c159d2b2fb9560f4a59794607181a852355fe10031680a87` +
					`9f501e78cc691b37fc00c50bba0ab84282881153aa1fe6282d92ad20f756` +
					`97e6a4f59e21e3de163e94a98fa34ecce535121eb2c9f53d99a368498fe6` +
					`effeeee7e2daff40bff239c85b6e4a3bfd0507245f9cea36b0f27c97d98a` +
					`09da9c`),
		}},
	}, {
		name: "large_data_all_zeroes-packet",
		data: `0201010258` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff00ff` +
			`a22b00fd` +
			``,
		want: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: unhex(`025800000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000000000000000000000000000000000000000000000000000000000` +
				`000000`),
		}},
	}, {
		name: "large_data-packet_from_uart_frame_aborted",
		data: `0201010258af16915160fa3f24e0015844fc50b7090588653ccc30baf46c` +
			`b06ef2f0d8c5fe60311e52e6e0218d80f7785abe5a1819668ac9a88edb4b` +
			`35ff4ac9abdd27fef881f9e8c8e520df818abc071ca11c9b4b36c38f43a5` +
			`ca00ff8e8d133d4470c86c4b876a096337e4b73a5157da172374c6839a6c` +
			`93c3f94ceb5c996793683b1676165c2615f52604fdd8463b2ed89f814435` +
			`18b239b8ed8cf793f1439f04db00ff4d3ba7f49eba7b669f141295ad7a6d` +
			`ea1e814c977a61e139106ee599f1df32e7a591bbf6a25883a0bac2d6fb0e` +
			`63be275e6e590be86b1abe8454f9dde225e4ae75950c9705c2a6751cc3f8` +
			`68fe6c0d743f13a2601c5b308573029e51a146cd0499969833a14edbb815` +
			`1b53245d4c8e70c7254f55e9e47373443f4fbc71961a9e7ffa53438b0c89` +
			`9196442ce32d4375328a5dca32cfdb78a254362cd6ecdd0328b142240a3d` +
			`436ba55d481caf46637de7718edf1b6a6ef3b7fcae4d6882e5b6fc5b8508` +
			`3e112388cd4a5edf816adac0d720395908ae6a0c8707c2b0220a21c3b3b0` +
			`2a45ba34fe0799b3c7bb9af274268facbe9922cb7334249e9cb4ac414ed7` +
			`a3591243c1b65a11c26e2bc431667b1396b6393f6f930c2b5641c500ffcf` +
			`1443b9eade4adbcac7a209b15eb75a893330d335f0371b8e07fc85b800ff` +
			`77518a59356c13a100ff67c159d2b2fb9560f4a59794607181a852355fe1` +
			`00ff31680a879f501e78cc691b37fc00ffc50bba0ab84282881153aa1fe6` +
			`282d92ad20f75697e6a4f59e21e3de163e94a98fa34ecce535121eb2c9f5` +
			`3d99a368498fe6effeeee7e2daff40bff239c85b6e4a3bfd0507245f9cea` +
			`36b0f27c97d98a09da9c16fe00fa`,
		want: nil,
	}, {
		name: "large_data-packet_from_uart_framing_error",
		data: `0201010258af16915160fa3f24e0015844fc50b7090588653ccc30baf46c` +
			`b06ef2f0d8c5fe60311e52e6e0218d80f7785abe5a1819668ac9a88edb4b` +
			`35ff4ac9abdd27fef881f9e8c8e520df818abc071ca11c9b4b36c38f43a5` +
			`ca00ff8e8d133d4470c86c4b876a096337e4b73a5157da172374c6839a6c` +
			`93c3f94ceb5c996793683b1676165c2615f52604fdd8463b2ed89f814435` +
			`18b239b8ed8cf793f1439f04db00ff4d3ba7f49eba7b669f141295ad7a6d` +
			`ea1e814c977a61e139106ee599f1df32e7a591bbf6a25883a0bac2d6fb0e` +
			`63be275e6e590be86b1abe8454f9dde225e4ae75950c9705c2a6751cc3f8` +
			`68fe6c0d743f13a2601c5b308573029e51a146cd0499969833a14edbb815` +
			`1b53245d4c8e70c7254f55e9e47373443f4fbc71961a9e7ffa53438b0c89` +
			`919600fe`,
		want: nil,
	}, {
		name: "empty-packet",
		data: `00fd`,
		want: nil,
	}, {
		name: "too_short-packet",
		data: `01f1e100fd`, // $f1e1 is valid FCS for $01.
		want: nil,
	}, {
		name: "enq_ack-packet_sequence",
		data: `0201812dff00fd` + `010282ba0800fd`,
		want: []llap.Packet{
			{Header: llap.Header{DstNode: 2, SrcNode: 1, Kind: llap.TypeEnq}},
			{Header: llap.Header{DstNode: 1, SrcNode: 2, Kind: llap.TypeAck}},
		},
	}, {
		name: "enq_abort_ack-packet_sequence",
		data: `0201812dff00fd` + `010200fa` + `010282ba0800fd`,
		want: []llap.Packet{
			{Header: llap.Header{DstNode: 2, SrcNode: 1, Kind: llap.TypeEnq}},
			{Header: llap.Header{DstNode: 1, SrcNode: 2, Kind: llap.TypeAck}},
		},
	}, {
		name: "enq_invalid_ack-packet_sequence",
		data: `0201812dff00fd` + `01f1e100fd` + `010282ba0800fd`,
		want: []llap.Packet{
			{Header: llap.Header{DstNode: 2, SrcNode: 1, Kind: llap.TypeEnq}},
			{Header: llap.Header{DstNode: 1, SrcNode: 2, Kind: llap.TypeAck}},
		},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			buf := bytes.NewBuffer([]byte(unhex(tt.data)))
			d := NewDecoder(buf)
			var packets []llap.Packet
			for {
				pak := llap.Packet{}
				if err := d.Decode(&pak); err == io.EOF {
					break
				} else if err != nil {
					panic(err)
				}
				packets = append(packets, pak)
			}
			assert.Equal(tt.want, packets)
		})
	}
}

const reset = `0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000` +
	`0000000000000000000000000000000000000000000000000000000000000000`

func TestEncode(t *testing.T) {
	for _, tt := range []struct {
		name    string
		packets []llap.Packet
		want    string
	}{{
		name: "enq-packet",
		packets: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeEnq,
			},
		}},
		want: reset + `010201812dff`,
	}, {
		name: "small_data-packet",
		packets: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: []byte{0x00, 0x02},
		}},
		want: reset + `0102010100028abc`,
	}, {
		name: "small_data_with_flags-packet",
		packets: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: []byte{0xfc, 0x02},
		}},
		want: reset + `01020101fc022269`,
	}, {
		name: "large_data-packet",
		packets: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: unhex(
				`0258af16915160fa3f24e0015844fc50b7090588653ccc30baf46c` +
					`b06ef2f0d8c5fe60311e52e6e0218d80f7785abe5a1819668ac9a88edb4b` +
					`35ff4ac9abdd27fef881f9e8c8e520df818abc071ca11c9b4b36c38f43a5` +
					`ca008e8d133d4470c86c4b876a096337e4b73a5157da172374c6839a6c93` +
					`c3f94ceb5c996793683b1676165c2615f52604fdd8463b2ed89f81443518` +
					`b239b8ed8cf793f1439f04db004d3ba7f49eba7b669f141295ad7a6dea1e` +
					`814c977a61e139106ee599f1df32e7a591bbf6a25883a0bac2d6fb0e63be` +
					`275e6e590be86b1abe8454f9dde225e4ae75950c9705c2a6751cc3f868fe` +
					`6c0d743f13a2601c5b308573029e51a146cd0499969833a14edbb8151b53` +
					`245d4c8e70c7254f55e9e47373443f4fbc71961a9e7ffa53438b0c899196` +
					`442ce32d4375328a5dca32cfdb78a254362cd6ecdd0328b142240a3d436b` +
					`a55d481caf46637de7718edf1b6a6ef3b7fcae4d6882e5b6fc5b85083e11` +
					`2388cd4a5edf816adac0d720395908ae6a0c8707c2b0220a21c3b3b02a45` +
					`ba34fe0799b3c7bb9af274268facbe9922cb7334249e9cb4ac414ed7a359` +
					`1243c1b65a11c26e2bc431667b1396b6393f6f930c2b5641c500cf1443b9` +
					`eade4adbcac7a209b15eb75a893330d335f0371b8e07fc85b80077518a59` +
					`356c13a10067c159d2b2fb9560f4a59794607181a852355fe10031680a87` +
					`9f501e78cc691b37fc00c50bba0ab84282881153aa1fe6282d92ad20f756` +
					`97e6a4f59e21e3de163e94a98fa34ecce535121eb2c9f53d99a368498fe6` +
					`effeeee7e2daff40bff239c85b6e4a3bfd0507245f9cea36b0f27c97d98a` +
					`09da9c`),
		}},
		want: reset + `01020101` +
			`0258af16915160fa3f24e0015844fc50b7090588653ccc30baf46c` +
			`b06ef2f0d8c5fe60311e52e6e0218d80f7785abe5a1819668ac9a88edb4b` +
			`35ff4ac9abdd27fef881f9e8c8e520df818abc071ca11c9b4b36c38f43a5` +
			`ca008e8d133d4470c86c4b876a096337e4b73a5157da172374c6839a6c93` +
			`c3f94ceb5c996793683b1676165c2615f52604fdd8463b2ed89f81443518` +
			`b239b8ed8cf793f1439f04db004d3ba7f49eba7b669f141295ad7a6dea1e` +
			`814c977a61e139106ee599f1df32e7a591bbf6a25883a0bac2d6fb0e63be` +
			`275e6e590be86b1abe8454f9dde225e4ae75950c9705c2a6751cc3f868fe` +
			`6c0d743f13a2601c5b308573029e51a146cd0499969833a14edbb8151b53` +
			`245d4c8e70c7254f55e9e47373443f4fbc71961a9e7ffa53438b0c899196` +
			`442ce32d4375328a5dca32cfdb78a254362cd6ecdd0328b142240a3d436b` +
			`a55d481caf46637de7718edf1b6a6ef3b7fcae4d6882e5b6fc5b85083e11` +
			`2388cd4a5edf816adac0d720395908ae6a0c8707c2b0220a21c3b3b02a45` +
			`ba34fe0799b3c7bb9af274268facbe9922cb7334249e9cb4ac414ed7a359` +
			`1243c1b65a11c26e2bc431667b1396b6393f6f930c2b5641c500cf1443b9` +
			`eade4adbcac7a209b15eb75a893330d335f0371b8e07fc85b80077518a59` +
			`356c13a10067c159d2b2fb9560f4a59794607181a852355fe10031680a87` +
			`9f501e78cc691b37fc00c50bba0ab84282881153aa1fe6282d92ad20f756` +
			`97e6a4f59e21e3de163e94a98fa34ecce535121eb2c9f53d99a368498fe6` +
			`effeeee7e2daff40bff239c85b6e4a3bfd0507245f9cea36b0f27c97d98a` +
			`09da9c` +
			`16fe`,
	}, {
		name: "large_data_all_zeroes-packet",
		packets: []llap.Packet{{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: unhex(
				`025800000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000000000000000000000000000000000000000000000000000000000` +
					`000000`),
		}},
		want: reset + `01020101` +
			`025800000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000000000000000000000000000000000000000000000000000000000` +
			`000000a22b`,
	}, {
		name: "enq_ack-packet_sequence",
		packets: []llap.Packet{
			{Header: llap.Header{DstNode: 2, SrcNode: 1, Kind: llap.TypeEnq}},
			{Header: llap.Header{DstNode: 1, SrcNode: 2, Kind: llap.TypeAck}},
		},
		want: reset + `010201812dff` + `01010282ba08`,
	}} {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			buf := bytes.Buffer{}
			e := NewEncoder(&buf)
			for _, pak := range tt.packets {
				if err := e.Encode(pak); err != nil {
					panic(err)
				}
			}
			assert.Equal(unhex(tt.want), buf.Bytes())
		})
	}
}

func TestEncodeBadPackets(t *testing.T) {
	for _, tt := range []struct {
		name    string
		packet  llap.Packet
		wantErr string
	}{{
		name: "invalid-type",
		packet: llap.Packet{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    13,
			},
		},
		wantErr: `invalid packet type: $0d`,
	}, {
		name: "non-empty-enq",
		packet: llap.Packet{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeEnq,
			},
			Payload: []byte{0x00, 0x02},
		},
		wantErr: `control frame packet with payload`,
	}, {
		name: "non-empty-ack",
		packet: llap.Packet{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeAck,
			},
			Payload: []byte{0x00, 0x02},
		},
		wantErr: `control frame packet with payload`,
	}, {
		name: "length-mismatch",
		packet: llap.Packet{
			Header: llap.Header{
				DstNode: 2,
				SrcNode: 1,
				Kind:    llap.TypeDDP,
			},
			Payload: []byte{0x00, 0x04, 0x03},
		},
		wantErr: `DDP packet length mismatch: 3 vs. 4`,
	}} {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			buf := bytes.Buffer{}
			e := NewEncoder(&buf)
			if err := e.Encode(tt.packet); assert.Error(err) {
				assert.Equal(tt.wantErr, err.Error())
			}
			assert.Equal(unhex(reset), buf.Bytes())
		})
	}
}

func TestSetNodeIDs(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.Buffer{}
	e := NewEncoder(&buf)
	if assert.NoError(e.SetNodeIDs(NewNodeSet(1, 2, 3, 4, 5, 254))) {
		assert.Equal(unhex(
			reset+`02`+
				`3e00000000000000`+
				`0000000000000000`+
				`0000000000000000`+
				`0000000000000040`), buf.Bytes())
	}
}

func unhex(s string) []byte {
	data := []byte{}
	for i := 0; i < len(s); i += 2 {
		n, err := strconv.ParseUint(s[i:i+2], 16, 8)
		if err != nil {
			panic(err)
		}
		data = append(data, byte(n))
	}
	return data
}
