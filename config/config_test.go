package config

import "testing"

const (
	sample1 = `
    # This is your bao config file, define 
    # a port and any env variables 
    PORT 8000
    THIS=that
    THAT="this"`
	sample2 = `
        # This is your bao config file, define 
    # a port and any env variables 
    # PORT 8001
    #THIS=that
        PORT 8000
        THIS = that
        THAT="th#is"`
)

func TestParse(t *testing.T) {
	response, err := Parse(sample1)
	if err != nil {
		panic(err)
	}
	if response.Port != 8000 {
		t.Errorf("Wrong port number")
	}
	if len(response.EnvVars) != 2 {
		t.Errorf("Not enough env vars")
	}
	//more detailed verification here

	response, err = Parse(sample2)
	if err != nil {
		panic(err)
	}
	if response.Port != 8000 {
		t.Errorf("Wrong port number")
	}
	if len(response.EnvVars) != 2 {
		t.Errorf("Not enough env vars")
	}
}
