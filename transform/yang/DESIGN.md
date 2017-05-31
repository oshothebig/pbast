# Conversion from YANG to Protocol Buffers
This documents describes how to convert YANG ([RFC6020](https://tools.ietf.org/html/rfc6020)) model definition to
[Protocol Buffers](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec) definition.

## Analyses on YANG statements
### module statement
A `module` statement defines a single data model that essentially works as a namespace or package
in modern programming languages. Actually, a `module` statement always has a `namespace` sub-statement.

`import` statements and `include` statements indicates that the module uses constructs defined in
other files. The former is for constructs defined in another module and the latter is for constructs
defined in an own submodule.

`organization`, `contact`, `description`, `reference` are categorized as meta information on the module.
These would be suitable to be expressed as comments in other languages if there is no equivalent.

While `revision` would be able to be treated as meta information like the above as well, it provides
additional information on versioning or dependency management of the module. The information could be
useful to build a package/dependency management tool which we can see in other languages.

### submodule statement
A `submodule` statement is almost equal to a `module` statement except that a submodule belongs to the
parent module. That means a submodule shares the same namespace as its parent module. A `submodule`
statement doesn't have a `namespace` sub-statement.

`import`, `include`, `organization`, `contact`, `description`, `reference` and `revision` statement
works the same as those in a module.

### container statment
A `container` statement expresses a data node in the data model tree. It works like as a class or struct
in other languages. In addition, a `container` statement implies a field of the parent node that the
container belongs to. The name of a container almost equally works as a type name and a field name as well.

`container`, `leaf`, `list`, `leaf-list`, `uses`, `choice`, and `anyxml` can be used to define child
nodes of the container. At this time, `choice` and `anyxml` are not considered for simplicity.

### leaf statement
A `leaf` statement expresses a value and its type. It works like as a field in a class or struct in other
languages. The name of a list almost equally works as a field name of the specified type.

If a `default` value is defined in the `leaf` statement, the definition should be converted as much as
possible in other languages.

### list statement
A `list` statement expresses a data node in the data model tree. While it works like
as a class or struct in other language similar to `container`, its entry can be specified by the `key`
if it is defined. A `list` can be considered as a data structure like map or dictionary in other
languages if key(s) are defined.

### leaf-list statement
A `leaf-list` statement expresses an array of a particular type. It works like a field of an array or
a list in a class/struct in other languages. The name of a leaf-list alsmost equally works as a field
name of the specified type.

### grouping statement
A `grouping` statement is used to define a reusable block of nodes. It works like a definition
of a class or struct in other language, but the child nodes are directly extracted where `uses` statement
is used.

### uses statement
A `uses` statement is used to reference a `grouping` definition. It extracts the child nodes of the
specified `grouping` where it appears.

### rpc statement
A `rpc` statement defines a RPC (Remote Procedure Call) operation. It is defined directly under a module or
a submodule. The argument of a `rpc` statement is like a method or function name in other languages.
A `input` sub-statement or `output` sub-statement may be defined under a `rpc` statement. They define
an input parameters or an output return values in other languages. Both `input` statement and `output`
statement may have child nodes, then it might be better to treat the statement as a new type definition
like input as a new "request" type and output as a new "response" type.

### notification statement
A `notification` statement defines a notification. It can't be directly mapped to a method or function
in other languages, but it can be thought as a function returning streaming values asynchronously.
It may have child nodes, then it might be better to trea the statement as a new type definition as the
return type whose values are asynchronously sent.

### Naming conventions
We can see the following typical naming conventions for arguments or values of constructs in YANG
although they would not cover all cases.

- Dash separated string
    - Example: "openconfig-bgp"
    - Seen in `module`, `submodule`, `import`, `include`, `container`, `leaf`, `list`, `leaf-list`, `uses`, `grouping`, `rpc`, `notification` in IETF data models
- Colon separeted string
    - Example: "urn:ietf:params:xml:ns:yang:iana-if-type"
    - Seen in `namespace` defined as IETF's data models
- URL like string
    - Example: "http://brocade.com/ns/interfacedatatype"
    - Seen in `namespace` in vendor specific data models and OpenConfig data models
- CamelCase
    - Example: "CompanyService"
    - Seen in `module` in vendor specific data models, but could seen in other constructs (although not fully investigated)
- camelCase
    - Example: "boardOfDirector"
    - Seen in vendor specific data models (not fully investigated)
- underscore_case
    - Example: "interface_priority"
    - Seen in vendor specific data models (not fully investigated)
- Unformatted string
    - Example: "groupmanager"
    - Seen in vendor specific data models (not fully investigated)

## Conversion guideline
### module
- A `module` is converted to a ".proto" file
- The package of the proto file is generated from the module's `namespace`
- The argument of the `namespace` should be transformed to a dot seperated string to use it as the package name

### submodule
- A `submodule` is converted to a ".proto" file similar to a `module` if the `submodule` is the top level node for conversion
- The package of the proto file is the same as the package of the `module` which the `submodule` belongs to
- The package name is generated in the same manner as the package name for a module

### container
- A `container` is converted to a `message` definition in Protocol Buffers and a field of the parent node
- The name of the `message`, which should be CamelCase, is generated from the argument of the `container`
- The name of the field, which should be underscore_case, is generated from the argument of the `container`

### leaf
- A `leaf` is converted to a field of the parent node which the `leaf` belongs to
- The type of the field is that specified in the `leaf` statement
- The name of the field, which should be underscore_caes, is generated from the argument of the `leaf`

### list
- Guideline for a `leaf` is the same as that for a `container`

### leaf-list
- A `leaf-list` is converted to a field of the parent node wich the `leaf` belongs to
- The field is marked as "repeated"
- The type of the field is that specified in the `leaf-list` statement
- The name of the field, which should be underscore_caes, is generated from the argument of the `leaf-list`

### grouping
- The child nodes of a `grouping` are extracted where the parent node that `uses` is used
- A `grouping` itself can be ignored for defining a new message type or a new field

### rpc
- A `rpc` is converted to a `rpc` definition in Protocol Buffers
- The name of the `rpc` in Protocol Buffers, which should be CamelCase, is generated from the argument of the `rpc` statement
- The name of the `service` enclosing the `rpc` definition, which should be CamelCase, is generated from the `module` name enclosing the `rpc` statement
- `input` sub-statement is converted to the incoming message type of the `rpc` in Protocol Buffers
    - The name of the message type ends with "Request"
- `output` sub-statement is converted to the returning message type of the `rpc` in Protocol Buffers
    - The name of the message type ends with "Response"
- If `input` or `output` sub-statement is not defined, the corresponding message type is [Empty](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Empty)

### notification
- A `notification` is converted to a `rpc` definition
- The name of the `rpc`, which should be CamelCase, is generated from the argument of the `notification` statement
- The name of the `service` enclosing the `rpc` definition, which should be CamelCase, is generated from the `module` name enclosing the `notification` statement
- The incoming message type is [Empty](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#google.protobuf.Empty)
- The notification message type is defined from the `notification` statement

### Type mapping
The following table shows how YANG built-in type is mapped to Protocol Buffers' built-in type.

YANG type | Protocol Buffer type
----------|----------------------
int8      |int32
int16     |int32
int32     |int32
int64     |int64
uint8     |uint32
uint16    |unit32
uint32    |uint32
uint64    |uint64
decimal64 |N/A
string    |string
boolean   |bool
enum      |enum
bits      |N/A
binary    |bytes
leafref   |N/A
identityref|N/A
empty     |N/A
union     |N/A
instance-identifier|N/A

Some of YANG types, such as int32 and uint32, completely fit to Protocol Buffers' types.
While some types like int8 and uint8 in YANG don't have a direct corresponding types in Protocol Buffers,
we can substitute a larger interger type for those types. The other types are completely YANG unique types
(shown as N/A) and we would need to think about how to compose it from Protocol Buffers' built-in types.