package harbor

import (
    "encoding/json"
    //"fmt"
    //"reflect"
)

type SSOURCE struct {
    ADDR string `json:"addr"`
    INSTANCEID string `json:"instanceID"`
}

type PULLTARGET struct {
    MEDIATYPE string `json:"mediaType"`
    SIZE int `json:"size"`
    DIGEST string `json:"digest"`
    LENGTH int `json:"length"`
    REPOSITORY string `json:"repository"`
    URL string `json:"url"`
}

type PUSHTARGET struct {
    MEDIATYPE string `json:"mediaType"`
    SIZE int `json:"size"`
    DIGEST string `json:"digest"`
    LENGTH int `json:"length"`
    REPOSITORY string `json:"repository"`
    URL string `json:"url"`
    TAG string `json:"tag"`
}

type PULLEVENT struct {
    ID string `json:"id"`
    TIMESTAMP string `json:"timestamp"`
    ACTION string `json:"action"`
    TARGET PULLTARGET `json:"target"`
    REQUEST struct {
            ID string `json:"id"`
            ADDR string `json:"addr"`
            HOST string `json:"host"`
            METHOD string `json:"method"`
            USERAGENT string `json:"useragent"`
        } `json:"request"`
    ACTOR map[string]string `json:"name"`
    SOURCE SSOURCE `json:"source"`
}

type PUSHEVENT struct {
    ID string `json:"id"`
    TIMESTAMP string `json:"timestamp"`
    ACTION string `json:"action"`
    TARGET PUSHTARGET `json:"target"`
    REQUEST struct {
            ID string `json:"id"`
            ADDR string `json:"addr"`
            HOST string `json:"host"`
            METHOD string `json:"method"`
            USERAGENT string `json:"useragent"`
        } `json:"request"`
    ACTOR map[string]string `json:"name"`
    SOURCE SSOURCE `json:"source"`
}

type PUSHEVENTS struct {
    EVENTS []PUSHEVENT `json:"events"`
}

type PULLEVENTS struct {
    EVENTS []PULLEVENT `json:"events"`
}

func UnMarshal(jsonData []byte) (*PUSHEVENTS, error){
    pushevents := &PUSHEVENTS{}
    err := json.Unmarshal(jsonData, pushevents)
    if err != nil {
        return pushevents, err
    } else {
        //fmt.Println(pushevents.EVENTS[0].ACTION)
        return pushevents, nil
    }
}
