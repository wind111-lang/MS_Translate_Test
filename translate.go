package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(fmt.Sprintf("../%s.env", os.Getenv("key")))
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

	//fmt.Print("//IMPORTANT PLEASE READ//\n\nCheck your subscriptionKey and location.\nSubscriptionKey: "+AzureData.subscriptionKey, "\nlocation: "+AzureData.location, "\n\n")

	u, _ := url.Parse(uri)
	v, _ := url.Parse(uri)
	q := u.Query()
	r := v.Query()
	//q.Add("from", "ja")
	q.Add("to", "en")
	r.Add("from", "en")
	r.Add("to", "ja")
	u.RawQuery = q.Encode()
	v.RawQuery = r.Encode()
	//再翻訳するために２つ作成

	//Create an anonymous struct for your request body and encode it to JSON

	var insert string

	fmt.Println("再翻訳テスト：日本語で何か入力してください")
	fmt.Scanf("%s", &insert)

	body := []struct {
		Text string
	}{
		{Text: insert},
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

	// Decode the JSON response
	var result interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	// Format and print the response to terminal
	prettyJSON, _ := json.MarshalIndent(result, "", " ")
	fmt.Printf("%s\n", prettyJSON)

	// ////////////////////////////////////ここまで1回目翻訳///////////////////////////////////

	var a interface{} //prettyJSONの中身を取り出すためにインタフェースを定義

	err = json.Unmarshal(prettyJSON, &a)
	if err != nil {
		log.Fatal(err)
	}
	temp := a.([]interface{})
	temp2 := temp[0]
	temp3 := temp2.(map[string]interface{})["translations"]
	temp4 := temp3.([]interface{})
	temp5 := temp4[0]
	text_str := temp5.(map[string]interface{})["text"].(string)
	/*ちなみに、全部取り出そうと思うと、[]interface {}{map[string]interface {}
	{"translations":[]interface {}{map[string]interface {}
	{"text":"Hello.(example)", "to":"en"}}}} となる。なんだっそら！*/

	//Unmarshalしたものを"text"まで取り出す

	Datum := []struct { //再翻訳させるためにJSON用の構造体を再度用意
		Text string
	}{
		{Text: text_str}, //英語に翻訳された文
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

	var retranslate interface{}
	if err = json.NewDecoder(res2.Body).Decode(&retranslate); err != nil {
		log.Fatal(err)
	} //デコード処理

	// Format and print the response to terminal
	retranslate_res, _ := json.MarshalIndent(retranslate, "", " ")
	fmt.Printf("%s\n", retranslate_res) //再翻訳されたもの
	//JSONで出力
}
