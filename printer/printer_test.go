package pbast

import (
	"bytes"
	"testing"

	"github.com/oshothebig/pbast"
)

var syntax pbast.Syntax

var table = []struct {
	in       pbast.Node
	expected string
}{
	{
		syntax,
		`syntax = "proto3";
`,
	},
	{
		pbast.NewImport("org.foo"),
		`import "org.foo";
`,
	},
	{
		pbast.NewWeakImport("org.foo"),
		`import weak "org.foo";
`,
	},
	{
		pbast.NewPublicImport("org.foo"),
		`import public "org.foo";
`,
	},
	{
		pbast.Package{Name: "org.foo"},
		"package org.foo;\n",
	},
	{
		pbast.NewOption("human", "men"),
		"human = men;\n",
	},
	{
		pbast.NewMessage("human").
			AddField(pbast.NewMessageField("string", "firstName", 1)).
			AddField(pbast.NewMessageField("string", "lastName", 2)),
		`message human {
  string firstName = 1;
  string lastName = 2;
}
`,
	},
	{
		pbast.NewMessage("human").
			AddField(pbast.NewMessageField("string", "name", 1).
				AddOption(pbast.NewFieldOption("sex", "male"))),
		`message human {
  string name = 1 [sex = male];
}
`,
	},
	{
		pbast.NewMessageField("string", "name", 0),
		"string name = 0;\n",
	},
	{
		pbast.NewRepeatedMessageField("string", "name", 0),
		"repeated string name = 0;\n",
	},
	{
		pbast.NewMessageField("string", "name", 0).AddOption(pbast.NewFieldOption("age", "21")),
		"string name = 0 [age = 21];\n",
	},
	{
		pbast.NewMessageField("string", "name", 0).
			AddOption(pbast.NewFieldOption("age", "21")).
			AddOption(pbast.NewFieldOption("tall", "170")),
		"string name = 0 [age = 21, tall = 170];\n",
	},
	{
		pbast.NewEnum("sex").
			AddField(pbast.NewEnumField("male", 1)).
			AddField(pbast.NewEnumField("female", 2)),
		`enum sex {
  male = 1;
  female = 2;
}
`,
	},
	{
		pbast.NewEnumField("male", 1),
		"male = 1;\n",
	},
	{
		pbast.NewEnumField("male", 1).
			AddOption(pbast.NewEnumValueOption("age", "11")),
		"male = 1 [age = 11];\n",
	},
	{
		pbast.NewEnumField("male", 1).
			AddOption(pbast.NewEnumValueOption("age", "11")).
			AddOption(pbast.NewEnumValueOption("type", "human")),
		"male = 1 [age = 11, type = human];\n",
	},
	{
		pbast.NewService("get").
			AddRPC(pbast.NewRPC("name", pbast.NewReturnType("string"), pbast.NewReturnType("int"))),
		`service get {
  rpc name (string) returns (int);
}
`,
	},
	{
		pbast.NewService("get").
			AddRPC(pbast.NewRPC("name", pbast.NewReturnType("string"), pbast.NewReturnType("int"))).
			AddRPC(pbast.NewRPC("age", pbast.NewReturnType("string"), pbast.NewReturnType("int"))),
		`service get {
  rpc name (string) returns (int);
  rpc age (string) returns (int);
}
`,
	},
	{
		pbast.NewFile().
			AddImport(pbast.NewImport("org.example")).
			AddPackage(pbast.NewPackage("org.foo")).
			AddMessage(pbast.NewMessage("human").
				AddField(pbast.NewMessageField("string", "firstName", 1)).
				AddField(pbast.NewMessageField("string", "lastName", 2))).
			AddMessage(pbast.NewMessage("animal").
				AddField(pbast.NewMessageField("string", "name", 1)).
				AddField(pbast.NewMessageField("int", "age", 2))).
			AddEnum(pbast.NewEnum("sex").
				AddField(pbast.NewEnumField("male", 1)).
				AddField(pbast.NewEnumField("female", 2))).
			AddService(pbast.NewService("get").
				AddRPC(pbast.NewRPC("name", pbast.NewReturnType("string"), pbast.NewReturnType("int"))).
				AddRPC(pbast.NewRPC("age", pbast.NewReturnType("string"), pbast.NewReturnType("int")))),
		`syntax = "proto3";
import "org.example";
package org.foo;

message human {
  string firstName = 1;
  string lastName = 2;
}

message animal {
  string name = 1;
  int age = 2;
}

enum sex {
  male = 1;
  female = 2;
}

service get {
  rpc name (string) returns (int);
  rpc age (string) returns (int);
}
`,
	},
}

func TestWrite(t *testing.T) {
	for x, d := range table {
		buf := new(bytes.Buffer)
		Fprint(buf, d.in)

		if bytes.Compare(buf.Bytes(), []byte(d.expected)) != 0 {
			t.Errorf("#%d:\ngot\n%s\nwant\n%s", x, buf.Bytes(), d.expected)
		}
	}
}

func BenchmarkWrite(b *testing.B) {
	w := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		for _, d := range table {
			w.Reset()

			Fprint(w, d.in)
		}
	}
}
