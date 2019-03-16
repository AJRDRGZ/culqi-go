package culqi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiVersion = "v2.0"
	baseURL    = "https://api.culqi.com/v2"
)

// Errors API
var (
	errInvalidRequest = errors.New("La petición tiene una sintaxis inválida")
	errAuthentication = errors.New("La petición no pudo ser procesada debido a problemas con las llaves")
	errCard           = errors.New("No se pudo realizar el cargo a una tarjeta")
	errLimitAPI       = errors.New("Estás haciendo muchas peticiones rápidamente al API o superaste tu límite designado")
	errResource       = errors.New("El recurso no puede ser encontrado, es inválido o tiene un estado diferente al permitido")
	errAPI            = errors.New("Error interno del servidor de Culqi")
	errUnexpected     = errors.New("Error inesperado, el código de respuesta no se encuentra controlado")
	ErrParameter      = errors.New("Algún parámetro de la petición es inválido")
	ErrResponse       = errors.New("")
)

// WrapperResponse respuesta generica para respuestas GetAll
type WrapperResponse struct {
	Paging struct {
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Cursors  struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
	} `json:"paging"`
}

type errorResponse struct {
	Object          string `json:"object"`
	Type            string `json:"type"`
	Code            string `json:"code"`
	MerchantMessage string `json:"merchant_message"`
	UserMessage     string `json:"user_message"`
}

func do(method, endpoint string, params url.Values, body io.Reader) ([]byte, error) {
	if len(params) != 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+keyInstance.SecretKey)

	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	obj, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	rErr := &errorResponse{}
	err = json.Unmarshal(obj, rErr)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 400:
		err = errInvalidRequest
	case 401:
		err = errAuthentication
	case 422:
		err = ErrParameter
	case 402:
		err = errCard
	case 429:
		err = errLimitAPI
	case 404:
		err = errResource
	case 500, 503, 504:
		err = errAPI
	}

	if err != nil {
		ErrResponse = fmt.Errorf("%v | %s | %s", err, rErr.MerchantMessage, rErr.UserMessage)
		return nil, ErrResponse
	}

	if res.StatusCode >= 200 && res.StatusCode <= 206 {
		return obj, nil
	}

	return nil, fmt.Errorf("%v:(%d) %s", errUnexpected, res.StatusCode, string(obj))
}
