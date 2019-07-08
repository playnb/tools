package echo_tmp

import (
	"cell/common/mustang/network"
	"cell/common/tools/protoc-gen-lrpc-go/echo"
	"cell/service/h5-gate/protocol/inf"
)

var Handles = &echoHandles{}

type echoHandles struct {
	Notify echoNotifyHandles
	Call   echoCallHandles
}

////////////////////////////////////////////////////////////////////////////////
//GEN: 生产的代码
type echoNotifyHandles struct {
	NotifyMsg func(agent network.IAgent, cmd *echo.NotifyMsg)
}

//GEN: 生产的代码
type echoCallHandles struct {
	SayHello func(agent network.IAgent, cmd *echo.HelloRequest) *echo.HelloResponse
}

////////////////////////////////////////////////////////////////////////////////
//GEN: 生产的代码
func _notify_NotifyMsg(agent network.IAgent, md *inf.MsgDef) {
	if md.IsNotifyMsg() {
		cmd := &echo.NotifyMsg{}
		cmd.Unmarshal(md.Data)
		if Handles.Notify.NotifyMsg != nil {
			Handles.Notify.NotifyMsg(agent, cmd)
		}
		return
	}
	//TODO 不是RPC的消息，发来一个RPC请求，需要有个出错处理
}

func _rpc_SayHello(agent network.IAgent, md *inf.MsgDef) {
	if md.IsRpcRequest() {
		req := &echo.HelloRequest{}
		req.Unmarshal(md.Data)
		md.MessageName = "echo.SayHello"
		md.SetRpcResponse()
		if Handles.Call.SayHello != nil {
			resp := Handles.Call.SayHello(agent, req)
			if resp != nil {
				md.Data, _ = resp.Marshal()
			} else {
				md.Data = nil
				md.ErrorCode = 10002
			}
		} else {
			md.Data = nil
			md.ErrorCode = 10001
		}
		agent.WriteMsgWithoutPack(md.PackMsg(nil))

		return
	}
	//TODO 是RPC的消息，发来一个Notify请求，需要有个出错处理
}

////////////////////////////////////////////////////////////////////////////////
//GEN: 生产的代码
func RegEchoHandle(p interface{ RegHandle(name string, f inf.HandleFunc) }) {
	p.RegHandle("echo.NotifyMsg", _notify_NotifyMsg)
	p.RegHandle("echo.SayHello", _rpc_SayHello)
}
