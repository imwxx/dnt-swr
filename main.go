package main

import (
    "dnt-swr/route"
    "dnt-swr/ginit"
)

func main() {
    listenInfo := ginit.InitServer()
    route.RouteAll(listenInfo)

}
