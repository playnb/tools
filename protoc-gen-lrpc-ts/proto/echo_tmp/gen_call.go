package echo_tmp

import (
	"cell/common/tools/protoc-gen-lrpc-go/echo"
	"cell/service/h5-gate/protocol/inf"
)

func init() {
}

type IEchoCall interface {
	//GEN: 生产的代码
	EchoCall() *Call
}

type Call struct {
	inf.CallBase
}

//GEN: 生产的代码
func (c *Call) SayHello(cmd *echo.HelloRequest, callBack func(*echo.HelloResponse)) {
	md := &inf.MsgDef{}
	md.MessageID = c.GetCallSequence()
	md.SetRpcRequest()
	md.MessageName = "echo.SayHello"
	md.Data, _ = cmd.Marshal()
	c.Agent().WriteMsgWithoutPack(md.PackMsg(nil))

	c.AddCallback(md.MessageID, func(def *inf.MsgDef) {
		resp := &echo.HelloResponse{}
		resp.Unmarshal(def.Data)
		callBack(resp)
	})
}

//GEN: 生产的代码
func (c *Call) NotifyMsg(cmd *echo.NotifyMsg) {
	md := &inf.MsgDef{}
	md.MessageName = "echo.NotifyMsg"
	md.Data, _ = cmd.Marshal()
	c.Agent().WriteMsgWithoutPack(md.PackMsg(nil))
}
