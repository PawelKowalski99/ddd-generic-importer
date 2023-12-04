package counter

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	// extVerifier = emailverifier.NewVerifier().DisableCatchAllCheck().DisableSMTPCheck().DisableAutoUpdateDisposable().DisableDomainSuggest().DisableGravatarCheck()
	extVerifier = emailverifier.NewVerifier().EnableAutoUpdateDisposable().EnableCatchAllCheck().EnableDomainSuggest().EnableSMTPCheck()
)

type Result struct {
	Email  string
	Domain string
	IP     string
}

type Verifier interface {
	Verify(in string) (*Result, error)
}

type externalVerifier struct {
}

func (*externalVerifier) Verify(in string) (*Result, error) {
	res, err := extVerifier.Verify(in)
	if err != nil {
		return nil, err
	}

	return &Result{Email: res.Email, Domain: res.Syntax.Domain}, nil
}

type internalVerifier struct{}

func (*internalVerifier) Verify(in string) (*Result, error) {

	// Small sleeper to check how amount of workers makes difference
	time.Sleep(time.Nanosecond * time.Duration(rand.Intn(20)))

	res := strings.Split(in, "@")
	if len(res) != 2 {
		return nil, errors.New("mail is not correct")

	}

	return &Result{Email: in, Domain: res[1]}, nil
}
