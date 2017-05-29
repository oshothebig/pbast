package pbast

type Service struct {
	Name    string
	Options []Option
	RPCs    []RPC
}

func NewService(name string) Service {
	return Service{
		Name: name,
	}
}

func (s Service) AddOptions(o Option) Service {
	ns := Service(s)
	s.Options = append(s.Options, o)
	return ns
}

func (s Service) AddRPC(r RPC) Service {
	ns := Service(s)
	ns.RPCs = append(s.RPCs, r)
	return ns
}

type RPC struct {
	Name    string
	Input   ReturnType
	Output  ReturnType
	Options []Option
}

func NewRPC(name string, input ReturnType, output ReturnType) RPC {
	return RPC{
		Name:   name,
		Input:  input,
		Output: output,
	}
}

func (r RPC) AddOption(o Option) RPC {
	nr := RPC(r)
	nr.Options = append(nr.Options, o)
	return nr
}

type ReturnType struct {
	Name       string
	Streamable bool
}

func NewReturnType(name string) ReturnType {
	return ReturnType{
		Name: name,
	}
}
