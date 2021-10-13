package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"translate/structs"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(fmt.Sprintf("%s.env", os.Getenv("key")))
	if err != nil {
		fmt.Println("Error loading environment")
		log.Fatal(err)
	}
	//envファイル読み込み処理

	Key := os.Getenv("subscriptionKey")
	location := os.Getenv("location")
	endpoint := os.Getenv("endpoint")
	uri := endpoint + os.Getenv("uri")
	//envファイルで読み込んだものを代入

	//IMPORTANT PLEASE READ Check your subscriptionKey and location.

	u, _ := url.Parse(uri)
	v, _ := url.Parse(uri)
	q := u.Query()
	r := v.Query()
	q.Add("from", "ja")
	q.Add("to", "en")
	r.Add("from", "en")
	r.Add("to", "ja")
	u.RawQuery = q.Encode()
	v.RawQuery = r.Encode()
	//再翻訳するために２つ作成

	//Create an anonymous struct for your request body and encode it to JSON
	var msg string
	fmt.Println("入力して")
	fmt.Scanf("%s", &msg)
	text := msg

	body := []struct {
		Text string `json:"text"`
	}{
		{Text: text},
	}

	b, _ := json.Marshal(body)

	// Build the HTTP POST request
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	// Add required headers to the request
	req.Header.Add("Ocp-Apim-Subscription-Key", Key)
	req.Header.Add("Ocp-Apim-Subscription-Region", location)
	req.Header.Add("Content-Type", "application/json")

	// Call the Translator API
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var arr []structs.TranslationRes

	translations, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(translations, &arr)
	if err != nil {
		log.Fatal(err)
	}

	text_str := arr[0].Translation[0].Text
	//fmt.Println(text_str)
	///////////////////////////////////////ここまで1回目翻訳///////////////////////////////////

	Datum := []struct { //再翻訳するためのjson構造体を用意する
		Text string `json:"text"`
	}{
		{Text: text_str},
	}
	c, _ := json.Marshal(Datum)

	// Build the HTTP POST request
	req2, err := http.NewRequest("POST", v.String(), bytes.NewBuffer(c))
	if err != nil {
		log.Fatal(err)
	}

	// Add required headers to the request
	req2.Header.Add("Ocp-Apim-Subscription-Key", Key)
	req2.Header.Add("Ocp-Apim-Subscription-Region", location)
	req2.Header.Add("Content-Type", "application/json")
	//翻訳と再翻訳で２回行われるため、Azure上の使用文字数は２倍になる

	// Call the Translator API
	res2, err := http.DefaultClient.Do(req2)
	if err != nil {
		log.Fatal(err)
	}

	translations2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(translations2, &arr)
	if err != nil {
		log.Fatal(err)
	}
	text_str2 := arr[0].Translation[0].Text

	fmt.Println("英語翻訳後: ", text_str)
	fmt.Println("再翻訳後: ", text_str2)
}
