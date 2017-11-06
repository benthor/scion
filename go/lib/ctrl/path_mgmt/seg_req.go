// Copyright 2017 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file contains the Go representation of segment requests.

package path_mgmt

import (
	"fmt"

	"github.com/netsec-ethz/scion/go/lib/addr"
	"github.com/netsec-ethz/scion/go/lib/common"
	"github.com/netsec-ethz/scion/go/proto"
)

type SegReqId uint64

func (id SegReqId) String() string {
	return fmt.Sprintf("%016x", uint64(id))
}

var _ proto.Cerealizable = (*SegReq)(nil)

type SegReq struct {
	Id       SegReqId
	RawSrcIA addr.IAInt `capnp:"srcIA"`
	RawDstIA addr.IAInt `capnp:"dstIA"`
	Flags    struct {
		Sibra     bool
		CacheOnly bool
	}
}

func NewSegReqFromRaw(b common.RawBytes) (*SegReq, error) {
	s := &SegReq{}
	return s, proto.ParseFromRaw(s, s.ProtoId(), b)
}

func (s *SegReq) SrcIA() *addr.ISD_AS {
	return s.RawSrcIA.IA()
}

func (s *SegReq) DstIA() *addr.ISD_AS {
	return s.RawDstIA.IA()
}

func (s *SegReq) ProtoId() proto.ProtoIdType {
	return proto.SegReq_TypeID
}

func (s *SegReq) Write(b common.RawBytes) (int, error) {
	return proto.WriteRoot(s, b)
}

func (s *SegReq) String() string {
	return fmt.Sprintf("Id: %s %s -> %s, Flags: %v", s.Id, s.SrcIA(), s.DstIA(), s.Flags)
}
