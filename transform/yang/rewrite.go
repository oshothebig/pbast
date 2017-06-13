package yang

import (
	"strings"

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

func CompleteZeroInEnum(f *pbast.File) *pbast.File {
	if len(f.Messages) == 0 && len(f.Enums) == 0 {
		return f
	}

	newFile := *f
	var enums []*pbast.Enum
	for _, e := range f.Enums {
		enums = append(enums, completeZeroIfAbsent(e))
	}
	newFile.Enums = enums

	var messages []*pbast.Message
	for _, m := range f.Messages {
		messages = append(messages, completeZeroInMessage(m))
	}
	newFile.Messages = messages

	return &newFile
}

func completeZeroInMessage(m *pbast.Message) *pbast.Message {
	if len(m.Messages) == 0 && len(m.Enums) == 0 {
		return m
	}

	newMessage := *m
	var enums []*pbast.Enum
	for _, e := range m.Enums {
		enums = append(enums, completeZeroIfAbsent(e))
	}
	newMessage.Enums = enums

	var messages []*pbast.Message
	for _, m := range m.Messages {
		messages = append(messages, completeZeroInMessage(m))
	}
	newMessage.Messages = messages

	return &newMessage
}

func completeZeroIfAbsent(e *pbast.Enum) *pbast.Enum {
	const suffix = "_DEFAULT"
	for _, v := range e.Fields {
		// there is 0, no need to complete
		if v.Index == 0 {
			return e
		}
	}

	field := pbast.NewEnumField(strings.ToUpper(e.Name)+suffix, 0)
	newEnum := *e
	newEnum.Fields = append([]*pbast.EnumField{field}, newEnum.Fields...)
	return &newEnum
}
