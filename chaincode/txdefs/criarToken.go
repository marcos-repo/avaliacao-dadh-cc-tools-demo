package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// POST Method
var CriarToken = tx.Transaction{
	Tag:         "createNewToken",
	Label:       "Criar Novo Token",
	Description: "Criar Novo Token",
	Method:      "POST",
	Callers:     []string{"$org3MSP", "$orgMSP"}, // Only org3 can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "id",
			Label:       "ID",
			Description: "Identificador do Token",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "proprietario",
			Label:       "Proprietario",
			Description: "Proprietario do Token",
			DataType:    "->proprietario",
			Required:    true,
		},
		{
			Tag:         "quantidade",
			Label:       "Quantidade",
			Description: "Quantidade do Token",
			DataType:    "number",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		id, _ := req["id"].(string)
		proprietario, _ := req["proprietario"]
		quantidade, _ := req["quantidade"].(float64)

		if quantidade <= 0 {
			return nil, errors.WrapError(nil, "A quantidade deve ser maior que zero.")
		}

		tokenMap := make(map[string]interface{})
		tokenMap["@assetType"] = "token"
		tokenMap["id"] = id
		tokenMap["proprietario"] = proprietario
		tokenMap["quantidade"] = quantidade

		tokenAsset, err := assets.NewAsset(tokenMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		// Salva o novo Token
		_, err = tokenAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		tokenJSON, nerr := json.Marshal(tokenAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return tokenJSON, nil
	},
}
