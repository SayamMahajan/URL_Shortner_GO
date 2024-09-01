package main
import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type UrlDetails struct {
	Id          string    `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortUrl    string    `json:"short_url"`
	CreatedDate time.Time `json:"date"`
}

 //for temperary storage of data
var Info = make(map[string]UrlDetails)

//function to generate short url of 10 element
func generateShortUrl(OriginalUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalUrl)) //convert url string into byte slice
	data := hasher.Sum(nil)
	fmt.Println("Hash sum", data)
	hexhash := hex.EncodeToString(data) //give a string url
	fmt.Println("ShortURL is", hexhash[:10])
	return hexhash[:10]
}

//function to make the short url store in info/db
func createUrl(originalUrl string) string {
	shortUrl := generateShortUrl(originalUrl)
	id := shortUrl
	Info[id] = UrlDetails{
		Id:          id,
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
		CreatedDate: time.Now(),
	}
	return shortUrl
}

//function to retrieve the url using url id
func getUrl(id string) (UrlDetails, error) {
	url, ok := Info[id]
	if !ok {
		return UrlDetails{}, errors.New("URL not found")
	}
	return url, nil
}

//function to handle root page request response
func rootPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to URL SHORTNER...")
}

//function to give the short url to user 
func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Url string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortUrl := createUrl(data.Url)
	
	respone := struct {
		ShortUrl string `json:"short_url"`
	}{ShortUrl : shortUrl}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respone)
}

//function to handle redirect request
func redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/redirect/"):] 
    url, err := getUrl(id)
    if err != nil {
        http.Error(w, "URL not found", http.StatusNotFound)
    }

    http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}


func main() {
	fmt.Println("Welcome to URL Shortner...")

	// Handling get request
	http.HandleFunc("/", rootPageHandler)

	//handling short url request
	http.HandleFunc("/shortner", shortUrlHandler)

	//handling redirect request
	http.HandleFunc("/redirect/", redirectUrlHandler)

	//server setup
	fmt.Println("Server is live on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error while starting server", err)
	}
}
