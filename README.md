# dining_philosophers
A Go'ish rendition of the dining philosophers problem.

So I wrote this to do a presentation on Go for some coworkers, and thought I would share. 

This is a cooperative solution (patterned on Chandy/Misra, 1984) to the problem (no auditor) with 5 diners and 5 chopsticks. I choose chopsticks because I felt it made more sense to need 2 of them rather than 2 forks.

Since each diner is independent no locks or mutex's are used and therefore the diners must rely on state to determine whether or not they can eat. This can be represented in several ways; hunger, plate ready, utensils clean or dirty, etc. I choose to keep state in the chopsticks as dirty and clean.

In this solution, channels are used to send / take forks from other diners. The main variant here is the introduction of the last man standing deadlock, which all the other diners are done, and the last guy is waiting for a clean chopstick... the is remediation is a forced clean in the event a diner is on his last meal, and nobody is expected to send a clean chopstick his way since they are all done to.

More on this solution and the problem can be found at:

https://en.wikipedia.org/wiki/Dining_philosophers_problem

The purpose for posting this application is for demonstration purposes only and it's not fit for production use. I wrote it in one evening to go along with a set of slides I presented to some co workers, it did not go through any type of peer review, unit testing, etc... if there isn't a bug in it, I'll be shocked! On that note, please don't bug me with update requests... it shows off Go in some nifty ways and that's all it's designed to do. This also means you can't use it for your CS
class either!

Hope you enjoy....
