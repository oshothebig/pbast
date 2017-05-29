package pbast

import "testing"
import "bytes"

var syntax Syntax

var table = []struct {
	in       Node
	expected string
}{
	{
		syntax,
		`syntax = "proto3";
`,
	},
	{
		NewImport("org.foo"),
		`import "org.foo";
`,
	},
	{
		NewWeakImport("org.foo"),
		`import weak "org.foo";
`,
	},
	{
		NewPublicImport("org.foo"),
		`import public "org.foo";
`,
	},
	{
		Package{Name: "org.foo"},
		"package org.foo;\n",
	},
	{
		NewOption("human", "men"),
		"human = men;\n",
	},
	{
		NewMessage("human").
			AddField(NewMessageField("string", "firstName", 1)).
			AddField(NewMessageField("string", "lastName", 2)),
		`message human {
  string firstName = 1;
  string lastName = 2;
}
`,
	},
	{
		NewMessage("human").
			AddField(NewMessageField("string", "name", 1).
				AddOption(NewFieldOption("sex", "male"))),
		`message human {
  string name = 1 [sex = male];
}
`,
	},
	{
		NewMessageField("string", "name", 0),
		"string name = 0;\n",
	},
	{
		NewRepeatedMessageField("string", "name", 0),
		"repeated string name = 0;\n",
	},
	{
		NewMessageField("string", "name", 0).AddOption(NewFieldOption("age", "21")),
		"string name = 0 [age = 21];\n",
	},
	{
		NewMessageField("string", "name", 0).
			AddOption(NewFieldOption("age", "21")).
			AddOption(NewFieldOption("tall", "170")),
		"string name = 0 [age = 21, tall = 170];\n",
	},
	{
		NewEnum("sex").
			AddField(NewEnumField("male", 1)).
			AddField(NewEnumField("female", 2)),
		`enum sex {
  male = 1;
  female = 2;
}
`,
	},
	{
		NewEnumField("male", 1),
		"male = 1;\n",
	},
	{
		NewEnumField("male", 1).
			AddOption(NewEnumValueOption("age", "11")),
		"male = 1 [age = 11];\n",
	},
	{
		NewEnumField("male", 1).
			AddOption(NewEnumValueOption("age", "11")).
			AddOption(NewEnumValueOption("type", "human")),
		"male = 1 [age = 11, type = human];\n",
	},
	{
		NewService("get").
			AddRPC(NewRPC("name", NewReturnType("string"), NewReturnType("int"))),
		`service get {
  rpc name (string) returns (int);
}
`,
	},
	{
		NewService("get").
			AddRPC(NewRPC("name", NewReturnType("string"), NewReturnType("int"))).
			AddRPC(NewRPC("age", NewReturnType("string"), NewReturnType("int"))),
		`service get {
  rpc name (string) returns (int);
  rpc age (string) returns (int);
}
`,
	},
	{
		NewFile().
			AddImport(NewImport("org.example")).
			AddPackage(NewPackage("org.foo")).
			AddMessage(NewMessage("human").
				AddField(NewMessageField("string", "firstName", 1)).
				AddField(NewMessageField("string", "lastName", 2))).
			AddMessage(NewMessage("animal").
				AddField(NewMessageField("string", "name", 1)).
				AddField(NewMessageField("int", "age", 2))).
			AddEnum(NewEnum("sex").
				AddField(NewEnumField("male", 1)).
				AddField(NewEnumField("female", 2))).
			AddService(NewService("get").
				AddRPC(NewRPC("name", NewReturnType("string"), NewReturnType("int"))).
				AddRPC(NewRPC("age", NewReturnType("string"), NewReturnType("int")))),
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
		p := Printer{}
		p.Fprint(buf, d.in)

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
			p := Printer{}

			p.Fprint(w, d.in)
		}
	}
}
