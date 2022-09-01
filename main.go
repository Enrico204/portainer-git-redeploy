package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Stack struct {
	EndpointId  int64                 `json:"EndpointID"`
	Environment []RedeploySettingsEnv `json:"Env"`
	GitConfig   struct {
		ReferenceName  string `json:"ReferenceName"`
		Authentication struct {
			Username string `json:"Username"`
			Password string `json:"Password"`
		} `json:"Authentication"`
	}
}

type RedeploySettingsEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RedeploySettings struct {
	RepositoryUsername       string                `json:"repositoryUsername"`
	RepositoryReferenceName  string                `json:"repositoryReferenceName"`
	RepositoryPassword       string                `json:"repositoryPassword"`
	RepositoryAuthentication bool                  `json:"repositoryAuthentication"`
	Environment              []RedeploySettingsEnv `json:"env"`
}

func main() {
	var err error
	var tlsCertData []byte

	var portainerURL = flag.String("url", "", "Portainer URL (env: PORTAINER_URL)")
	var apiKey = flag.String("access-token", "", "Portainer Access Token (UNSAFE, use environment variable PORTAINER_ACCESS_TOKEN)")
	var stackId = flag.Int64("stack-id", 0, "Portainer Stack ID (env: PORTAINER_STACK_ID)")
	var timeout = flag.Duration("timeout", 120*time.Second, "HTTP timeout (env: PORTAINER_HTTP_TIMEOUT)")
	var tlsCert = flag.String("tls-cert", "", "Use an additional TLS certificate / CA as trusted (env: PORTAINER_SSL_CERT_FILE)")

	flag.Parse()

	if os.Getenv("PORTAINER_URL") != "" {
		*portainerURL = os.Getenv("PORTAINER_URL")
	}
	if os.Getenv("PORTAINER_ACCESS_TOKEN") != "" {
		*apiKey = os.Getenv("PORTAINER_ACCESS_TOKEN")
	}
	if os.Getenv("PORTAINER_STACK_ID") != "" {
		*stackId, err = strconv.ParseInt(os.Getenv("PORTAINER_STACK_ID"), 10, 64)
		if err != nil {
			panic(err)
		}
	}
	if *portainerURL == "" || *apiKey == "" || *stackId == 0 {
		flag.PrintDefaults()
		return
	}

	if os.Getenv("PORTAINER_HTTP_TIMEOUT") != "" {
		*timeout, err = time.ParseDuration(os.Getenv("PORTAINER_HTTP_TIMEOUT"))
		if err != nil {
			panic(err)
		}
	}
	if os.Getenv("PORTAINER_SSL_CERT_FILE") != "" {
		*tlsCert = os.Getenv("PORTAINER_SSL_CERT_FILE")
	}
	if os.Getenv("PORTAINER_SSL_CERT") != "" {
		tlsCertData = []byte(os.Getenv("PORTAINER_SSL_CERT"))
	} else if *tlsCert != "" {
		tlsCertData, err = ioutil.ReadFile(*tlsCert)
		if err != nil {
			panic(err)
		}
	}

	httpc, err := newClient(*timeout, tlsCertData)
	if err != nil {
		panic(err)
	}

	stack, err := getStack(httpc, *portainerURL, *apiKey, *stackId)
	if err != nil {
		panic(err)
	}

	err = redeploy(httpc, *portainerURL, *apiKey, *stackId, stack.EndpointId, RedeploySettings{
		RepositoryUsername:       stack.GitConfig.Authentication.Username,
		RepositoryReferenceName:  stack.GitConfig.ReferenceName,
		RepositoryPassword:       stack.GitConfig.Authentication.Password,
		RepositoryAuthentication: true,
		Environment:              stack.Environment,
	})
	if err != nil {
		panic(err)
	}
}

func newClient(timeout time.Duration, tlsCert []byte) (*http.Client, error) {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	} else if rootCAs == nil {
		// Empty root CA?
		rootCAs = x509.NewCertPool()
	}

	if len(tlsCert) > 0 {
		if ok := rootCAs.AppendCertsFromPEM(tlsCert); !ok {
			return nil, errors.New("no certificate found to append")
		}
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
		Timeout: timeout,
	}, nil
}

func getStack(httpc *http.Client, portainerURL string, apiKey string, stackId int64) (Stack, error) {
	var ret Stack

	// First, we need to retrieve the current environment and endpointId
	url := fmt.Sprintf("%s/api/stacks/%d", portainerURL, stackId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ret, err
	}

	req.Header.Set("X-API-Key", apiKey)

	resp, err := httpc.Do(req)
	if err != nil {
		return ret, err
	}

	if resp.StatusCode != http.StatusOK {
		return ret, errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&ret)
	_ = resp.Body.Close()
	return ret, err
}

func redeploy(httpc *http.Client, portainerURL string, apiKey string, stackId int64, endpointId int64, settings RedeploySettings) error {
	buf, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	// First, we need to retrieve the current environment and endpointId
	url := fmt.Sprintf("%s/api/stacks/%d/git/redeploy?endpointId=%d", portainerURL, stackId, endpointId)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
