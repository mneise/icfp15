all :
	go build -o play_icfp2015

test : all
	./play_icfp2015 -f p0.json

