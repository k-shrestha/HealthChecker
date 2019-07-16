package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"urlHealthChecker/UrlHealthChecker/HealthChecker/model"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

//Init initialises database
func Init() {
	var err error
	db, err = gorm.Open("mysql", "root:@tcp(127.0.0.1:3307)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("Connection not established")
		fmt.Println(err)
	}

	db.AutoMigrate(&model.URLConfigModel{})
	db.AutoMigrate(&model.URLStatusModel{})
}

//AddURL adds urls data to database
func AddURL(c *gin.Context) {

	var urlCheck []model.URLConfigModel
	c.BindJSON(&urlCheck)

	for item := range urlCheck {
		db.Create(&urlCheck[item])
		c.JSON(200, urlCheck[item])
	}
}

//CheckStatus hit the url and check its status
func CheckStatus(url string, crawlTime int, frequency int, failureThreshold int) {

	timeout := time.Duration(crawlTime) * time.Second
	httpClient := http.Client{
		Timeout: timeout,
	}

	for i := 1; i <= failureThreshold; i++ {

		resp, err := httpClient.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode != 200 {
			res := fmt.Sprintf(`%s is down on %d trial`, url, i)
			color.Red(res)
			urlStat := model.URLStatusModel{URL: url, StatusCode: resp.StatusCode, FailureCount: i}
			db.Save(&urlStat)
		} else {
			res := fmt.Sprintf(`%s is up`, url)
			color.HiGreen(res)
			urlStat := model.URLStatusModel{URL: url, StatusCode: resp.StatusCode, FailureCount: i}
			db.Save(&urlStat)
			defer resp.Body.Close()
			break
		}
		time.Sleep(time.Duration(frequency) * time.Second)
	}

}

//FetchStatus calls checkstatus for each urls
func FetchStatus(c *gin.Context) {
	var urlConfigs []model.URLConfigModel
	db.Find(&urlConfigs)
	if len(urlConfigs) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No url found!"})
		return
	}

	for _, item := range urlConfigs {
		go CheckStatus(item.URL, item.CrawlTime, item.Frequency, item.FailureThreshold)
	}

}

//UpdateURLData updates the url data
func UpdateURLData(c *gin.Context) {
	var urlConfig model.URLConfigModel
	urlConfigID := c.Param("id")
	db.First(&urlConfig, urlConfigID)
	if urlConfig.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No url found!"})
		return
	}
	db.Model(&urlConfig).Update("url", c.PostForm("url"))
	crawlTime, _ := strconv.Atoi(c.PostForm("crawlTime"))
	frequency, _ := strconv.Atoi(c.PostForm("frequency"))
	failureThreshold, _ := strconv.Atoi(c.PostForm("failureThreshold"))
	db.Model(&urlConfig).Update("crawlTime", crawlTime)
	db.Model(&urlConfig).Update("frequency", frequency)
	db.Model(&urlConfig).Update("failureThreshold", failureThreshold)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "url updated successfully!"})
}
