package yang

import "github.com/oshothebig/pbast"

// Message used to convert decimal64 type
var decimal64 = pbast.NewMessage("Decimal64").
	AddField(pbast.NewMessageField(pbast.Int64, "value", 1)).
	AddField(pbast.NewMessageField(pbast.UInt32, "fraction_digits", 2))

// Message used to convert leafref type
var leafRef = pbast.NewMessage("LeafRef").
	AddField(pbast.NewMessageField(pbast.String, "path", 1))

// Message used to convert identityref type
var identityRef = pbast.NewMessage("IdentityRef").
	AddField(pbast.NewMessageField(pbast.String, "value", 1))

// Message used to convert instance-identifier type
var instanceIdentifier = pbast.NewMessage("InstanceIdentifier").
	AddField(pbast.NewMessageField(pbast.String, "path", 1))
