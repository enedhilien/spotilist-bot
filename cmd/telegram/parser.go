package telegram

import "regexp"

var pattern = "https:\\/\\/open\\.spotify\\.com\\/track\\/([a-zA-Z0-9]*)"

func ParseTrack(message string) (bool, string){
	r := regexp.MustCompile(pattern)
	if match := r.MatchString(message); match{
		return true, r.FindStringSubmatch(message)[1]
	}else{
		return false, ""
	}
}
