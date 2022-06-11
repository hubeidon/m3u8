package initial


func Run(level string){
	InitLogger(level)

	go CompositeVideo()
}