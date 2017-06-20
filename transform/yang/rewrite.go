package yang

import (
	"strings"

	"github.com/oshothebig/pbast"
)

func RemoveRootMessage(f *pbast.File) *pbast.File {
	if len(f.Messages) == 0 {
		return f
	}

	index, root := rootMessage(f)
	if root == nil {
		return f
	}

	messages := make([]*pbast.Message, 0, len(f.Messages)-1)
	messages = append(messages, f.Messages[:index]...)
	messages = append(messages, f.Messages[index+1:]...)

	names := newStringSet()
	for _, x := range messages {
		names.add(x.Name)
	}
	for _, x := range f.Enums {
		names.add(x.Name)
	}
	for _, x := range f.Services {
		names.add(x.Name)
	}

	for _, m := range root.Messages {
		// naming overlap is found, then root message cann't be removed
		if names.contains(m.Name) {
			return f
		}
	}

	if len(root.Fields) != 0 {
		return f
	}

	ret := *f
	if len(root.Messages) == 0 {
		ret.Messages = messages
		return &ret
	}

	ret.Enums = append(ret.Enums, root.Enums...)
	newMessages := make([]*pbast.Message, 0, len(f.Messages)+len(root.Messages)-1)
	newMessages = append(newMessages, f.Messages[:index]...)
	newMessages = append(newMessages, root.Messages...)
	newMessages = append(newMessages, f.Messages[index+1:]...)
	ret.Messages = newMessages

	return &ret
}

func rootMessage(f *pbast.File) (int, *pbast.Message) {
	for i, m := range f.Messages {
		if m.Name == "Root" {
			return i, m
		}
	}

	return -1, nil
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
