package utility

import (
	"testing"
)

type emailValidation struct {
    arg1 string
	expected bool
}

var emailValidations = []emailValidation{
    {"prefix@domainname.com", true},
    {"prefix123@domainname.com", true},
    {"invaliddomainname.com", false},
    {"invalid2@domainname", false},
    
}

func TestEmailValidation(t *testing.T){
    for _, test := range emailValidations{
        if output := IsEmailValid(test.arg1); output != test.expected {
            t.Errorf("Output %v not equal to expected %v", output, test.expected)
        }
    }
}

func TestEnvFallback(t *testing.T){
    got := GetEnv("ENV_TOKEN", "fallbackToken")
    want := "fallbackToken"

    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}

