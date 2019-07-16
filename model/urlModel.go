package model

import "github.com/jinzhu/gorm"

type (
	//URLConfigModel is model
	URLConfigModel struct {
		gorm.Model
		URL              string `gorm:"primary_key" json:"url"`
		CrawlTime        int    `json:"crawlTime"`
		Frequency        int    `json:"frequency"`
		FailureThreshold int    `json:"failureThreshold"`
	}

	//TransformedurlConfig for reading and storing
	TransformedurlConfig struct {
		ID               int    `json:"id"`
		URL              string `json:"url"`
		CrawlTIme        int    `json:"crawlTime"`
		Frequency        int    `json:"frequency"`
		FailureThreshold int    `json:"failureThreshold"`
	}
)

type (
	//URLStatusModel is model
	URLStatusModel struct {
		gorm.Model
		URL          string `json:"url"`
		StatusCode   int    `json:"stCode"`
		FailureCount int    `json:"failureCount"`
	}

	//TransformedurlStatus for reading and storing
	TransformedurlStatus struct {
		ID           int    `json:"id"`
		URL          string `json:"url"`
		StatusCode   int    `json:"stCode"`
		FailureCount int    `json:"failureCount"`
	}
)
