package payment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/Shabrinashsf/go-xendit-payment-webhook/dto"
)

func SendXenditInvoice(invoice dto.XenditInvoice) (map[string]interface{}, error) {
	//Crete HTTP request to Xendit API
	apiKey := os.Getenv("XENDIT_API_KEY")
	urlXendit := "https://api.xendit.co/v2/invoices"
	authToken := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))

	// Make header
	headers := map[string]string{
		"Authorization": "Basic " + authToken,
		"Content-Type":  "application/json",
	}

	data, err := json.Marshal(invoice)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlXendit, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Get Response
	bodyResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	// Parsing response to json
	var bodyResponseJSON map[string]interface{}
	if err := json.Unmarshal(bodyResponse, &bodyResponseJSON); err != nil {
		return nil, err
	}

	return bodyResponseJSON, nil
}
