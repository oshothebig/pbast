package yang

import (
	"github.com/oshothebig/pbast"
)

func RemoveRootMessage(f *pbast.File) *pbast.File {
	if len(f.Messages) == 0 {
		return f
	}

	var ms []*pbast.Message
	for _, m := range f.Messages {
		if m.Name != "Root" {
			ms = append(ms, m)
			continue
		}

		if len(m.Enums) != 0 || len(m.Fields) != 0 {
			ms = append(ms, m)
			continue
		}

		if len(m.Messages) > 1 {
			ms = append(ms, m)
			continue
		}

		if len(m.Messages) == 0 {
			continue
		}

		ms = append(ms, m.Messages[0])
	}

	ret := *f
	ret.Messages = ms
	return &ret
}
