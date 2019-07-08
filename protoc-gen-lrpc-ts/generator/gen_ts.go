package generator

import (
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"strings"
)

//--------------------------------------
type PackageTemp struct {
	notify   *NotifyTemp
	service  []*RpcServiceTemp
	handle   *HandleTemp
	PackName string
}

func (pt *PackageTemp) String() string {
	ret := `
namespace #PackName#.net {`
	ret += pt.notify.String()
	for _, m := range pt.service {
		ret += m.String()
	}
	ret += pt.handle.String()
	ret += `
}
`
	ret = strings.Replace(ret, "#PackName#", pt.PackName, -1)
	return ret
}

//--------------------------------------
type HandleTemp struct {
	notify  []*HandleNotifyTemp
	service []*HandleServiceTemp
}

func (ht *HandleTemp) String() string {
	ret := `
    export class Handle extends processor.HandleBase {
        constructor() {
            super()
        }
		//生成代码部分
        protected init() {`
	for _, v := range ht.notify {
		ret += v.Declare()
	}
	for _, v := range ht.service {
		ret += v.Declare()
	}
	ret += `
		}
`
	for _, v := range ht.notify {
		ret += v.Define()
	}
	for _, v := range ht.service {
		ret += v.Define()
	}
	ret += `
	}
`
	return ret
}

//--------------------------------------

type HandleNotifyTemp struct {
	MessageName string
}

func (hnt *HandleNotifyTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#MessageName#", hnt.MessageName, -1)
	return ret
}
func (hnt *HandleNotifyTemp) Declare() string {
	ret := `
			this.codecs.set("#PackName#.#MessageName#", { decode: (c) => { return #PackName#.#MessageName#.decode(c) }, encode: (c) => { return #PackName#.#MessageName#.encode(c).finish() } })
`
	return hnt.Replace(ret)
}

func (hnt *HandleNotifyTemp) Define() string {
	ret := `
        public #MessageName#(f: processor.Func<#PackName#.#MessageName#, void>) {
            this.handles.set("#PackName#.#MessageName#", {
                action: f,
                func: undefined,
                reqCodec: this.codecs.get("#PackName#.#MessageName#"),
                respCodec: undefined
            })
        }
`
	return hnt.Replace(ret)
}

//--------------------------------------
type HandleServiceTemp struct {
	methods     []*HandleServiceMethodTemp
	ServiceName string
}

func (hst *HandleServiceTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#ServiceName#", hst.ServiceName, -1)
	return ret
}
func (hst *HandleServiceTemp) Declare() string {
	ret := ``
	for _, m := range hst.methods {
		ret += m.Declare()
	}
	return hst.Replace(ret)
}
func (hst *HandleServiceTemp) Define() string {
	ret := ``
	for _, m := range hst.methods {
		ret += m.Define()
	}
	return hst.Replace(ret)
}

//--------------------------------------

type HandleServiceMethodTemp struct {
	MethodName string
	InName     string
	OutName    string
}

func (hst *HandleServiceMethodTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#MethodName#", hst.MethodName, -1)
	ret = strings.Replace(ret, "#InName#", hst.InName, -1)
	ret = strings.Replace(ret, "#OutName#", hst.OutName, -1)
	return ret
}

func (hst *HandleServiceMethodTemp) Declare() string {
	ret := `
            this.codecs.set("#PackName#.#ServiceName#.#MethodName#", { decode: (c) => { return #OutName#.decode(c) }, encode: (c) => { return #InName#.encode(c).finish() } })
`
	return hst.Replace(ret)
}

func (hst *HandleServiceMethodTemp) Define() string {
	ret := `
        public #ServiceName#_#MethodName#(f: processor.Func<#InName#, #OutName#>) {
            this.handles.set("#PackName#.#ServiceName#.#MethodName#", {
                action: undefined,
                func: f,
                reqCodec: this.codecs.get("#InName#"),
                respCodec: this.codecs.get("#OutName#")
            })
        }
`
	return hst.Replace(ret)
}

//--------------------------------------
type NotifyTemp struct {
	methods []*NotifyMethodTemp
}

func (nt *NotifyTemp) Replace(org string) string {
	return org
}

func (nt *NotifyTemp) String() string {
	ret := `
    //通知类消息
    export class Notify extends processor.NotifyBase {
        constructor(net: processor.Network) {
            super(net)
        }
`
	for _, m := range nt.methods {
		ret += m.String()
	}
	ret += `
    }
`
	return nt.Replace(ret)
}

//--------------------------------------
type NotifyMethodTemp struct {
	MessageName string
}

func (nmt *NotifyMethodTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#MessageName#", nmt.MessageName, -1)
	return ret
}
func (nmt *NotifyMethodTemp) String() string {
	ret := `
		//如果发送失败，则返回false
        public #MessageName#(cmd: #PackName#.#MessageName#): boolean {
            return this.sendMsg("#PackName#.#MessageName#", #PackName#.#MessageName#.encode(cmd).finish())
        }
`
	return nmt.Replace(ret)
}

//--------------------------------------
type RpcServiceTemp struct {
	ServiceName string
	methods     []*RpcServiceMethodTemp
}

func (rst *RpcServiceTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#ServiceName#", rst.ServiceName, -1)
	return ret
}

func (rst *RpcServiceTemp) String() string {
	ret := `
    //函数类消息
    export class #ServiceName# extends processor.CallBase {
        //TODO: 这里只是成功的情况，还有各种失败的情况没有处理(比如超时)
        constructor(handle: processor.IRpcSub, net: processor.Network) {
            super(handle, net)
        }
`
	for _, m := range rst.methods {
		ret += m.String()
	}
	ret += `
    }
`
	return rst.Replace(ret)
}

//--------------------------------------
type RpcServiceMethodTemp struct {
	MethodName string
	InName     string
	OutName    string
}

func (rsmt *RpcServiceMethodTemp) Replace(org string) string {
	ret := org
	ret = strings.Replace(ret, "#MethodName#", rsmt.MethodName, -1)
	ret = strings.Replace(ret, "#InName#", rsmt.InName, -1)
	ret = strings.Replace(ret, "#OutName#", rsmt.OutName, -1)
	return ret
}

func (rsmt *RpcServiceMethodTemp) String() string {
	ret := `
        //生成代码部分
        public #MethodName#(cmd: #InName#): Promise<#OutName#> {
            return new Promise<#OutName#>((resolve, reject) => {
                try {
                    let md = this.makeMsg("#PackName#.#ServiceName#.#MethodName#", #InName#.encode(cmd).finish())
                    md.RpcRequest = true
                    if (!this.net.write(md.packMsg())) {
                        reject(new processor.NetError(10000));        // 失败                    
                        return
                    }
                    this.processor.addRpcSub(md.MessageID, resolve, reject)
                }
                catch (error) {
                    reject(new processor.NetError(error))        // 失败
                }
            })
        }`

	return rsmt.Replace(ret)
}

//------------------------------------------------------------------------------------------------------------------
func (g *Generator) generateTS(pack *Package) {
	f := &plugin.CodeGeneratorResponse_File{}
	g.Response.File = append(g.Response.File, f)
	f.Name = proto.String(pack.Name + ".gen.ts")

	outStr := ""

	{
		pt := &PackageTemp{}
		pt.PackName = pack.Name

		pt.notify = &NotifyTemp{}
		pt.handle = &HandleTemp{}
		for _, m := range pack.Defines {
			if !strings.HasPrefix(m.GetName(), "Msg") {
				continue
			}
			
			nmt := &NotifyMethodTemp{}
			nmt.MessageName = m.GetName()
			pt.notify.methods = append(pt.notify.methods, nmt)

			hnt := &HandleNotifyTemp{}
			hnt.MessageName = m.GetName()
			pt.handle.notify = append(pt.handle.notify, hnt)
		}
		for _, c := range pack.Service {
			st := &RpcServiceTemp{}
			pt.service = append(pt.service, st)
			st.ServiceName = c.GetName()

			hst := &HandleServiceTemp{}
			hst.ServiceName = c.GetName()
			pt.handle.service = append(pt.handle.service, hst)

			for _, m := range c.GetMethod() {
				rmt := &RpcServiceMethodTemp{}
				st.methods = append(st.methods, rmt)
				rmt.MethodName = m.GetName()
				rmt.InName = strings.TrimPrefix(m.GetInputType(), ".")
				rmt.OutName = strings.TrimPrefix(m.GetOutputType(), ".")

				hsmt := &HandleServiceMethodTemp{}
				hsmt.MethodName = m.GetName()
				hsmt.InName = strings.TrimPrefix(m.GetInputType(), ".")
				hsmt.OutName = strings.TrimPrefix(m.GetOutputType(), ".")
				hst.methods = append(hst.methods, hsmt)
			}
		}

		outStr = pt.String()
	}

	f.Content = proto.String(outStr)
}
