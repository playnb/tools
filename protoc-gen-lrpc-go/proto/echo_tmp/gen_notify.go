package echo_tmp

import (
	"cell/common/mustang/network"
	"cell/common/tools/protoc-gen-lrpc-go/echo"
	"cell/service/h5-gate/protocol/inf"
)

type IEchoNotify interface {
	//GEN: 生产的代码
	EchoNotify() *Notify
}

type Notify struct {
	agent network.IAgent
}

func (n *Notify) Init(agent network.IAgent) {
	n.agent = agent
}

//GEN: 生产的代码
func (n *Notify) NotifyMsg(cmd *echo.NotifyMsg) {
	md := &inf.MsgDef{}
	md.MessageName = "echo.NotifyMsg"
	md.Data, _ = cmd.Marshal()
	n.agent.WriteMsgWithoutPack(md.PackMsg(nil))
}
