package ichika

import (
	"fmt"
)

// HelpCommandFunc shows default darkness HelpCommandFunc message
func HelpCommandFunc() {
	fmt.Println(`My name is Darkness.
My calling is that of a crusader.
Do Shometing Gwazy!

If you don't have a darkness website yet, start with
creating it with new followed by the directory name

  $> darkness new axel

Here are the commands you can use, -help is supported:
  build - build the entire directory
  serve - build HTTP and serve them
  megumin - blow up the directory!!
  clean - megumin but super boring
  misa - supercharge your website
  lalatina - pls dont
  aqua - ...

Don't hold back! You have no choice!`)
}
