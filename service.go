package pbast

type Service struct {
	Name    string
	Options []*Option
	RPCs    []*RPC
}

func NewService(name string) *Service {
	return &Service{
		Name: name,
	}
}

func (s *Service) AddOptions(o *Option) *Service {
	if o == nil {
		return s
	}
	s.Options = append(s.Options, o)
	return s
}

func (s *Service) AddRPC(r *RPC) *Service {
	if r == nil {
		return s
	}
	s.RPCs = append(s.RPCs, r)
	return s
}

type RPC struct {
	Name    string
	Input   *ReturnType
	Output  *ReturnType
	Options []*Option
}

func NewRPC(name string, input *ReturnType, output *ReturnType) *RPC {
	return &RPC{
		Name:   name,
		Input:  input,
		Output: output,
	}
}

func (r *RPC) AddOption(o *Option) *RPC {
	if o == nil {
		return r
	}
	r.Options = append(r.Options, o)
	return r
}

type ReturnType struct {
	Name       string
	Streamable bool
}

func NewReturnType(name string) *ReturnType {
	return &ReturnType{
		Name: name,
	}
}
