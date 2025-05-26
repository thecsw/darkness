#!/bin/bash
set -e

# Create temp directory with test file
TEST_DIR=$(mktemp -d)
mkdir -p $TEST_DIR/content

# Create a test org file with nested elements
cat > $TEST_DIR/content/test.org <<EOL
#+title: Test Nested Source Code in Details
#+date: 2023-05-14

* Testing nested elements

#+begin_details Click to see code
This is text inside a details block.

#+begin_src go
func main() {
	fmt.Println("Hello from inside details!")
}
#+end_src

More text within details.

#+begin_quote
A quote inside details
#+end_quote

#+end_details

Normal paragraph outside details.
EOL

# Create a minimal darkness config
cat > $TEST_DIR/darkness.toml <<EOL
name = "Test Site"
[author]
name = "Test Author"
EOL

# Build the site
echo "Building test site..."
cd $TEST_DIR && go run /home/runner/work/darkness/darkness/darkness.go build

# Check the generated HTML
echo "Checking generated HTML..."
if grep -A 20 "<details>" $TEST_DIR/out/test.html | grep -q "<div class=\"coding\""; then
    echo "SUCCESS: Source code block is inside details block!"
    exit 0
else
    echo "FAILURE: Source code block is not inside details block!"
    exit 1
fi