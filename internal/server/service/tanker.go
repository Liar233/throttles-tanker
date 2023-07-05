package service

import (
	"time"

	"github.com/Liar233/throttles-tank/pkg/tanker"
)

type TankerConfig struct {
	NetworkTankCapacity  uint32 `mapstructure:"network.capacity"`
	NetworkTankPerUpdate uint32 `mapstructure:"network.per_update"`

	LoginTankCapacity  uint32 `mapstructure:"login.capacity"`
	LoginTankPerUpdate uint32 `mapstructure:"login.per_update"`

	PasswordTankCapacity  uint32 `mapstructure:"password.capacity"`
	PasswordTankPerUpdate uint32 `mapstructure:"password.per_update"`
}

type Tanker struct {
	addressTank   tanker.Tanker
	loginTank     tanker.Tanker
	passwordTank  tanker.Tanker
	timeStampChan <-chan time.Time
	brake         chan bool
}

func (t *Tanker) Check(address, login, password string) bool {

	if !t.addressTank.Check(address) {

		return false
	}

	if !t.passwordTank.Check(password) {

		return false
	}

	return t.loginTank.Check(login)
}

func (t *Tanker) Run() {

	t.brake = make(chan bool)

	go func() {
		for {
			select {
			case _, ok := <-t.timeStampChan:

				if ok {
					t.addressTank.Flush()
					t.passwordTank.Flush()
					t.loginTank.Flush()
				}

				break
			case <-t.brake:
				return
			}
		}
	}()
}

func (t *Tanker) Reset(address, login string) {

	t.addressTank.Remove(address)
	t.loginTank.Remove(login)
}

func (t *Tanker) Close() {

	defer close(t.brake)
}

func NewTanker(
	config TankerConfig,
	ticker <-chan time.Time,
) *Tanker {

	networkBuilder := func() tanker.Bucket {

		return tanker.NewTokenBucket(
			config.NetworkTankCapacity,
			config.NetworkTankPerUpdate,
		)
	}

	loginBuilder := func() tanker.Bucket {

		return tanker.NewTokenBucket(
			config.LoginTankCapacity,
			config.LoginTankPerUpdate,
		)
	}

	passwordBuilder := func() tanker.Bucket {

		return tanker.NewTokenBucket(
			config.PasswordTankCapacity,
			config.PasswordTankPerUpdate,
		)
	}

	return &Tanker{
		addressTank:   tanker.NewTank(networkBuilder),
		loginTank:     tanker.NewTank(loginBuilder),
		passwordTank:  tanker.NewTank(passwordBuilder),
		timeStampChan: ticker,
		brake:         make(chan bool),
	}
}
