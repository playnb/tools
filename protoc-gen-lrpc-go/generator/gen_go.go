package generator

import (
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"strings"
)

/*
#Package#	Echo
#package#	echo

#RpcServiceName#	HelloService
#RpcMethodName#		SayHello
#RpcReqName#		HelloRequest
#RpcRespName#		HelloResponse

#MessageName#	NotifyMsg

#Service_Declare#	[TMP_Service_Declare]
#Service_Define#	[TMP_Service_Define]
#Service_Method#	[TMP_Service_Method]
#Notify_Send#		[TMP_Notify_Send]

#Handle_Service_Declare#	[TMP_Handle_Service_Declare]
#Handle_Service_Define#			[TMP_Handle_Service_Define]
#Service_Method_Define		[TMP_Handle_Service_Method_Define]
#NotifyHandles#	[TMP_Handle_Notify_Define]

#TMP_Handle_Notify# [TMP_Handle_Notify]
#TMP_Handle_Service# [TMP_Handle_Service]

#RegNotifyHandle#	[TMP_Reg_Notify]
#RegServiceHandle#	[TMP_Reg_Service]


TMP_Imports

TMP_Call_Define
TMP_Notify_Define

TMP_Handle_Define

TMP_Reg_Define
*/

type PackTemp struct {
	Name string
	name string
}

func (t *PackTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#Package#", t.Name, -1)
	ret = strings.Replace(ret, "#package#", t.name, -1)
	return ret
}
func (t *PackTemp) TrimPackageName(org string) string {
	ret := org
	if strings.HasPrefix(ret, t.name+".") {
		ret = strings.TrimPrefix(ret, t.name+".")
	} else if strings.HasPrefix(ret, "."+t.name+".") {
		ret = strings.TrimPrefix(ret, "."+t.name+".")
	}
	return ret
}

type RpcServiceTemp struct {
	*PackTemp
	RpcServiceName string
}

func (t *RpcServiceTemp) Replace(org string) string {
	ret := t.PackTemp.Replace(org)
	ret = strings.ReplaceAll(ret, "#RpcServiceName#", t.RpcServiceName)
	return ret
}

type RpcMethodTemp struct {
	*RpcServiceTemp
	RpcMethodName string
	RpcReqName    string
	RpcRespName   string
}

func (t *RpcMethodTemp) Replace(org string) string {
	ret := t.RpcServiceTemp.Replace(org)
	ret = strings.Replace(ret, "#RpcMethodName#", t.RpcMethodName, -1)
	ret = strings.Replace(ret, "#RpcReqName#", t.RpcReqName, -1)
	ret = strings.Replace(ret, "#RpcRespName#", t.RpcRespName, -1)
	return ret
}

type MessageTemp struct {
	*PackTemp
	MessageName string
}

func (t *MessageTemp) Replace(org string) string {
	ret := t.PackTemp.Replace(org)
	ret = strings.Replace(ret, "#MessageName#", t.MessageName, -1)
	return ret
}

type GoTemp struct {
	PackTemp *PackTemp
	Keys     map[string]string
}

func (gt *GoTemp) init() {
	gt.Keys = make(map[string]string)
	gt.PackTemp = &PackTemp{}
}
func (gt *GoTemp) Replace(org string) string {
	ret := gt.PackTemp.Replace(org)
	for k, v := range gt.Keys {
		ret = strings.Replace(ret, k, v, -1)
	}
	return ret
}

const TMP_Imports = `
package #package#

import (
	"cell/common/mustang/network"
	"cell/common/mustang/processor"
	//"sync"
	//"sync/atomic"
)

func init() {
}
`

const TMP_Call_Define = `
#Service_Define#
`

const TMP_Service_Declare = `
	#Package##RpcServiceName#() *#RpcServiceName#`

const TMP_Service_Define = `

type I#Package##RpcServiceName# interface {
	//GEN: 生产的代码
	#Package##RpcServiceName#() *#RpcServiceName#
}

type #RpcServiceName# struct {
	processor.CallBase
}
`
const TMP_Service_Method = `
//GEN: 生成#RpcServiceName#的RPC调用"#RpcMethodName#"代码
func (c *#RpcServiceName#) #RpcMethodName#(cmd *#RpcReqName#, callBack func(*#RpcRespName#)) {
	md := &processor.MsgDef{}
	md.MessageID = c.GetCallSequence()
	md.SetRpcRequest()
	md.MessageName = "#package#.#RpcServiceName#.#RpcMethodName#"
	md.Data, _ = cmd.Marshal()
	c.Agent().WriteMsgWithoutPack(md.PackMsg(nil))

	c.AddCallback(md.MessageID, func(def *processor.MsgDef) {
		resp := &#RpcRespName#{}
		resp.Unmarshal(def.Data)
		callBack(resp)
	})
}
`
const TMP_Notify_Define = `
//发送消息的对象
var Send = &_send{}

type _send struct {
}

type I#Package#Notify interface {
	#Package#Notify() *Notify  
}

type Notify struct {
	processor.NotifyBase
}
#Notify_Send#
`
const TMP_Notify_Send = `
//GEN: 通知"#MessageName#"消息的代码
func (*_send) #MessageName#(agents []network.SampleAgent, cmd *#MessageName#) {
	md := &processor.MsgDef{}
	md.MessageName = "#package#.#MessageName#"
	md.Data, _ = cmd.Marshal()
	for _,v := range agents {
		if v != nil {
			v.WriteMsgWithoutPack(md.PackMsg(nil))
		}
	}
}
func (c *Notify) #MessageName#(cmd *#MessageName#) {
	Send.#MessageName#([]network.SampleAgent{c.Agent()}, cmd)
}

func (cmd *#MessageName#) PackMsg() *processor.MsgDef {
	md := &processor.MsgDef{}
	md.MessageName = "#package#.#MessageName#"
	md.Data, _ = cmd.Marshal()
	return md
}
func (cmd *#MessageName#) MsgName() string {
	return "#package#.#MessageName#"
}
`

const TMP_Handle_Define = `
type Handles struct {
	Notify #package#NotifyHandles
#Handle_Service_Declare#
}

type #package#NotifyHandles struct {
	//GEN: 生产的代码
#NotifyHandles#
}

#Handle_Service_Define#
`

const TMP_Handle_Service_Declare = `
	#RpcServiceName#   #package##RpcServiceName#Handles
`
const TMP_Handle_Service_Define = `
type #package##RpcServiceName#Handles struct {
	//GEN: 生产的代码
#Service_Method_Define#
}
`
const TMP_Handle_Notify_Define = `
	#MessageName# func(agent network.SampleAgent, cmd *#MessageName#, clientData interface{})`
const TMP_Handle_Service_Method_Define = `
	#RpcMethodName# func(agent network.SampleAgent, req *#RpcReqName#, clientData interface{}) *#RpcRespName#`

const TMP_Handle_Notify = `
//GEN: 生产的代码
func _notify_#MessageName#(handle *Handles, agent network.SampleAgent, md *processor.MsgDef, clientData interface{}) bool {
	if md.IsNotifyMsg() {
		cmd := &#MessageName#{}
		cmd.Unmarshal(md.Data)
		if handle.Notify.#MessageName# != nil {
			handle.Notify.#MessageName#(agent, cmd, clientData)
			return true
		}
	}
	//TODO 不是RPC的消息，发来一个RPC请求，需要有个出错处理
	return false
}
`
const TMP_Handle_Service = `
func _rpc_#RpcServiceName#_#RpcMethodName#(handle *Handles, agent network.SampleAgent, md *processor.MsgDef, clientData interface{}) bool {
	if md.IsRpcRequest() {
		req := &#RpcReqName#{}
		req.Unmarshal(md.Data)
		md.MessageName = "#package#.#RpcServiceName#.#RpcMethodName#"
		md.SetRpcResponse()
		if handle.#RpcServiceName#.#RpcMethodName# != nil {
			resp := handle.#RpcServiceName#.#RpcMethodName#(agent, req, clientData)
			if resp != nil {
				md.Data, _ = resp.Marshal()
			} else {
				md.Data = nil
				md.ErrorCode = 10002
			}
			agent.WriteMsgWithoutPack(md.PackMsg(nil))
			return true
		} else {
			md.Data = nil
			md.ErrorCode = 10001
		}
	}
	//TODO 是RPC的消息，发来一个Notify请求，需要有个出错处理
	return false
}
`

const TMP_Reg_Define = `
#TMP_Handle_Notify# 

#TMP_Handle_Service# 

type rawHandle func(handle *Handles, agent network.SampleAgent, md *processor.MsgDef, clientData interface{}) bool

func RegProtocolHandle(p interface {
	RegHandle(name string, f processor.HandleFunc)
	GetHandles() *Handles
}) {
	packFunc := func(h *Handles, f rawHandle) processor.HandleFunc {
		return func(agent network.SampleAgent, md *processor.MsgDef, clientData interface{}) bool {
			return f(h, agent, md, clientData)
		}
	}
#RegNotifyHandle#
#RegServiceHandle#
}
`
const TMP_Reg_Notify = `
	p.RegHandle("#package#.#MessageName#",  packFunc(p.GetHandles(), _notify_#MessageName#))`
const TMP_Reg_Service = `
	p.RegHandle("#package#.#RpcServiceName#.#RpcMethodName#",  packFunc(p.GetHandles(), _rpc_#RpcServiceName#_#RpcMethodName#))`

func (g *Generator) generateGo(pack *Package) {
	f := &plugin.CodeGeneratorResponse_File{}
	g.Response.File = append(g.Response.File, f)
	f.Name = proto.String(pack.Name + ".gen.go")

	{
		gt := &GoTemp{}
		gt.init()

		gt.PackTemp.Name = strings.ToUpper(pack.Name[0:1]) + pack.Name[1:]
		gt.PackTemp.name = pack.Name

		/*
		   #Package#	Echo
		   #package#	echo

		   #RpcServiceName#	HelloService
		   #RpcMethodName#		SayHello
		   #RpcReqName#		HelloRequest
		   #RpcRespName#		HelloResponse

		   #MessageName#	NotifyMsg

		   #Service_Declare#	[TMP_Service_Declare]
		   #Service_Define#	[TMP_Service_Define]
		   #Service_Method#	[TMP_Service_Method]

		   #Handle_Service_Declare#	[TMP_Handle_Service_Declare]
		   #Handle_Service_Define#			[TMP_Handle_Service_Define]
		   #Service_Method_Define#		[TMP_Handle_Service_Method_Define]

		   #TMP_Handle_Service# [TMP_Handle_Service]

		   #RegServiceHandle#	[TMP_Reg_Service]


		   TMP_Imports

		   TMP_Call_Define
		   TMP_Notify_Define

		   TMP_Handle_Define

		   TMP_Reg_Define
		*/

		for _, m := range pack.Defines {
			if !strings.HasPrefix(m.GetName(), "Msg") {
				continue
			}
			mt := &MessageTemp{}
			mt.PackTemp = gt.PackTemp
			mt.MessageName = m.GetName()

			gt.Keys["#Notify_Send#"] += mt.Replace(TMP_Notify_Send)
			gt.Keys["#NotifyHandles#"] += mt.Replace(TMP_Handle_Notify_Define)
			gt.Keys["#TMP_Handle_Notify#"] += mt.Replace(TMP_Handle_Notify)
			gt.Keys["#RegNotifyHandle#"] += mt.Replace(TMP_Reg_Notify)

		}
		for _, c := range pack.Service {
			st := &RpcServiceTemp{}
			st.PackTemp = gt.PackTemp
			st.RpcServiceName = c.GetName()
			gt.Keys["#Service_Define#"] += st.Replace(TMP_Service_Define)
			gt.Keys["#Handle_Service_Define#"] += st.Replace(TMP_Handle_Service_Define)
			gt.Keys["#Service_Declare#"] += st.Replace(TMP_Service_Declare)
			gt.Keys["#Handle_Service_Declare#"] += st.Replace(TMP_Handle_Service_Declare)

			for _, m := range c.GetMethod() {
				ct := &RpcMethodTemp{}
				ct.RpcServiceTemp = st
				ct.RpcMethodName = m.GetName()
				ct.RpcReqName = gt.PackTemp.TrimPackageName(m.GetInputType())
				ct.RpcRespName = gt.PackTemp.TrimPackageName(m.GetOutputType())

				gt.Keys["#Service_Method#"] += ct.Replace(TMP_Service_Method)
				gt.Keys["#Service_Method_Define#"] += ct.Replace(TMP_Handle_Service_Method_Define)
				gt.Keys["#TMP_Handle_Service#"] += ct.Replace(TMP_Handle_Service)
				gt.Keys["#RegServiceHandle#"] += ct.Replace(TMP_Reg_Service)
			}
			gt.Keys["#Service_Define#"] += gt.Keys["#Service_Method#"]
			gt.Keys["#Service_Method#"] = ""
		}

		outStr := ""
		outStr += TMP_Imports

		outStr += TMP_Call_Define
		outStr += TMP_Notify_Define

		outStr += TMP_Handle_Define

		outStr += TMP_Reg_Define

		f.Content = proto.String(gt.Replace(gt.Replace(outStr)))
	}

}
