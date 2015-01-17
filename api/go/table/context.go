// Copyright 2015 stevejiang. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package table

import (
	"github.com/stevejiang/gotable/api/go/table/proto"
)

// Connection Context to GoTable server.
// It's NOT safe to use in multiple goroutines.
type Context struct {
	cli  *Client
	dbId uint8
}

type Call struct {
	Done chan *Call  // Reply channel
	args interface{} // Request args
	err  error
	pkg  []byte
	seq  uint64
	cmd  uint8
}

// Get the underling connection Client of the Context.
func (c *Context) Client() *Client {
	return c.cli
}

func (c *Context) Use(dbId uint8) {
	c.dbId = dbId
}

func (c *Context) Ping() error {
	call := c.GoPing(nil)
	if call.err != nil {
		return call.err
	}

	_, err := (<-call.Done).Reply()
	return err
}

func (c *Context) Get(tableId uint8, rowKey, colKey []byte,
	cas uint32) (*OneReply, error) {
	call := c.GoGet(tableId, rowKey, colKey, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) ZGet(tableId uint8, rowKey, colKey []byte,
	cas uint32) (*OneReply, error) {
	call := c.GoZGet(tableId, rowKey, colKey, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) Set(tableId uint8, rowKey, colKey, value []byte, score int64,
	cas uint32) (*OneReply, error) {
	call := c.GoSet(tableId, rowKey, colKey, value, score, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) ZSet(tableId uint8, rowKey, colKey, value []byte, score int64,
	cas uint32) (*OneReply, error) {
	call := c.GoZSet(tableId, rowKey, colKey, value, score, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) Del(tableId uint8, rowKey, colKey []byte,
	cas uint32) (*OneReply, error) {
	call := c.GoDel(tableId, rowKey, colKey, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) ZDel(tableId uint8, rowKey, colKey []byte,
	cas uint32) (*OneReply, error) {
	call := c.GoZDel(tableId, rowKey, colKey, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) Incr(tableId uint8, rowKey, colKey []byte, score int64,
	cas uint32) (*OneReply, error) {
	call := c.GoIncr(tableId, rowKey, colKey, score, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) ZIncr(tableId uint8, rowKey, colKey []byte, score int64,
	cas uint32) (*OneReply, error) {
	call := c.GoZIncr(tableId, rowKey, colKey, score, cas, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*OneReply), nil
}

func (c *Context) MGet(args *MultiArgs) (*MultiReply, error) {
	call := c.GoMGet(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) ZmGet(args *MultiArgs) (*MultiReply, error) {
	call := c.GoZmGet(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) MSet(args *MultiArgs) (*MultiReply, error) {
	call := c.GoMSet(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) ZmSet(args *MultiArgs) (*MultiReply, error) {
	call := c.GoZmSet(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) MDel(args *MultiArgs) (*MultiReply, error) {
	call := c.GoMDel(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) ZmDel(args *MultiArgs) (*MultiReply, error) {
	call := c.GoZmDel(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) MIncr(args *MultiArgs) (*MultiReply, error) {
	call := c.GoMIncr(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) ZmIncr(args *MultiArgs) (*MultiReply, error) {
	call := c.GoZmIncr(args, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*MultiReply), nil
}

func (c *Context) Scan(tableId uint8, rowKey, colKey []byte,
	asc bool, num int) (*ScanReply, error) {
	call := c.GoScan(tableId, rowKey, colKey, asc, num, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*ScanReply), nil
}

func (c *Context) ZScan(tableId uint8, rowKey, colKey []byte, score int64,
	asc, orderByScore bool, num int) (*ScanReply, error) {
	call := c.GoZScan(tableId, rowKey, colKey, score, asc, orderByScore, num, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*ScanReply), nil
}

func (c *Context) scanMore(zop bool, last *ScanReply) (*ScanReply, error) {
	if last.End || len(last.Reply) == 0 {
		return nil, ErrScanEnded
	}
	var a = last.Reply[len(last.Reply)-1].KeyValue
	var call *Call
	if zop {
		call = c.GoZScan(a.TableId, a.RowKey, a.ColKey, a.Score,
			last.ctx.asc, last.ctx.orderByScore, last.ctx.num, nil)
	} else {
		call = c.GoScan(a.TableId, a.RowKey, a.ColKey,
			last.ctx.asc, last.ctx.num, nil)
	}
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*ScanReply), nil
}

func (c *Context) ScanMore(last *ScanReply) (*ScanReply, error) {
	return c.scanMore(false, last)
}

func (c *Context) ZScanMore(last *ScanReply) (*ScanReply, error) {
	return c.scanMore(true, last)
}

func (c *Context) Dump(scope, tableId uint8) (*DumpReply, error) {
	var dbId = c.dbId
	// Only scan full DB can escape c.dbId
	if ScopeFullDB == scope {
		dbId = 0
		tableId = 0
	}
	rec := &DumpRecord{dbId, 0, &proto.KeyValue{tableId, nil, nil, nil, 0, 0}}

	call := c.goDump(scope, 0, rec, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}
	return r.(*DumpReply), nil
}

func (c *Context) DumpMore(last *DumpReply) (*DumpReply, error) {
	if last.End {
		return nil, ErrScanEnded
	}

	var rec *DumpRecord
	var unitId = last.ctx.unitId
	if len(last.Reply) == 0 {
		unitId += 1
		rec = &DumpRecord{last.ctx.dbId, 0,
			&proto.KeyValue{last.ctx.tableId, nil, nil, nil, 0, 0}}
	} else {
		rec = &last.Reply[len(last.Reply)-1]
	}

	var call = c.goDump(last.ctx.scope, unitId, rec, nil)
	if call.err != nil {
		return nil, call.err
	}

	r, err := (<-call.Done).Reply()
	if err != nil {
		return nil, err
	}

	var s = r.(*DumpReply)
	if s.End {
		return s, nil
	}

	if len(s.Reply) == 0 {
		return c.DumpMore(s)
	} else {
		return s, nil
	}
}

// Get, Set, Del, Incr, ZGet, ZSet, ZDel, ZIncr
func (c *Context) goOneOp(zop bool, args *OneArgs, cmd uint8,
	done chan *Call) *Call {
	call := c.cli.newCall(cmd, done)
	if call.err != nil {
		return call
	}

	var p proto.PkgOneOp
	p.Seq = call.seq
	p.DbId = c.dbId
	p.Cmd = call.cmd
	p.TableId = args.TableId
	p.RowKey = args.RowKey
	p.ColKey = args.ColKey

	if args.Cas != 0 {
		p.Cas = args.Cas
		p.CtrlFlag |= proto.CtrlCas
	}
	if proto.CmdSet == cmd || proto.CmdIncr == cmd {
		if args.Score != 0 {
			p.Score = args.Score
			p.CtrlFlag |= proto.CtrlScore
		}
	}
	if proto.CmdSet == cmd {
		if args.Value != nil {
			p.Value = args.Value
			p.CtrlFlag |= proto.CtrlValue
		}
	}

	// ZGet, ZSet, ZDel, ZIncr
	if zop {
		p.ColSpace = proto.ColSpaceScore1
		p.CtrlFlag |= proto.CtrlColSpace
	}

	var err error
	call.pkg, _, err = p.Encode(nil)
	if err != nil {
		c.cli.errCall(call, err)
		return call
	}

	// put request pkg to sending channel
	c.cli.sending <- call

	return call
}

func (c *Context) GoPing(done chan *Call) *Call {
	var oa OneArgs
	return c.goOneOp(false, &oa, proto.CmdPing, done)
}

func (c *Context) GoGet(tableId uint8, rowKey, colKey []byte, cas uint32,
	done chan *Call) *Call {
	var args = NewGetArgs(tableId, rowKey, colKey, cas)
	return c.goOneOp(false, args, proto.CmdGet, done)
}

func (c *Context) GoZGet(tableId uint8, rowKey, colKey []byte, cas uint32,
	done chan *Call) *Call {
	var args = NewGetArgs(tableId, rowKey, colKey, cas)
	return c.goOneOp(true, args, proto.CmdGet, done)
}

func (c *Context) GoSet(tableId uint8, rowKey, colKey, value []byte, score int64,
	cas uint32, done chan *Call) *Call {
	var args = NewSetArgs(tableId, rowKey, colKey, value, score, cas)
	return c.goOneOp(false, args, proto.CmdSet, done)
}

func (c *Context) GoZSet(tableId uint8, rowKey, colKey, value []byte, score int64,
	cas uint32, done chan *Call) *Call {
	var args = NewSetArgs(tableId, rowKey, colKey, value, score, cas)
	return c.goOneOp(true, args, proto.CmdSet, done)
}

func (c *Context) GoDel(tableId uint8, rowKey, colKey []byte,
	cas uint32, done chan *Call) *Call {
	var args = NewDelArgs(tableId, rowKey, colKey, cas)
	return c.goOneOp(false, args, proto.CmdDel, done)
}

func (c *Context) GoZDel(tableId uint8, rowKey, colKey []byte,
	cas uint32, done chan *Call) *Call {
	var args = NewDelArgs(tableId, rowKey, colKey, cas)
	return c.goOneOp(true, args, proto.CmdDel, done)
}

func (c *Context) GoIncr(tableId uint8, rowKey, colKey []byte, score int64,
	cas uint32, done chan *Call) *Call {
	var args = NewIncrArgs(tableId, rowKey, colKey, score, cas)
	return c.goOneOp(false, args, proto.CmdIncr, done)
}

func (c *Context) GoZIncr(tableId uint8, rowKey, colKey []byte, score int64,
	cas uint32, done chan *Call) *Call {
	var args = NewIncrArgs(tableId, rowKey, colKey, score, cas)
	return c.goOneOp(true, args, proto.CmdIncr, done)
}

// MGet, MSet, MDel, MIncr, ZMGet, ZMSet, ZMDel, ZMIncr
func (c *Context) goMultiOp(zop bool, args *MultiArgs, cmd uint8, done chan *Call) *Call {
	call := c.cli.newCall(cmd, done)
	if call.err != nil {
		return call
	}

	var p proto.PkgMultiOp
	p.Seq = call.seq
	p.DbId = c.dbId
	p.Cmd = call.cmd

	p.Kvs = make([]proto.KeyValueCtrl, len(args.Args))
	for i := 0; i < len(args.Args); i++ {
		p.Kvs[i].TableId = args.Args[i].TableId
		p.Kvs[i].RowKey = args.Args[i].RowKey
		p.Kvs[i].ColKey = args.Args[i].ColKey

		if args.Args[i].Cas != 0 {
			p.Kvs[i].Cas = args.Args[i].Cas
			p.Kvs[i].CtrlFlag |= proto.CtrlCas
		}
		if proto.CmdMSet == cmd || proto.CmdIncr == cmd {
			if args.Args[i].Score != 0 {
				p.Kvs[i].Score = args.Args[i].Score
				p.Kvs[i].CtrlFlag |= proto.CtrlScore
			}
		}
		if proto.CmdMSet == cmd {
			if args.Args[i].Value != nil {
				p.Kvs[i].Value = args.Args[i].Value
				p.Kvs[i].CtrlFlag |= proto.CtrlValue
			}
		}

		// ZMGet, ZMSet, ZMDel, ZMIncr
		if zop {
			p.Kvs[i].ColSpace = proto.ColSpaceScore1
			p.Kvs[i].CtrlFlag |= proto.CtrlColSpace
		}
	}

	var err error
	call.pkg, _, err = p.Encode(nil)
	if err != nil {
		c.cli.errCall(call, err)
		return call
	}

	c.cli.sending <- call

	return call
}

func (c *Context) GoMGet(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(false, args, proto.CmdMGet, done)
}

func (c *Context) GoZmGet(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(true, args, proto.CmdMGet, done)
}

func (c *Context) GoMSet(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(false, args, proto.CmdMSet, done)
}

func (c *Context) GoZmSet(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(true, args, proto.CmdMSet, done)
}

func (c *Context) GoMDel(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(false, args, proto.CmdMDel, done)
}

func (c *Context) GoZmDel(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(true, args, proto.CmdMDel, done)
}

func (c *Context) GoMIncr(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(false, args, proto.CmdMIncr, done)
}

func (c *Context) GoZmIncr(args *MultiArgs, done chan *Call) *Call {
	return c.goMultiOp(true, args, proto.CmdMIncr, done)
}

func (c *Context) goScan(zop bool, tableId uint8, rowKey, colKey []byte,
	score int64, asc, orderByScore bool, num int, done chan *Call) *Call {
	call := c.cli.newCall(proto.CmdScan, done)
	if call.err != nil {
		return call
	}

	if num < 1 {
		c.cli.errCall(call, ErrInvScanNum)
		return call
	}

	var p proto.PkgScanReq
	p.Seq = call.seq
	p.DbId = c.dbId
	p.Cmd = call.cmd
	if asc {
		p.Direction = 0
	} else {
		p.Direction = 1
	}
	p.Num = uint16(num)
	p.TableId = tableId
	p.RowKey = rowKey
	p.ColKey = colKey

	// ZScan
	if zop {
		if orderByScore {
			p.ColSpace = proto.ColSpaceScore1
		} else {
			p.ColSpace = proto.ColSpaceScore2
		}
		p.CtrlFlag |= proto.CtrlColSpace

		if score != 0 {
			p.Score = score
			p.CtrlFlag |= proto.CtrlScore
		}
	}

	var err error
	call.pkg, _, err = p.Encode(nil)
	if err != nil {
		c.cli.errCall(call, err)
		return call
	}

	call.args = scanContext{asc, orderByScore, num}
	c.cli.sending <- call

	return call
}

func (c *Context) GoScan(tableId uint8, rowKey, colKey []byte,
	asc bool, num int, done chan *Call) *Call {
	return c.goScan(false, tableId, rowKey, colKey, 0, asc, false, num, done)
}

func (c *Context) GoZScan(tableId uint8, rowKey, colKey []byte, score int64,
	asc, orderByScore bool, num int, done chan *Call) *Call {
	return c.goScan(true, tableId, rowKey, colKey, score, asc, orderByScore, num,
		done)
}

func (c *Context) goDump(scope uint8, unitId uint16, rec *DumpRecord,
	done chan *Call) *Call {
	call := c.cli.newCall(proto.CmdDump, done)
	if call.err != nil {
		return call
	}

	var p proto.PkgDumpReq
	p.Seq = call.seq
	p.DbId = c.dbId
	p.Cmd = call.cmd
	p.Scope = scope
	p.UnitId = unitId
	p.KeyValue = *rec.KeyValue

	var err error
	call.pkg, _, err = p.Encode(nil)
	if err != nil {
		c.cli.errCall(call, err)
		return call
	}

	call.args = dumpContext{scope, unitId, rec.DbId, rec.TableId}
	c.cli.sending <- call

	return call
}

func (call *Call) Reply() (interface{}, error) {
	if call.err != nil {
		return nil, call.err
	}

	switch call.cmd {
	case proto.CmdPing:
		fallthrough
	case proto.CmdIncr:
		fallthrough
	case proto.CmdDel:
		fallthrough
	case proto.CmdSet:
		fallthrough
	case proto.CmdGet:
		var p proto.PkgOneOp
		_, err := p.Decode(call.pkg)
		if err != nil {
			call.err = err
			return nil, call.err
		}

		if !isNormalErrorCode(p.ErrCode) {
			call.err = newErrCode(p.ErrCode)
			return nil, call.err
		}

		return &OneReply{p.ErrCode, &p.KeyValue}, nil

	case proto.CmdMIncr:
		fallthrough
	case proto.CmdMDel:
		fallthrough
	case proto.CmdMSet:
		fallthrough
	case proto.CmdMGet:
		var p proto.PkgMultiOp
		_, err := p.Decode(call.pkg)
		if err != nil {
			call.err = err
			return nil, call.err
		}

		var r MultiReply
		r.Reply = make([]OneReply, len(p.Kvs))
		for i := 0; i < len(p.Kvs); i++ {
			r.Reply[i].ErrCode = p.Kvs[i].ErrCode
			r.Reply[i].KeyValue = &p.Kvs[i].KeyValue
		}
		return &r, nil

	case proto.CmdScan:
		var p proto.PkgScanResp
		_, err := p.Decode(call.pkg)
		if err != nil {
			call.err = err
			return nil, call.err
		}

		var r ScanReply
		r.ctx = call.args.(scanContext)
		r.End = (p.End != 0)
		r.Reply = make([]OneReply, len(p.Kvs))
		for i := 0; i < len(p.Kvs); i++ {
			r.Reply[i].ErrCode = 0
			r.Reply[i].KeyValue = &p.Kvs[i].KeyValue
		}
		return &r, nil

	case proto.CmdDump:
		var p proto.PkgScanResp
		_, err := p.Decode(call.pkg)
		if err != nil {
			call.err = err
			return nil, call.err
		}

		var r DumpReply
		r.ctx = call.args.(dumpContext)
		r.End = (p.End != 0)
		r.Reply = make([]DumpRecord, len(p.Kvs))
		for i := 0; i < len(p.Kvs); i++ {
			r.Reply[i].DbId = p.Kvs[i].DbIdExt
			r.Reply[i].ColSpace = p.Kvs[i].ColSpace
			r.Reply[i].KeyValue = &p.Kvs[i].KeyValue
		}
		return &r, nil
	}

	return nil, ErrUnknownCmd
}
