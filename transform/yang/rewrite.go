package yang

import (
	"strings"

	"github.com/oshothebig/pbast"
)

func LiftMessage(f *pbast.File) *pbast.File {
	types := aggregateMessages(f)
	targets := extractIdenticalMessages(types)
	traversed := map[string]bool{}
	for _, m := range f.Messages {
		liftMessage(m, f, targets, traversed)
	}
	return f
}

func liftMessage(msg *pbast.Message, f *pbast.File, targets stringSet, traversed map[string]bool) {
	if len(msg.Messages) == 0 {
		return
	}

	var children []*pbast.Message
	for _, child := range msg.Messages {
		liftMessage(child, f, targets, traversed)
		if targets.contains(child.Name) {
			if !traversed[child.Name] {
				f.AddMessage(child)
				traversed[child.Name] = true
			}
		} else {
			children = append(children, child)
		}
	}
	msg.Messages = children
}

func extractIdenticalMessages(types map[string][]*pbast.Message) stringSet {
	set := newStringSet()
	for name, msgs := range types {
		if allMessagesIdentical(msgs) {
			set.add(name)
		}
	}

	return set
}

func allMessagesIdentical(msgs []*pbast.Message) bool {
	if len(msgs) == 1 {
		return true
	}

	for i := 0; i < len(msgs); i++ {
		for j := i + 1; j < len(msgs); j++ {
			if !pbast.IsSameType(msgs[i], msgs[j]) {
				return false
			}
		}
	}
	return true
}

func aggregateMessages(f *pbast.File) map[string][]*pbast.Message {
	msgs := map[string][]*pbast.Message{}
	traverseMessages(f.Messages, msgs)
	return msgs
}

func traverseMessages(msgs []*pbast.Message, count map[string][]*pbast.Message) {
	if len(msgs) == 0 {
		return
	}

	head, tail := msgs[0], msgs[1:]
	count[head.Name] = append(count[head.Name], head)
	traverseMessages(head.Messages, count)
	traverseMessages(tail, count)
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
	for _, v := range e.Fields {
		// there is 0, no need to complete
		if v.Index == 0 {
			return e
		}
	}

	field := pbast.NewEnumField("DEFAULT", 0)
	newEnum := *e
	newEnum.Fields = append([]*pbast.EnumField{field}, newEnum.Fields...)
	return &newEnum
}

func AppendPrefixForEnumValueStartingWithNumber(f *pbast.File) *pbast.File {
	if len(f.Messages) == 0 && len(f.Enums) == 0 {
		return f
	}

	newFile := *f
	var enums []*pbast.Enum
	for _, e := range f.Enums {
		enums = append(enums, appendPrefixInEnum(e))
	}
	newFile.Enums = enums

	var messages []*pbast.Message
	for _, m := range f.Messages {
		messages = append(messages, appendPrefixInMessage(m))
	}
	newFile.Messages = messages

	return &newFile
}

func appendPrefixInMessage(m *pbast.Message) *pbast.Message {
	if len(m.Messages) == 0 && len(m.Enums) == 0 {
		return m
	}

	newMessage := *m
	var enums []*pbast.Enum
	for _, e := range m.Enums {
		enums = append(enums, appendPrefixInEnum(e))
	}
	newMessage.Enums = enums

	var messages []*pbast.Message
	for _, m := range m.Messages {
		messages = append(messages, appendPrefixInMessage(m))
	}
	newMessage.Messages = messages

	return &newMessage

}

func appendPrefixInEnum(e *pbast.Enum) *pbast.Enum {
	newEnum := *e
	var fields []*pbast.EnumField
	const prefix = "NUM_"
	for _, f := range e.Fields {
		if strings.IndexAny(f.Name, "0123456789") != 0 {
			fields = append(fields, f)
			continue
		}

		newField := *f
		newField.Name = prefix + newField.Name
		fields = append(fields, &newField)
	}
	newEnum.Fields = fields

	return &newEnum
}
