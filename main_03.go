package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arivalRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarberDoneChan  chan bool
	ClientChan      chan string
	Open            bool
}

func main() {
	rand.Seed(time.Now().UnixNano())

	color.Yellow("The sleeping barber problem")
	color.Yellow("---------------------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientChan:      clientChan,
		BarberDoneChan:  doneChan,
		Open:            true,
	}

	color.Green("The shop is open for the day")

	shop.addBarber("Frank")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)

		shopClosing <- true
		shop.closeShopForDay()
		closed <- true

	}()

	i := 0

	go func() {

		for {
			randomMillSecond := rand.Int() % (2 * arivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMillSecond)):
				shop.addClient(fmt.Sprintf("Client  #%d", i))
			}
		}

	}()

	<-closed

}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		color.Yellow("%s goes to waiting room to check client", barber)

		for {
			if len(shop.ClientChan) == 0 {
				color.Yellow("There is nothing to do , so %s takes a nap", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientChan
			if shopOpen {
				if isSleeping {
					color.Yellow("%s wakes %s up.", client, barber)
					isSleeping = false
				}
				shop.cutHair(barber, client)
			} else {
				shop.SendBarberHome(barber)
				return
			}

		}

	}()
}

func (shop *BarberShop) cutHair(barber string, client string) {
	color.Green("%s cutting hair of %s", barber, client)
	time.Sleep(shop.HairCutDuration)
	color.Green("%s is finished cutting %s's hair.", barber, client)

}

func (shop *BarberShop) SendBarberHome(barber string) {
	color.Cyan("%s is going home", barber)
	shop.BarberDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	color.Cyan("closing shop for the day")
	close(shop.ClientChan)

	shop.Open = false

	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarberDoneChan
	}

	close(shop.BarberDoneChan)

	color.Green("The barber shop now closed for the day")
}

func (shop *BarberShop) addClient(client string) {
	color.Green("%s arrives", client)
	if shop.Open {
		select {
		case shop.ClientChan <- client:
			color.Yellow("%s takes a seat in the waiting room", client)
		default:
			color.Red("The waiting room is full %s leaves ", client)
		}

	} else {
		color.Red("The shop is already closed,so %s leaves", client)
	}
}
