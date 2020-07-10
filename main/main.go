package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gutil/bhx"
	"io/ioutil"
	"mux"
	"net/http"
)

type DataType struct {
	JsonData string
	CryptResult []byte
	EncryptData []byte
}

type Rates struct {
	XMLName xml.Name `xml:"rates"`
	Items []Item `xml:"item"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	From string `xml:"from"`
	To string `xml:"to"`
	In string `xml:"in"`
	Out string `xml:"out"`
	Amount string `xml:"amount"`
	Minamount string `xml:"minamount"`
	Maxamount string `xml:"maxamount"`
	Param string `xml:"param"`
	City string `xml:"city"`
}

func main () {
	resp, err := http.Get("https://test.cryptohonest.ru/request-exportxml.xml")
	if err != nil {
		fmt.Println(err)
		return
	}

	byteAll, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(byteAll))

	var allRates Rates
	var Data DataType

	xml.Unmarshal([]byte(byteAll), &allRates)
	ByteJsonData, _ := json.Marshal(allRates)
	Data.JsonData = string(ByteJsonData)
	//fmt.Println(string(jsonData.JsonData))

	r := mux.NewRouter()
	r.HandleFunc("/", msgHandler)
	r.HandleFunc("/courses", Data.courcesHandler)
	r.HandleFunc("/postform", Data.postformHandler)
	http.Handle("/", r)

	fmt.Println("Sever is listening...")
	http.ListenAndServe("localhost:8181", r)
}

func msgHandler (resp http.ResponseWriter, r *http.Request) {
	http.ServeFile(resp, r, "static/courses.html")
	/*m := "Go to localhost:8181/courses"
	fmt.Fprint(resp, m)*/
}

func (j DataType) courcesHandler (resp http.ResponseWriter, r *http.Request) {
	fmt.Fprint(resp, j.JsonData)
	//http.ServeFile(resp, r, "static/courses.html")
	/*data := DataType{JsonData: j.JsonData}
	tmpl, _ := template.ParseFiles("templates/courses.html")
	err := tmpl.Execute(resp, data)
	if err != nil{
		fmt.Println(err)
	}*/
}

func (j DataType) postformHandler (resp http.ResponseWriter, r *http.Request) {
	text := []byte(r.FormValue("Textstring"))
	//cyphertext := []byte(r.FormValue("Cyphertextstring"))
	keyString := []byte(r.FormValue("Encryptkey"))
	Hash := bhx.GetSha256Hash(keyString)
	nnonce := bhx.GetKeyNonce(bhx.BoxSharedKey(Hash))
	j.CryptResult, _ = bhx.Encrypt(text, bhx.BoxSharedKey(Hash), nnonce)
	/*if err != nil {
		fmt.Println(err)
	}*/
	j.EncryptData, _ = bhx.Decrypt(j.CryptResult, bhx.BoxSharedKey(Hash), nnonce)
	fmt.Fprintln(resp,"Результат шифрования: ", j.CryptResult)
	fmt.Fprintln(resp, "Результат дешифрования: ", string(j.EncryptData))
	/*data := JsonDataType{
		JsonData:    "fuck fuck fuck",
		CryptResult: "suka suka suka",
	}
	tmpl, _ := template.ParseFiles("templates/coursesEncrypt.html")
	tmpl.Execute(resp, data)*/
}
