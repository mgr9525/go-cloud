package gocloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type CloudExecs struct {
	Serv string
	Host string
}

func (c *CloudExecs) execHttp(path string, data *url.Values) (int, []byte, error) {
	host := c.Host
	if len(host) <= 0 && len(c.Serv) > 0 && Consul != nil {
		services, err := Consul.Agent().Services()
		if err != nil {
			return 0, nil, err
		}
		service := services[c.Serv]
		if service == nil {
			return 0, nil, errors.New("no service")
		}
		host = fmt.Sprintf("http://%s:%d", service.Address, service.Port)
	}
	if data == nil {
		data = &url.Values{}
	}
	res, err := http.PostForm(host+path, *data)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, bts, nil
}
func (c *CloudExecs) execHttpJSON(path string, body interface{}) (int, []byte, error) {
	host := c.Host
	if len(host) <= 0 && len(c.Serv) > 0 && Consul != nil {
		services, err := Consul.Agent().Services()
		if err != nil {
			return 0, nil, err
		}
		service := services[c.Serv]
		if service == nil {
			return 0, nil, errors.New("no service")
		}
		host = fmt.Sprintf("http://%s:%d", service.Address, service.Port)
	}
	if body == nil {
		body = &map[string]interface{}{}
	}

	bts, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequest("P0ST", host+path, bytes.NewBuffer(bts))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	byts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, byts, nil
}

func (c *CloudExecs) Exec(path string, pars *url.Values) (int, string, error) {
	code, bts, err := c.execHttp(path, pars)
	if err != nil {
		return code, "", err
	}
	return code, string(bts), nil
}
func (c *CloudExecs) ExecObj(path string, pars *url.Values, ret interface{}) error {
	code, bts, err := c.execHttp(path, pars)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New("code is:" + strconv.Itoa(code))
	}
	err = json.Unmarshal(bts, ret)
	if err != nil {
		return err
	}

	return nil
}

func (c *CloudExecs) ExecJSON(path string, pars interface{}) (int, string, error) {
	code, bts, err := c.execHttpJSON(path, pars)
	if err != nil {
		return code, "", err
	}
	return code, string(bts), nil
}
func (c *CloudExecs) ExecObjJSON(path string, pars interface{}, ret interface{}) error {
	code, bts, err := c.execHttpJSON(path, pars)
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New("code is:" + strconv.Itoa(code))
	}
	err = json.Unmarshal(bts, ret)
	if err != nil {
		return err
	}

	return nil
}
