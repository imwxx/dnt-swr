package handler

import (
	"fmt"
    "log"
    "strings"
    "io/ioutil"
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
    lib "dnt-swr/lib"
    harbor "dnt-swr/harbor"
)

var health = http.StatusOK

func HandleCheckHealth(c *gin.Context) {
	c.JSON(health, nil)
}

func verifKey(c *gin.Context) bool {
    ckey := c.Request.Header.Get("verification-key")
    skey := lib.AppConf.VERFKEY
    if ckey != skey {
        fmt.Println(c.Request.Header)
        return false
    }
    health = http.StatusOK
    return true
}
func HandleSetHealth(c *gin.Context) {
	statusMap := map[string]string{"200": "200", "503": "503"}
	ip := c.ClientIP()
	if ip == "127.0.0.1" {
		setStatus := c.Query("status")
		if _, ok := statusMap[setStatus]; ok {
			switch setStatus {
			case "200":
				health = http.StatusOK
				c.JSON(health, gin.H{"msg": "service will online"})
			case "503":
				health = http.StatusServiceUnavailable
				c.JSON(health, gin.H{"msg": "service will offline"})
			}
		}
	} else {
		c.JSON(health, gin.H{"msg": "invaild remote_addr"})
	}
}

func HandleRepoAction(c *gin.Context) {
    if ! verifKey(c) {
        c.JSON(http.StatusUnauthorized, gin.H{"status": false, "msg": "verifykey was error"})
        return
    }
    harborRepo := c.PostForm("harborRepo")
    swrRepo := c.PostForm("swrRepo")
    action := c.PostForm("action")
    if harborRepo == "" || swrRepo == "" || action == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "parameter should not be null"})
        return
    }
    var dbRsp lib.RESPONSEACTION
    switch action {
        case "on":
            dbRsp = lib.Insert(harborRepo, swrRepo, action, lib.MDBC)
            if ! dbRsp.Status {
                health = http.StatusBadRequest
            }
            c.JSON(health, gin.H{"status": dbRsp.Status, "msg": dbRsp.Msg})
        case "off":
            data := map[string]string{"harborRepo":harborRepo, "swrRepo":swrRepo}
            dbRsp = lib.OffRepoSync(data, lib.MDBC)
            if ! dbRsp.Status {
                health = http.StatusBadRequest 
            }
            c.JSON(health, gin.H{"status": dbRsp.Status, "msg": dbRsp.Msg})
        case "check":
            data := map[string]string{"harborRepo":harborRepo, "swrRepo":swrRepo}
            dbRsp = lib.CheckRepoSync(data, lib.MDBC)
            if ! dbRsp.Status {
                health = http.StatusBadRequest
            }
            c.JSON(health, gin.H{"status": dbRsp.Status, "msg": dbRsp.Msg})
        case "del":
            data := map[string]string{"harborRepo":harborRepo, "swrRepo":swrRepo}
            dbRsp = lib.DelRepoSync(data, lib.MDBC)
            c.JSON(health, gin.H{"status": dbRsp.Status, "msg": dbRsp.Msg})
            if ! dbRsp.Status {
                health = http.StatusInternalServerError
            }
        default:
            c.JSON(http.StatusBadRequest, gin.H{"status": false, "msg": "action parameter just support ['on', 'off', 'check']"})
    }
}

func HandleDnt(c *gin.Context) {
    data, err := ioutil.ReadAll(c.Request.Body)
    if err != nil {
        log.Fatal(err)
    }
    resp, err := harbor.UnMarshal(data)
    if err != nil {
        fmt.Println(err)
    } else {
        events := resp.EVENTS[0]
        action := events.ACTION
        if action == "push" && events.TARGET.TAG != "" {
            tag := events.TARGET.TAG
            repository := events.TARGET.REPOSITORY
            harborRepo := strings.Split(repository, "/")[0]
            harborProject := strings.Split(repository, "/")[1]
            host := events.REQUEST.HOST
            addr := events.REQUEST.ADDR
            harborUri := host + "/" + repository + ":" + tag
            crepo := lib.AllowSync(harborRepo, lib.MDBC)
            if crepo.Status == true {
                swrRepo := crepo.Res[harborRepo]
                swrUri := lib.AppConf.HWDOMAIN + "/" + swrRepo + "/" + harborProject + ":" + tag
                cmd := lib.RedisClient.HSet(harborUri, swrUri, addr)
                fmt.Println(cmd)
            } else {
                fmt.Printf("%s repo was not allow\n", harborRepo)
            }
        }
    }
    c.JSON(health, gin.H{"msg": "hello"})
}

