package photolab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	API_ENDPOINT        = "http://api-soft.photolab.me"
	API_UPLOAD_ENDPOINT = "http://upload-soft.photolab.me/upload.php"
	API_ENDPOINT_PROXY  = "http://api-proxy-soft.photolab.me"
)

type Steps struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Id        int64    `json:"id"`
	ImageUrls []string `json:"image_urls"`
}

type Image struct {
	Url      string  `xml:",innerxml"`
	Order    int     `xml:"order,attr"`
	Rectf    string  `xml:"rectf,attr,omitempty"`
	Rotation float32 `xml:"rotation,attr,omitempty"`
	Flip     int     `xml:"flip,attr,omitempty"`
}

type ImageRequest struct {
	Url    string
	Rotate int
	Flip   int
	Crop   string
}

type PhotolabClient struct {
}

func (p *PhotolabClient) PhotolabStepsAdvanced(comboId int64) (steps Steps, err error) {
	formData := url.Values{}
	formData.Set("combo_id", fmt.Sprint(comboId))

	endpoint := fmt.Sprintf("%s/photolab_steps_advanced.php", API_ENDPOINT)

	data, err := p.query(endpoint, formData)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data, &steps); err != nil {
		return
	}
	return
}

func (p *PhotolabClient) TemplateProcess(templateName string, ims []ImageRequest) (s string, err error) {
	formData := url.Values{}
	formData.Set("template_name", templateName)
	for i, v := range ims {
		t := fmt.Sprintf("%d", i+1)
		formData.Set("template_name", templateName)
		formData.Set(fmt.Sprintf("image_url[%s]", t), v.Url)
		formData.Set(fmt.Sprintf("rotate[%s]", t), fmt.Sprint(v.Rotate))
		formData.Set(fmt.Sprintf("crop[%s]", t), v.Crop)
		formData.Set(fmt.Sprintf("flip[%s]", t), fmt.Sprint(v.Flip))
	}
	endpoint := fmt.Sprintf("%s/template_process.php", API_ENDPOINT)
	data, err := p.query(endpoint, formData)
	s = string(data)
	return
}

func (p *PhotolabClient) PhotolabProcess(templateName string, ims []ImageRequest) (s string, err error) {
	formData := url.Values{}
	formData.Set("template_name", templateName)
	for i, v := range ims {
		t := fmt.Sprintf("%d", i+1)
		formData.Set("template_name", templateName)
		formData.Set(fmt.Sprintf("image_url[%s]", t), v.Url)
		formData.Set(fmt.Sprintf("rotate[%s]", t), fmt.Sprint(v.Rotate))
		formData.Set(fmt.Sprintf("crop[%s]", t), v.Crop)
		formData.Set(fmt.Sprintf("flip[%s]", t), fmt.Sprint(v.Flip))
	}
	endpoint := fmt.Sprintf("%s/template_process.php", API_ENDPOINT)
	data, err := p.query(endpoint, formData)
	s = string(data)
	return
}

func (p *PhotolabClient) ImageUpload(image string) (s string, err error) {
	file, err := os.Open(image)
	if err != nil {
		return
	}
	defer file.Close()
	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multiPartWriter.CreateFormFile("file1", "image.jpg")
	if err != nil {
		return
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return
	}
	fieldWriter, err := multiPartWriter.CreateFormField("no_resize")
	if err != nil {
		return
	}
	_, err = fieldWriter.Write([]byte("1"))
	if err != nil {
		return
	}
	multiPartWriter.Close()

	req, err := http.NewRequest("POST", API_UPLOAD_ENDPOINT, &requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("StatusCode %d from: %s", resp.StatusCode, req.URL.String())
		msg += " Response:" + string(data)
		err = fmt.Errorf(msg)
		return
	}
	s = string(data)
	return
}

func (p *PhotolabClient) TemplateUpload(resources string) (s string, err error) {
	file, err := os.Open(resources)
	if err != nil {
		return
	}
	defer file.Close()
	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)
	fileWriter, err := multiPartWriter.CreateFormFile("resources", "resources.zip")
	if err != nil {
		return
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return
	}
	multiPartWriter.Close()
	endpoint := fmt.Sprintf("%s/template_upload.php", API_ENDPOINT_PROXY)
	req, err := http.NewRequest("POST", endpoint, &requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("StatusCode %d from: %s", resp.StatusCode, req.URL.String())
		msg += " Response:" + string(data)
		err = fmt.Errorf(msg)
		return
	}
	s = string(data)
	return
}

func (p *PhotolabClient) DowloadFile(endpoint, dst string) (err error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("StatusCode %d from: %s", resp.StatusCode, endpoint)
		msg += " Response:" + string(data)
		err = fmt.Errorf(msg)
		return
	}
	err = ioutil.WriteFile(dst, data, 0666)
	return
}

func (p *PhotolabClient) query(endpoint string, formValues url.Values) (data []byte, err error) {
	httpReq, err := http.NewRequest("POST", endpoint, strings.NewReader(formValues.Encode()))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("StatusCode %d from: %s", resp.StatusCode, httpReq.URL.String())
		msg += " Response:" + string(data)
		err = fmt.Errorf(msg)
		return
	}
	return
}
