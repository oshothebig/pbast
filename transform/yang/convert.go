package yang

import (
	"fmt"

	"github.com/openconfig/goyang/pkg/yang"
	"github.com/oshothebig/pbast"
)

// e must be YANG module
func Transform(e *yang.Entry) pbast.File {
	if _, ok := e.Node.(*yang.Module); !ok {
		return pbast.File{}
	}

	return transformModule(entry{e})
}

func transformModule(e entry) pbast.File {
	namespace := e.Namespace().Name
	elems := guessElements(namespace)
	fmt.Println(elems)
	f := pbast.NewFile(pbast.NewPackageWithElements(elems))

	f = f.AddService(transformRPCs(e))

	return f
}

func transformRPCs(e entry) pbast.Service {
	rpcs := e.rpcs()
	if len(rpcs) == 0 {
		return pbast.Service{}
	}

	s := pbast.NewService(CamelCase(e.Name))
	for _, rpc := range rpcs {
		s = s.AddRPC(transformRPC(rpc))
	}

	return s
}

func transformRPC(e entry) pbast.RPC {
	method := CamelCase(e.Name)

	return pbast.NewRPC(
		method,
		pbast.NewReturnType(method+"Request"),
		pbast.NewReturnType(method+"Response"),
	)
}
