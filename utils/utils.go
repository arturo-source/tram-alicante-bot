package utils

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	SERVICE_URL = "https://www.tramalicante.es/wp-admin/admin-ajax.php"
)

func DoRequest(v url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", SERVICE_URL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.8,en-US;q=0.5,en;q=0.3")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://www.tramalicante.es/ca/consulta-horaris-i-planificador/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Origin", "https://www.tramalicante.es")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "PHPSESSID=b52tfnp4tb6adciebh8jp4ju21; wp-wpml_current_language=ca")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
