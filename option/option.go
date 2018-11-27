package option

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ClientOption is an option for herschel client
type ClientOption interface {
	GetClient() (*http.Client, error)
}

// WithConfigFileAndTokenFile returns a ClientOption that loads config and token from given file paths
func WithConfigFileAndTokenFile(configFile string, tokenFile string) ClientOption {
	return withConfigFileAndTokenFile{configFile: configFile, tokenFile: tokenFile}
}

type withConfigFileAndTokenFile struct {
	configFile string
	tokenFile  string
}

func (w withConfigFileAndTokenFile) GetClient() (*http.Client, error) {
	config, err := getConfigFromFile(w.configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config from file from %s", w.configFile)
	}

	token, err := getTokenFromFile(w.tokenFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load token from file from %s", w.tokenFile)
	}

	return config.Client(context.Background(), token), nil
}

func getConfigFromFile(filePath string) (*oauth2.Config, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return google.ConfigFromJSON(bytes, "https://www.googleapis.com/auth/spreadsheets")
}

func getTokenFromFile(filePath string) (*oauth2.Token, error) {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	if err = json.NewDecoder(f).Decode(token); err != nil {
		return nil, err
	}
	return token, nil
}

// WithConfigAndToken returns a ClientOption that uses given config and token
func WithConfigAndToken(config *oauth2.Config, token *oauth2.Token) ClientOption {
	return withConfigAndToken{config, token}
}

type withConfigAndToken struct {
	config *oauth2.Config
	token  *oauth2.Token
}

func (w withConfigAndToken) GetClient() (*http.Client, error) {
	return w.config.Client(context.Background(), w.token), nil
}

// WithConfigReaderAndTokenReader returns a ClientOption that loads config and token from given readers
func WithConfigReaderAndTokenReader(configReader io.Reader, tokenReader io.Reader) ClientOption {
	return withConfigReaderAndTokenReader{configReader: configReader, tokenReader: tokenReader}
}

type withConfigReaderAndTokenReader struct {
	configReader io.Reader
	tokenReader  io.Reader
}

func (w withConfigReaderAndTokenReader) GetClient() (*http.Client, error) {
	config, err := getConfigFromReader(w.configReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	token, err := getTokenFromReader(w.tokenReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load token")
	}

	return config.Client(context.Background(), token), nil
}

// WithServiceAccountCredentials returns a ClientOption that loads credentials from given file path
func WithServiceAccountCredentials(file string) ClientOption {
	return withServiceAccountCredentials{credentialsFile: file}
}

type withServiceAccountCredentials struct {
	credentialsFile string
}

func (w withServiceAccountCredentials) GetClient() (*http.Client, error) {
	data, err := ioutil.ReadFile(w.credentialsFile)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	return conf.Client(context.Background()), nil
}

func getConfigFromReader(r io.Reader) (*oauth2.Config, error) {
	buffer := new(bytes.Buffer)
	io.Copy(buffer, r)

	return google.ConfigFromJSON(buffer.Bytes(), "https://www.googleapis.com/auth/spreadsheets")
}

func getTokenFromReader(r io.Reader) (*oauth2.Token, error) {
	token := &oauth2.Token{}
	if err := json.NewDecoder(r).Decode(token); err != nil {
		return nil, err
	}
	return token, nil
}
