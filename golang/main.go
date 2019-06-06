package main

import (
	"fmt"
	"github.com/PhotoLabDevelopment/hackathon-sdk/golang/photolab"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	api := new(photolab.PhotolabClient)

	resoursesFilename := "resources.zip"
	if _, err := os.Stat(resoursesFilename); os.IsNotExist(err) {
		if err := api.DowloadFile("http://soft.photolab.me/samples/resources.zip", resoursesFilename); err != nil {
			log.Fatalln(resoursesFilename, "api.DowloadFile err:", err)
		}
	}

	contentFilename := "girl.jpg"
	if _, err := os.Stat(contentFilename); os.IsNotExist(err) {
		if err := api.DowloadFile("http://soft.photolab.me/samples/girl.jpg", contentFilename); err != nil {
			log.Fatalln(contentFilename, "api.DowloadFile err:", err)
		}
	}

	contentUrl, err := api.ImageUpload(contentFilename)
	if err != nil {
		log.Fatalln("api.ImageUpload err:", err)
	}
	fmt.Printf("content url: %s\n", contentUrl)

	templateName, err := api.TemplateUpload(resoursesFilename)
	if err != nil {
		log.Fatalln("api.TemplateUpload err: %v", err)
	}
	fmt.Printf("template name: %s\n", templateName)

	a := []photolab.ImageRequest{}
	a = append(a, photolab.ImageRequest{
		Url:    contentUrl,
		Flip:   0,
		Rotate: 0,
		Crop:   "0,0,1,1",
	})
	resultUrl, err := api.TemplateProcess(templateName, a)
	if err != nil {
		log.Fatalln(templateName, "api.TemplateProcess", err)
	}
	fmt.Printf("for template_name: %s, result_url: %s\n", templateName, resultUrl)

	for _, comboId := range []int64{5635874, 3124589} {
		var resultUrl string
		originalContentUrl := contentUrl
		fmt.Println("===")

		fmt.Println("start process combo_id:", comboId)
		steps, err := api.PhotolabStepsAdvanced(comboId)
		if err != nil {
			log.Fatalln("api.PhotolabStepsAdvanced", comboId, "err:", err)
			return
		}

		for _, v := range steps.Steps {
			templateName := fmt.Sprint(v.Id)
			a := []photolab.ImageRequest{}
			for _, imageUrl := range v.ImageUrls {
				if imageUrl == "" {
					imageUrl = originalContentUrl
				}
				a = append(a, photolab.ImageRequest{
					Url:    imageUrl,
					Flip:   0,
					Rotate: 0,
					Crop:   "0,0,1,1",
				})
			}
			if len(a) == 0 {
				a = append(a, photolab.ImageRequest{
					Url:    originalContentUrl,
					Flip:   0,
					Rotate: 0,
					Crop:   "0,0,1,1",
				})
			}
			resultUrl, err = api.TemplateProcess(templateName, a)
			if err != nil {
				log.Fatalln("api.TemplateProcess", templateName, "err:", err)
				return
			}
			originalContentUrl = resultUrl
			fmt.Printf("--- for template_name: %s, result_url: %s\n", templateName, resultUrl)
		}
		fmt.Println("result_url:", resultUrl)
	}
}
