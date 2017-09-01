package pbast

import (
	"reflect"
	"testing"
)

func TestLiftMessage(t *testing.T) {
	table := []struct {
		in       *File
		expected *File
	}{
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "A",
						Messages: []*Message{
							&Message{
								Name: "B",
							},
						},
					},
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
					&Message{Name: "B"},
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "A",
						Messages: []*Message{
							&Message{
								Name: "B",
								Messages: []*Message{
									&Message{Name: "C"},
								},
							},
							&Message{Name: "D"},
						},
					},
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
					&Message{Name: "B"},
					&Message{Name: "D"},
					&Message{Name: "C"},
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("C"))).
						AddMessage(NewMessage("D")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B"),
					NewMessage("D"),
					NewMessage("C"),
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B"))).
						AddMessage(NewMessage("C")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").
						AddMessage(NewMessage("B")),
					NewMessage("C"),
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B").
								AddMessage(NewMessage("C")))).
						AddMessage(NewMessage("D")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").AddMessage(NewMessage("B")),
					NewMessage("D"),
					NewMessage("C"),
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B").
								AddMessage(NewMessage("C")))).
						AddMessage(NewMessage("C")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").
						AddMessage(NewMessage("B")).
						AddMessage(NewMessage("C")),
					NewMessage("C")},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B").
								AddMessage(NewMessage("C").
									AddMessage(NewMessage("D"))))).
						AddMessage(NewMessage("C")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").
						AddMessage(NewMessage("B")).
						AddMessage(NewMessage("C")),
					NewMessage("C"),
					NewMessage("D"),
				},
			},
		},
	}

	for x, d := range table {
		if actual := LiftMessage(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}
