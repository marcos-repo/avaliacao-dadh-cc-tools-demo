package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// Create a new Proprietario on channel
// POST Method
var CriarProprietario = tx.Transaction{
	Tag:         "criarProprietario",
	Label:       "Criar Novo Proprietario",
	Description: "Cria um Novo Proprietario",
	Method:      "POST",
	Callers:     []string{"$org3MSP", "$orgMSP"}, // Only org3 can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "id",
			Label:       "ID",
			Description: "Identificador do Proprietario",
			DataType:    "string",
			Required:    true,
		},
		{
			// Mandatory property
			Required: true,
			Tag:      "nome",
			Label:    "Nome",
			DataType: "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		id, _ := req["id"].(string)
		nome, _ := req["nome"].(string)

		proprietarioMap := make(map[string]interface{})
		proprietarioMap["@assetType"] = "proprietario"
		proprietarioMap["id"] = id
		proprietarioMap["nome"] = nome

		proprietarioAsset, err := assets.NewAsset(proprietarioMap)
		if err != nil {
			return nil, errors.WrapError(err, "Failed to create a new asset")
		}

		// Salva o Proprietário
		_, err = proprietarioAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "Error saving asset on blockchain")
		}

		proprietarioJSON, nerr := json.Marshal(proprietarioAsset)
		if nerr != nil {
			return nil, errors.WrapError(nil, "failed to encode asset to JSON format")
		}

		return proprietarioJSON, nil
	},
}
