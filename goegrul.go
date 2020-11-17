package goegrul

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type tokenJSON struct {
	T               string `json:"t"`
	CaptchaRequired bool   `json:"captchaRequired"`
}

type firmJSON struct {
	Rows []struct {
		Address          string `json:"a"`
		ShortName        string `json:"c"`
		ExpirationDate   string `json:"e"`
		Director         string `json:"g"`
		Cnt              string `json:"cnt"`
		INN              string `json:"i"`
		K                string `json:"k"`
		FullName         string `json:"n"`
		OGRN             string `json:"o"`
		KPP              string `json:"p"`
		RegistrationDate string `json:"r"`
		Token            string `json:"t"`
		Pg               string `json:"pg"`
		Tot              string `json:"tot"`
	} `json:"rows"`
}

func getToken(inn string) (string, error) {
	formData := url.Values{}
	formData.Add("vyp3CaptchaToken", "")
	formData.Add("page", "")
	formData.Add("query", inn)
	formData.Add("region", "")
	formData.Add("PreventChromeAutocomplete", "")

	resp, err := http.PostForm("https://egrul.nalog.ru/", formData)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := tokenJSON{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	return res.T, nil
}

//Firm - result for return
type Firm struct {
	Name     string
	FullName string
	Address  string
	INN      string
	KPP      string
	Director string
	Yurik    bool //юридическое лицо или физическое
	Expired  bool //не действует
}

func getFirm(t string, inn string) (Firm, error) {
	f := Firm{}
	resp, err := http.Get("https://egrul.nalog.ru/search-result/" + t)
	if err != nil {
		return f, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return f, err
	}

	res := firmJSON{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return f, err
	}
	for i := 0; i < len(res.Rows); i++ {
		if res.Rows[0].INN == inn {
			if res.Rows[i].K == "fl" {
				f.Yurik = false
				f.FullName = trimValues(res.Rows[i].FullName)
				f.INN = trimValues(res.Rows[i].INN)
				f.Expired = (trimValues(res.Rows[i].FullName) != "")
			} else {
				f.Yurik = true
				f.Address = trimValues(res.Rows[i].Address)
				f.Name = trimValues(res.Rows[i].ShortName)
				f.FullName = trimValues(res.Rows[i].FullName)
				f.INN = trimValues(res.Rows[i].INN)
				f.KPP = trimValues(res.Rows[i].KPP)
				f.Director = trimValues(res.Rows[i].Director)
				f.Expired = (trimValues(res.Rows[i].FullName) != "")
			}
			return f, nil
		}
	}
	return f, errors.New("organization not found")
}

//в адресах может быть кучка пробелов
func trimValues(s string) string {
	str := strings.ReplaceAll(s, "  ", " ")
	str = strings.ReplaceAll(str, "  ", " ")
	return strings.TrimSpace(str)
}

//GetDataFromEGRUL - качает данные с сервера налоговой
func GetDataFromEGRUL(inn string) (Firm, error) {
	var wg sync.WaitGroup
	var token string
	var err error
	//запросим в асинхронной горутине токен
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Выполеняем запрос
		token, err = getToken(inn)
	}()
	wg.Wait()
	//обработаем ошибку
	f := Firm{}
	if err != nil {
		return f, err
	}

	//запросим в асинхронной горутине фирму
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Выполеняем запрос
		f, err = getFirm(token, inn)
	}()
	wg.Wait()
	//обработаем ошибку
	if err != nil {
		return f, err
	}

	return f, nil
}
