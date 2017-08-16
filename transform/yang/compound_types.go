package yang

import "github.com/oshothebig/pbast"

var decimal64 = pbast.NewMessage("Decimal64").
	AddField(pbast.NewMessageField(pbast.Int64, "value", 1)).
	AddField(pbast.NewMessageField(pbast.UInt32, "fraction_digits", 2))

var leafRef = pbast.NewMessage("LeafRef").
	AddField(pbast.NewMessageField(pbast.String, "path", 1))

var identityRef = pbast.NewMessage("IdentityRef").
	AddField(pbast.NewMessageField(pbast.String, "value", 1))

var instanceIdentifier = pbast.NewMessage("InstanceIdentifier").
	AddField(pbast.NewMessageField(pbast.String, "path", 1))
