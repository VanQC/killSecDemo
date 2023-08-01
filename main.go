package main

import "killDemo/router"

func main() {

	r := router.Router()

	r.Run(":8088")
}
