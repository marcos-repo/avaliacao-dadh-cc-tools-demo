package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

var Token = assets.AssetType{
	Tag:         "token",
	Label:       "Token",
	Description: "Token Laboratorio",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "ID",
			DataType: "string",                      // Datatypes are identified at datatypes folder
			Writers:  []string{`org1MSP`, "orgMSP"}, // This means only org1 can create the asset (others can edit)
		},
		{
			// Mandatory property
			Required: true,
			Tag:      "proprietario",
			Label:    "Proprietario",
			DataType: "->proprietario",
		},
		{
			// Optional property
			Tag:          "quantidade",
			Label:        "Quantidade",
			DataType:     "number",
			DefaultValue: 0,
		},
		{
			// Optional property
			Tag:          "burned",
			Label:        "Burned",
			DataType:     "boolean",
			DefaultValue: false,
		},
	},
}
