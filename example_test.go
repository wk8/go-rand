package rand_test

import (
	"fmt"
	"github.com/wk8/go-rand"
)

func Example() {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}

	rand.Seed(42)

	fmt.Println("Magic 8-Ball says:")

	for i := 0; i < 3; i++ {
		fmt.Println(answers[rand.Intn(len(answers))])
	}

	// let's save the state, so we can then rewind to here
	state, err := rand.Marshall()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		fmt.Println(answers[rand.Intn(len(answers))])
	}

	// reset the state
	if err := rand.Unmarshall(state); err != nil {
		panic(err)
	}

	// the next 3 should be the exact same as the last 3
	for i := 0; i < 3; i++ {
		fmt.Println(answers[rand.Intn(len(answers))])
	}

	// Output:
	// Magic 8-Ball says:
	// As I see it yes
	// Outlook good
	// Yes
	// Reply hazy try again
	// Yes definitely
	// As I see it yes
	// Reply hazy try again
	// Yes definitely
	// As I see it yes
}
