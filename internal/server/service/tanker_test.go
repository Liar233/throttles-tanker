package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/Liar233/throttles-tank/pkg/tanker"
)

func TestTanker_Check(t *testing.T) {

	var entries [4][3]string
	entries[0] = [3]string{"192.168.0.1", "test_user", "password"}
	entries[1] = [3]string{"192.168.10.10", "test_user", "password"}
	entries[2] = [3]string{"192.168.0.1", "test_user1", "password"}
	entries[3] = [3]string{"192.168.0.1", "test_user", "password1"}

	addrTank := &tankStub{Truth: "192.168.0.1"}
	passTank := &tankStub{Truth: "password"}
	loginTank := &tankStub{Truth: "test_user"}

	var ticker chan time.Time

	tanker := NewTestedTanker(addrTank, loginTank, passTank, ticker)

	for i, entry := range entries {

		res := tanker.Check(entry[0], entry[1], entry[2])

		if i == 0 && !res || i > 0 && res {
			t.Log(i, res)
			t.Error("Tanker Check() method failed")
		}
	}
}

func TestTanker_Run(t *testing.T) {

	ticker := make(chan time.Time)

	defer close(ticker)

	signal := make(chan int)

	defer close(signal)

	addrTank := &tankStub{numb: 1, signal: signal}
	loginTank := &tankStub{numb: 2, signal: signal}
	passTank := &tankStub{numb: 3, signal: signal}

	tanker := NewTestedTanker(addrTank, loginTank, passTank, ticker)

	tanker.Run()

	ticker <- time.Now()

	counter := 0

	for i := 0; i < 3; i++ {
		counter += <-signal
	}

	if counter != 6 {

		t.Error("Tanker Flush() method failed")
	}
}

func TestTanker_Reset(t *testing.T) {

	var ticker chan time.Time

	addrTank := &tankStub{Truth: "address"}
	loginTank := &tankStub{Truth: "login"}
	passTank := &tankStub{Truth: "password"}

	tanker := NewTestedTanker(addrTank, loginTank, passTank, ticker)

	tanker.Reset("address", "login")

	valid := addrTank.Truth != "address!" || loginTank.Truth != "login!" || passTank.Truth != "password"

	if valid {

		t.Error("Method Reset() invalid")
	}
}

func NewTestedTanker(
	addrTank, loginTank, passTank tanker.Tanker,
	ticker <-chan time.Time,
) *Tanker {

	return &Tanker{
		addressTank:   addrTank,
		loginTank:     loginTank,
		passwordTank:  passTank,
		timeStampChan: ticker,
		brake:         make(chan bool),
	}
}

type tankStub struct {
	Truth  string
	signal chan int
	numb   int
}

func (t *tankStub) Check(id string) bool {

	return id == t.Truth
}

func (t *tankStub) Flush() {

	t.signal <- t.numb
}

func (t *tankStub) Remove(id string) {

	if id == t.Truth {

		t.Truth = fmt.Sprintf("%s!", id)
	}
}
