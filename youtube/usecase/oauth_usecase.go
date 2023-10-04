package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"tiktok_api/app/logger"
	"tiktok_api/youtube/repository/redis"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

var log = logger.NewLogrusLogger()
var ctx = context.Background()
var config *oauth2.Config

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	clientSecretPath := filepath.Join(path, "/app/config/client_secret.json")
	b, err := ioutil.ReadFile(clientSecretPath) //- current directory json file
	if err != nil {
		handleError(err, "Unable to read Youtube client secret file", "error")
	}

	//- If modifying these scopes, delete your previously saved credentials
	//- at ~/.credentials/youtube-go-quickstart.json
	config, err = google.ConfigFromJSON(b,
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeUploadScope,
		youtube.YoutubeScope,
	)
	if err != nil {
		handleError(err, "Google cannot load config from json file", "error")
	}
	log.Info("Load Youtube Client Secret file successfully")
}

// - FIXME: receive client_id, client_secret, project_id,...from client input

func Exec() *oauth2.Config {
	var config *oauth2.Config

	//- FIXME: Get file from config folder
	b, err := ioutil.ReadFile("client_secret.json") //- current directory json file
	if err != nil {
		handleError(err, "Unable to read client secret file", "error")
	}

	//- If modifying these scopes, delete your previously saved credentials
	//- at ~/.credentials/youtube-go-quickstart.json
	config, err = google.ConfigFromJSON(b,
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeUploadScope,
		youtube.YoutubeScope,
	)
	if err != nil {
		handleError(err, "Unable to parse client secret file to config", "error")
	}
	return config
}

// - getClient uses a Context and Config to retrieve a Token
// - then generate a Client. It returns the generated Client.
func GetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	log.Println("here at GetClient() function")
	cacheFile, err := tokenCacheFile()
	if err != nil {
		handleError(err, "Unable to get path to cached credential file", "error")
	}

	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		//- create new user from config
		clientKey, _ := redis.CreateNewYoutubeClient(config.ClientID, config.ClientSecret)
		tok = getTokenFromWeb(config, clientKey)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// - getTokenFromWeb uses Config to request a Token.
// - It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config, clientKey string) *oauth2.Token {
	state := clientKey
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		handleError(err, "Unable to read authorization code", "error")
	}

	tok, err := config.Exchange(ctx, code)
	if err != nil {
		handleError(err, "Unable to retrieve token from web", "error")
	}
	return tok
}

func YoutubeOAuthCodeExchange(clientKey string, code string) string {
	successURL := "https://www.meritize.com/wp-content/uploads/2018/06/keys-to-success.jpg"
	failURL := "https://www.freecodecamp.org/news/content/images/2023/02/image-126.png"
	tokens, err := config.Exchange(ctx, code)
	if err != nil {
		handleError(err, "Unable to retrieve token from web", "error")
	}

	//- check if clientKey is exist
	//- if exists, update value with clientKey
	//- if not exist, CSRF => panic
	isClientKeyExist := redis.IsExist(clientKey)

	if !isClientKeyExist {
		// panic("CSRF violation, your process is stopped here!")
		return failURL
	}

	//- if exist
	//- update tokens for clientKey
	youtubeOAuth := redis.GetClientByClientKey(clientKey)
	youtubeOAuth.AccessToken = tokens.AccessToken
	youtubeOAuth.RefreshToken = tokens.RefreshToken
	youtubeOAuth.Expiry = tokens.Expiry

	redis.UpdateYoutubeByClientKey(clientKey, youtubeOAuth)
	return successURL
}

func BuildClientFromTokens(clientKey string) *http.Client {
	//- FIXME: check if clientKey is exist

	//- get client tokens base on clientKey
	youtubeOAuth := redis.GetClientByClientKey(clientKey)
	tokens := &oauth2.Token{
		AccessToken:  youtubeOAuth.AccessToken,
		RefreshToken: youtubeOAuth.RefreshToken,
		Expiry:       youtubeOAuth.Expiry,
		TokenType:    "Bearer",
	}
	return config.Client(ctx, tokens)
}

func BuildServiceFromToken(clientKey string) *youtube.Service {
	//- FIXME: check if clientKey is exist

	//- get client tokens base on clientKey
	youtubeOAuth := redis.GetClientByClientKey(clientKey)
	tokens := &oauth2.Token{
		AccessToken:  youtubeOAuth.AccessToken,
		RefreshToken: youtubeOAuth.RefreshToken,
		Expiry:       youtubeOAuth.Expiry,
		TokenType:    "Bearer",
	}
	client := config.Client(ctx, tokens)
	service, err := youtube.New(client)
	if err != nil {
		handleError(err, "Unable to create Youtube service", "error")
	}
	return service
}

func GetAuthURL() (string, string) {
	clientKey, _ := redis.CreateNewYoutubeClient(config.ClientID, config.ClientSecret)
	state := clientKey
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, clientKey
}

// - tokenCacheFile generates credential file path/filename.
// - It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("client_credential.json")), err
}

// - tokenFromFile retrieves a Token from a given file path.
// - It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// - saveToken uses a file path to create a file and store the
// - token in it.
func saveToken(file string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		handleError(err, "Unable to cache oauth token", "error")
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
