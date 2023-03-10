/**
 * Created by goland.
 * @file   docker.go
 * @author 李锦 <Lijin@cavemanstudio.net>
 * @date   2023/1/10 11:11
 * @desc   docker.go
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getToken() string {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	// curl -s 'https://hub.docker.com/v2/users/login'   -H "Content-Type:application/json"   -X POST   -d '{"username": "kcmonkey","password": "hahaha8888"}'
	type Payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	data := Payload{
		Username: "kcmonkey",
		Password: "hahaha8888",
		// fill struct
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Panicln(err)
	}
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", "https://hub.docker.com/v2/users/login", body)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()
	var token = struct {
		Token string `json:"token"`
	}{}
	res, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(res, &token)
	return token.Token
}
func getLastVersion(project string, token string) string {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	// curl 'https://registry.hub.docker.com/v2/namespaces/kcmonkey/repositories/console/tags' \
	// -H "Content-Type:application/json" \
	// -H "Authorization: Bearer eyJ4NWMiOlsiTUlJQytUQ0NBcCtnQXdJQkFnSUJBREFLQmdncWhrak9QUVFEQWpCR01VUXdRZ1lEVlFRREV6dFNUVWxHT2xGTVJqUTZRMGRRTXpwUk1rVmFPbEZJUkVJNlZFZEZWVHBWU0ZWTU9rWk1WalE2UjBkV1dqcEJOVlJIT2xSTE5GTTZVVXhJU1RBZUZ3MHlNakF4TVRBeU1qSXhORGxhRncweU16QXhNalV5TWpJeE5EbGFNRVl4UkRCQ0JnTlZCQU1UTzFCWlExSTZTVkJhUWpwSVFsRlhPamRNUlZrNlFrRldRanBIU0RkYU9sVklWek02VWt0TFVqcEROMDR6T2xsTk5GSTZTVmhOU0RwS1ZGQkNNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQWxOSThhbTJQREZ2czZ3TnlZdndHZGVmV1VyV0FnTGdzdjMwNjJ2cnBuVTZDdFlKRGQ2ditjRFRtRTdRZU5ZRGhTMmlNcFN3Y2ZEQlwvZEVyUTB2eEhkTjJwMlwvODZmZ1wvU3lpSDJ4ZjBhVU45Q1dXbndCT2kyXC9LeEt2K2lsU2VDTUdhdHBGWDdKZjFxYjg3RDk1TE5UMG85T05OZjFPdGJ2NjlyT21cL1RIVFh3clVcL1d3U2ZVMlpLVGxMOElUV0ZEVzdPWSt4V3VCdFlKbXlYanFabGlkQW1DU3U3R2NGNDFQeXpvSkxRUzJyQnVycDl3NHFoMTBZNW1DSHFnbGhCNUZPWjlLNE9qVGlYVFByRVJOVnJwK1BVSEdyQWFUT0UwcEM4MlBwbllWYTczQ1JHbDBHQ3RcL1FyQnBWNGl2azB3MXR4S2RXU2JIM2dGa2pnaDU3S05wOFFJREFRQUJvNEd5TUlHdk1BNEdBMVVkRHdFQlwvd1FFQXdJSGdEQVBCZ05WSFNVRUNEQUdCZ1JWSFNVQU1FUUdBMVVkRGdROUJEdFFXVU5TT2tsUVdrSTZTRUpSVnpvM1RFVlpPa0pCVmtJNlIwZzNXanBWU0Zjek9sSkxTMUk2UXpkT016cFpUVFJTT2tsWVRVZzZTbFJRUWpCR0JnTlZIU01FUHpBOWdEdFNUVWxHT2xGTVJqUTZRMGRRTXpwUk1rVmFPbEZJUkVJNlZFZEZWVHBWU0ZWTU9rWk1WalE2UjBkV1dqcEJOVlJIT2xSTE5GTTZVVXhJU1RBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlFQTdIY1VyVm1namo1cE01MXhZVHd2eGE1VnRqd2hub0dRZjFxTU52UGVHeVlDSUFwYm5cL1hZXC9LUXlZYVZGdGMxa2xvSWZnN3hcL3hlbkZhbkp4XC9BdnFERnRYIl0sInR5cCI6IkpXVCIsImFsZyI6IlJTMjU2In0.eyJpc3MiOiJodHRwczpcL1wvYXBpLmRvY2tlci5jb21cLyIsImF1ZCI6WyJodHRwczpcL1wvaHViLmRvY2tlci5jb20iXSwianRpIjoiMGIzZGIyMmMtZDUyNC00YjY3LTlhNzYtMGYzNWFjODc5MjA5IiwidXNlcl9pZCI6ImRkMTQyM2QxMzE5YjQ0ZGM4MjAwNmRjYTc2NTcyNjcxIiwidXNlcm5hbWUiOiJrY21vbmtleSIsImlhdCI6MTY3MzI2MDU5MiwiaHR0cHM6XC9cL2h1Yi5kb2NrZXIuY29tIjp7ImVtYWlsIjoia2Ntb25rZXkxOTkyQGdtYWlsLmNvbSIsInJvbGVzIjpbXSwic291cmNlIjoiZG9ja2VyX3B3ZHxkZDE0MjNkMS0zMTliLTQ0ZGMtODIwMC02ZGNhNzY1NzI2NzEiLCJ1dWlkIjoiZGQxNDIzZDEtMzE5Yi00NGRjLTgyMDAtNmRjYTc2NTcyNjcxIiwidXNlcm5hbWUiOiJrY21vbmtleSIsInNlc3Npb25faWQiOiIwYjNkYjIyYy1kNTI0LTRiNjctOWE3Ni0wZjM1YWM4NzkyMDkifSwic291cmNlIjp7ImlkIjoiZGQxNDIzZDEtMzE5Yi00NGRjLTgyMDAtNmRjYTc2NTcyNjcxIiwidHlwZSI6InB3ZCJ9LCJzdWIiOiJkZDE0MjNkMTMxOWI0NGRjODIwMDZkY2E3NjU3MjY3MSIsInNlc3Npb25faWQiOiIwYjNkYjIyYy1kNTI0LTRiNjctOWE3Ni0wZjM1YWM4NzkyMDkiLCJleHAiOjE2NzU4NTI1OTJ9.a8qUreEZW4BPrn-UwYc8dPP74BFVhIcLc4TrcsdMjYFSbCBcBHAra8obW4n0Nhf4J7XILvag7upu0tb4FyahU7g_pIVzp1bVgWqzkJVnDRmc7wjjCYKNG4f_yG-ACFoxOySIiCPQMe3fxCD_98X7GZpExUJuCySpfjdKaZ9rS-Ld0awKE7G383mZsmVKLpDT5NbODC_JeJVKFRDAWccWVSOGuq39nQzkCVupJk1OBbXZZfo2EoFFEiYRgaqynQL47QsDGCGmnKZeTx8OxLL30qzxgzYOs6yqx1FiPojkDcRWEptROX9Sbs7hew6-Fe8jtiCFJLGbYZRjgb5sGKgZ_A" \
	// -H "Content-Type:application/json" | json_reformat | grep name
	req, err := http.NewRequest("GET", "https://registry.hub.docker.com/v2/namespaces/kcmonkey/repositories/"+project+"/tags", nil)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()
	res, _ := ioutil.ReadAll(resp.Body)
	version := struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}{}
	_ = json.Unmarshal(res, &version)
	return version.Results[0].Name
}
func newVersion(currentVersion string) string {
	currentVersion = strings.ReplaceAll(currentVersion, "v", "")
	currentVersion = strings.ReplaceAll(currentVersion, ".", "")
	v, _ := strconv.Atoi(currentVersion)
	newVersion := v + 1
	versionStr := "0" + strconv.Itoa(newVersion)
	if newVersion < 10 {
		versionStr = "00" + strconv.Itoa(newVersion)
	} else if newVersion > 100 {
		versionStr = strconv.Itoa(newVersion)
	}
	var versionArr []string
	for i := 0; i < len(versionStr); i++ {
		versionArr = append(versionArr, string(versionStr[i]))
	}
	version := strings.Join(versionArr, ".")
	return "v" + version
}

func main() {
	if len(os.Args) < 2 {
		log.Panicln("params error")
	}
	token := getToken()
	version := getLastVersion(os.Args[1], token)
	version = newVersion(version)
	fmt.Print(version)
}
