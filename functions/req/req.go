package req

import (
	"errors"
	"io/ioutil"
	"net/http"
)

//Post请求
func PostForm(reqUrl string, postData map[string][]string) ([]byte, error) {
	client := &http.Client{}
	//post请求
	resp, err := client.PostForm(reqUrl, postData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	return nil, errors.New("请求有误")
}

//Get请求
func Get(reqUrl string) ([]byte, error) {
	client := &http.Client{}
	//post请求
	resp, err := client.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	}

	return nil, errors.New("请求有误")
}
