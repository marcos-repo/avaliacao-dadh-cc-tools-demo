package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// GET method
var ContabilizarTokens = tx.Transaction{
	Tag:         "contabilizarTokens",
	Label:       "Contabilizar Tokens",
	Description: "Retorna Contabilização de Tokens",
	Method:      "GET",
	Callers:     []string{"$org1MSP", "$org2MSP", "$orgMSP"}, // Only org1 and org2 can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "proprietario",
			Label:       "Proprietario",
			Description: "Proprietario do TOken",
			DataType:    "->proprietario",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		proprietario, _ := req["proprietario"]

		query := map[string]interface{}{
			"selector": map[string]interface{}{
				"@assetType":   "token",
				"proprietario": proprietario,
			},
		}

		var err error
		response, err := assets.Search(stub, query, "", true)

		qtde := 0.0

		result := response.Result

		for i := 0; i < len(result); i++ {
			//Contabiliza apenas a quantidade dos Tokens que naõ foram queimados
			if result[i]["burned"].(bool) {
				continue
			}

			qtde += result[i]["quantidade"].(float64)
		}

		balance := make(map[string]interface{})
		balance["Saldo"] = qtde

		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "Erro contabilizando Tokens", 500)
		}

		responseJSON, err := json.Marshal(balance)
		if err != nil {
			return nil, errors.WrapErrorWithStatus(err, "error marshaling response", 500)
		}

		return responseJSON, nil
	},
}
