package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// POST Method
var TransferirToken = tx.Transaction{
	Tag:         "transferToken",
	Label:       "Transferir Token",
	Description: "Transferir Token",
	Method:      "PUT",
	Callers:     []string{`$org\dMSP`, "orgMSP"}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "token",
			Label:       "Token Origem",
			Description: "Token Origem",
			DataType:    "->token",
			Required:    true,
		},
		{
			Tag:         "idTokenDestino",
			Label:       "ID Token Destino",
			Description: "ID Token Destino",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "proprietario",
			Label:       "Proprietario do Token Destino",
			Description: "Proprietario do Token Destino",
			DataType:    "->proprietario",
			Required:    true,
		},
		{
			Tag:         "idNovoTokenOrigem",
			Label:       "ID Novo Token Origem",
			Description: "ID Novo Token Origem",
			DataType:    "string",
			Required:    true,
		},
		{
			Tag:         "quantidade",
			Label:       "Quantidade",
			Description: "Quantidade de Token",
			DataType:    "number",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		tokenKey, ok := req["token"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter token must be an asset")
		}

		idTokenDestino, ok := req["idTokenDestino"]
		if !ok {
			return nil, errors.WrapError(nil, "Parameter idTokenDestino must be an string")
		}

		idNovoTokenOrigem, ok := req["idNovoTokenOrigem"]
		if !ok {
			return nil, errors.WrapError(nil, "Parameter idNovoTokenOrigem must be an string")
		}

		proprietario, ok := req["proprietario"]
		if !ok {
			return nil, errors.WrapError(nil, "Parameter proprietario must be an asset")
		}

		quantidade, ok := req["quantidade"]
		if !ok {
			return nil, errors.WrapError(nil, "Parameter idToken must be an number")
		}

		// Retorna o Token
		tokenOrigem, err := tokenKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}
		tokenOrigemMap := (map[string]interface{})(*tokenOrigem)

		burned := tokenOrigemMap["burned"].(bool)
		//Valida se o Token já foi queimado
		if burned {
			return nil, errors.WrapError(nil, "Token já queimado")
		}

		saldo := tokenOrigemMap["quantidade"].(float64) - quantidade.(float64)
		//Valida Saldo da Transação
		if saldo <= 0 {
			return nil, errors.WrapError(nil, "Saldo insuficiente")
		}

		//Token Destino
		CreateToken(
			stub,
			idTokenDestino.(string),
			proprietario,
			quantidade.(float64))

		//Marca Token origem como Burned
		tokenOrigemMap["burned"] = true

		tokenOrigemMap, err = tokenOrigem.Update(stub, tokenOrigemMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}

		//Novo Token Origem
		novoTokenOrigem, err := CreateToken(
			stub,
			idNovoTokenOrigem.(string),
			tokenOrigemMap["proprietario"],
			saldo)

		tokenJSON, nerr := json.Marshal(novoTokenOrigem)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return tokenJSON, nil
	},
}

func CreateToken(stub *sw.StubWrapper, id string, proprietario interface{}, quantidade float64) (assets.Asset, errors.ICCError) {
	tokenDestinoMap := make(map[string]interface{})
	tokenDestinoMap["@assetType"] = "token"
	tokenDestinoMap["id"] = id
	tokenDestinoMap["proprietario"] = proprietario
	tokenDestinoMap["quantidade"] = quantidade

	newTokenDestinoAsset, err := assets.NewAsset(tokenDestinoMap)
	if err != nil {
		return nil, errors.WrapError(err, "Failed to create a new asset")
	}

	// Salva o novo Token
	_, err = newTokenDestinoAsset.PutNew(stub)
	if err != nil {
		return nil, errors.WrapError(err, "Error saving asset on blockchain")
	}

	return newTokenDestinoAsset, nil
}
