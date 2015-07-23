/*---------------------------------------------------------------------------- *\
   Dining Pphilosophers

   This rendintion of the problem was intended for a presentation I put together
   for some of my coworkers. It's a well known issue that has many many many
   implementations out on the net however, I did this with a slant twoards
   showing off some of Go's cool concurency features.

   Are there bugs... probably but, keep in mind this was written to show off
   language features, not to be the best possible solution to the problem.

   For more details on the problem:
   https://en.wikipedia.org/wiki/Dining_philosophers_problem

   No, you may not use this for your homework assignment

   Just in case there are any questions... this entire program is licensed under
   GPLv3, which you aught to have received a copy of with this source file. Just
   in case, you can find it here:

   http://www.gnu.org/licenses/gpl-3.0.html

   Enjoy

   Author:
   Karim Sharif <karim@sharifclan.com>
\* ----------------------------------------------------------------------------*/
package main

// some imports
import (
	"fmt"
	"math/rand"
	"time"
)

// some globals
var (
	philosophers []*Philosopher
)

// Some constants
const NUMBER_OF_MEALS = 130
const LEFT = 0
const RIGHT = 1

/* Our chopstick type */
type Chopstick struct {
	isClean bool
}

// Member method cleans the chopstick
func (self *Chopstick) Clean() {
	self.isClean = true
}

// Member method dirties the chopstick
func (self *Chopstick) Dirty() {
	self.isClean = false
}

/* Our struct type */
type Philosopher struct {
	name       string
	sink       chan *Chopstick
	next       *Philosopher
	meals      int
	done       bool
	chopsticks [2]*Chopstick
}

/* This is a member method, actually does the eating */
func (self *Philosopher) Eat() {
	if self.chopsticks[LEFT] != nil && self.chopsticks[LEFT].isClean &&
		self.chopsticks[RIGHT] != nil && self.chopsticks[RIGHT].isClean {
		self.meals++
		fmt.Println(self.name, "is eating meal number", self.meals)
		self.chopsticks[LEFT].Dirty()
		self.chopsticks[RIGHT].Dirty()
		if self.meals >= NUMBER_OF_MEALS {
			fmt.Println(self.name, "done eating all meals")
			self.done = true
		}
	}
}

/* take a chopstick from our neighbour if we can */
func (self *Philosopher) Take(sink chan *Chopstick) {

	//give up my left chopstick if dirty?
	if self.chopsticks[LEFT] != nil &&
		!self.chopsticks[LEFT].isClean {

		to_send := self.chopsticks[LEFT]
		self.chopsticks[LEFT] = nil
		to_send.Clean()
		sink <- to_send
		//fmt.Println(self.name, "sent left chopstick", to_send)

		return
	}

	//give up my right chopstick
	if self.chopsticks[RIGHT] != nil &&
		!self.chopsticks[RIGHT].isClean {

		to_send := self.chopsticks[RIGHT]
		self.chopsticks[RIGHT] = nil
		to_send.Clean()
		sink <- to_send
		//fmt.Println(self.name, "sent right chopstick")
		return
	}

	//if we get here, we didn't give up a chopstick... so pass along the request
	//only if we are not being greedy (i.e. have no chopstick)
	if self.chopsticks[LEFT] == nil && self.chopsticks[RIGHT] == nil {
		//fmt.Println(self, "passed request on to ", self.next.name)
		self.next.Take(sink)
		return
	}

}

/* This member starts the thread */
func (self *Philosopher) Start() {

	//Some method scoped variables.
	rand.Seed(42)

	//lambda guarentees the main routine knows we are done
	defer (func() {
		close(self.sink)
	})()

	//lambda drives our time out channel
	timeout := make(chan bool, 1)
	go func() {
		for {
			time.Sleep(1 * time.Millisecond)
			timeout <- true
		}
	}()

	//our main loop
	for {
		select {
		case csl := <-self.sink:
			if self.chopsticks[RIGHT] == nil { //make it the left chopstick
				self.chopsticks[RIGHT] = csl
				//fmt.Println(self.name, "received right chopstick:", csl)
			} else if self.chopsticks[LEFT] == nil { //make it the right chopstick
				self.chopsticks[LEFT] = csl
				//fmt.Println(self.name, "received left chopstick:", csl)
			} else { //pass it along
				ptr := self.next
				for ptr.done {
					ptr = ptr.next
				}
				ptr.sink <- csl
			}
			if !self.done {
				self.Eat()
			}
			break

		case <-timeout: //If we are done eating, pass on our self.chopsticks
			if !self.done && (self.chopsticks[LEFT] == nil || self.chopsticks[RIGHT] == nil) {
				self.next.next.Take(self.sink)
			}

			//avoid last man standing deadlock
			if !self.done &&
				self.chopsticks[LEFT] != nil &&
				self.chopsticks[RIGHT] != nil {
				if self.meals == (NUMBER_OF_MEALS-1) &&
					!self.chopsticks[LEFT].isClean ||
					!self.chopsticks[RIGHT].isClean {
					//fmt.Println("Invoking last man standing")
					self.chopsticks[LEFT].Clean()
					self.chopsticks[RIGHT].Clean()
					self.Eat()
				}
			}

			break
		}
	}
}

/* our entryp point into the program */
func main() {

	//Declare our philosophers
	philosophers = []*Philosopher{
		&Philosopher{"Plato", make(chan *Chopstick, NUMBER_OF_MEALS), nil, 0, false, [2]*Chopstick{nil, nil}},
		&Philosopher{"Aristotle", make(chan *Chopstick, NUMBER_OF_MEALS), nil, 0, false, [2]*Chopstick{nil, nil}},
		&Philosopher{"Cicero", make(chan *Chopstick, NUMBER_OF_MEALS), nil, 0, false, [2]*Chopstick{nil, nil}},
		&Philosopher{"Epimendies", make(chan *Chopstick, NUMBER_OF_MEALS), nil, 0, false, [2]*Chopstick{nil, nil}},
		&Philosopher{"Isocrates", make(chan *Chopstick, NUMBER_OF_MEALS), nil, 0, false, [2]*Chopstick{nil, nil}},
	}

	//Spin them all up
	for itr := 0; itr < len(philosophers); itr++ {

		next := itr + 1
		if next >= len(philosophers) {
			next = 0
		} //cant assign this in the constructor
		philosopher := philosophers[itr]
		philosopher.next = philosophers[next]

		go philosopher.Start() //start the philosopher waiting to eat
		cd := false
		if itr%2 == 0 {
			cd = true
		} //whatch what happens when the are all either false or true
		philosopher.sink <- &Chopstick{cd} //give him 1 chopstick (two are required to eat)

		fmt.Println("Initalized ", philosopher.name)

	}

	//wait for everythng to finish eating
	for pos := len(philosophers); pos > 0; {
		x := 0
		for _, philo := range philosophers {
			if !philo.done {
				x++
			}
		}
		if pos > x {
			fmt.Println(x, " Active eating philosophers left.")
		}
		pos = x
		time.Sleep(1 * time.Second)
	}
	fmt.Println("\n************************************\n\nEverybody is done eating ", NUMBER_OF_MEALS, " meals\n")
}
