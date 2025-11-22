# Quick Start Guide: Your First Day

**Welcome!** This guide will get you up and running in your first coding session.

**Time needed:** 1-2 hours

---

## Step 1: Install Go (15 minutes)

### Windows

1. Download Go from https://go.dev/dl/
2. Run the installer (`.msi` file)
3. Open Command Prompt and verify:
   ```cmd
   go version
   ```

### macOS

**Option A: Using Homebrew (recommended)**
```bash
brew install go@1.23
```

**Option B: Download installer**
1. Download from https://go.dev/dl/
2. Run the `.pkg` installer
3. Open Terminal and verify:
   ```bash
   go version
   ```

### Linux

```bash
# Download and extract
wget https://go.dev/dl/go1.23.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin

# Verify
go version
```

---

## Step 2: Set Up Your Workspace (10 minutes)

Create your project directory:

```bash
# Navigate to your preferred location
cd ~

# Create project directory
mkdir kuzu-go-learning
cd kuzu-go-learning

# Initialize Go module
go mod init github.com/yourusername/kuzu-go

# Create directory structure
mkdir -p learning/{phase0,phase1,phase2,phase3,phase4}
mkdir -p src/{storage,graph,query,transaction}
```

**Your directory structure should look like:**
```
kuzu-go-learning/
├── go.mod
├── learning/
│   ├── phase0/
│   ├── phase1/
│   ├── phase2/
│   ├── phase3/
│   └── phase4/
└── src/
    ├── storage/
    ├── graph/
    ├── query/
    └── transaction/
```

---

## Step 3: Write Your First Go Program (20 minutes)

### Create hello.go

Create a file `learning/phase0/hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    fmt.Println("I'm learning Go by building a database!")
}
```

### Run it!

```bash
# From your project root
cd learning/phase0
go run hello.go
```

**Expected output:**
```
Hello, World!
I'm learning Go by building a database!
```

### Build an executable

```bash
go build hello.go
./hello  # or hello.exe on Windows
```

**Congratulations!** You just wrote and ran your first Go program! 🎉

---

## Step 4: Simple Calculator (30 minutes)

Let's build something slightly more interesting. Create `learning/phase0/calculator.go`:

```go
package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func subtract(a, b int) int {
    return a - b
}

func multiply(a, b int) int {
    return a * b
}

func divide(a, b int) (int, int) {
    if b == 0 {
        return 0, 0
    }
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}

func main() {
    x := 10
    y := 3

    fmt.Printf("%d + %d = %d\n", x, y, add(x, y))
    fmt.Printf("%d - %d = %d\n", x, y, subtract(x, y))
    fmt.Printf("%d * %d = %d\n", x, y, multiply(x, y))

    q, r := divide(x, y)
    fmt.Printf("%d / %d = %d remainder %d\n", x, y, q, r)
}
```

### Run it:

```bash
go run calculator.go
```

**Expected output:**
```
10 + 3 = 13
10 - 3 = 7
10 * 3 = 30
10 / 3 = 3 remainder 1
```

### Experiment!

Try modifying the code:
1. Change `x` and `y` to different values
2. Add a `square()` function
3. Add a `power()` function

---

## Step 5: Write Your First Test (15 minutes)

Go has testing built-in! Create `learning/phase0/calculator_test.go`:

```go
package main

import "testing"

func TestAdd(t *testing.T) {
    result := add(2, 3)
    expected := 5

    if result != expected {
        t.Errorf("add(2, 3) = %d; want %d", result, expected)
    }
}

func TestSubtract(t *testing.T) {
    result := subtract(10, 3)
    expected := 7

    if result != expected {
        t.Errorf("subtract(10, 3) = %d; want %d", result, expected)
    }
}

func TestMultiply(t *testing.T) {
    result := multiply(4, 5)
    expected := 20

    if result != expected {
        t.Errorf("multiply(4, 5) = %d; want %d", result, expected)
    }
}

func TestDivide(t *testing.T) {
    q, r := divide(10, 3)

    if q != 3 || r != 1 {
        t.Errorf("divide(10, 3) = %d, %d; want 3, 1", q, r)
    }
}
```

### Run the tests:

```bash
go test
```

**Expected output:**
```
PASS
ok      github.com/yourusername/kuzu-go/learning/phase0 0.002s
```

**All tests passed!** ✅

### Run tests with more details:

```bash
go test -v
```

**Expected output:**
```
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
=== RUN   TestSubtract
--- PASS: TestSubtract (0.00s)
=== RUN   TestMultiply
--- PASS: TestMultiply (0.00s)
=== RUN   TestDivide
--- PASS: TestDivide (0.00s)
PASS
ok      github.com/yourusername/kuzu-go/learning/phase0 0.002s
```

---

## Step 6: Set Up Your Editor (10 minutes)

### VS Code (Recommended for beginners)

1. Download from https://code.visualstudio.com/
2. Install the **Go extension**:
   - Open VS Code
   - Press `Ctrl+Shift+X` (or `Cmd+Shift+X` on Mac)
   - Search for "Go"
   - Install the official Go extension by the Go Team

3. Open your project:
   ```bash
   code ~/kuzu-go-learning
   ```

4. Install Go tools when prompted (click "Install All")

### Verify it works:

1. Open `hello.go` in VS Code
2. You should see syntax highlighting
3. Hover over `fmt.Println` - you should see documentation
4. Press `F5` to run with debugger

---

## Step 7: What's Next? (5 minutes)

You've completed the quick start! Here's what you've learned:

- ✅ Installed Go
- ✅ Set up your workspace
- ✅ Wrote and ran a Go program
- ✅ Created functions
- ✅ Wrote and ran tests
- ✅ Set up your editor

### Choose Your Path:

**Complete Beginner?**
→ Start [Phase 0: Go Fundamentals](phase-0-go-fundamentals.md)

**Have some programming experience?**
→ Continue with [Phase 0: Lesson 0.2 (Functions, Packages)](phase-0-go-fundamentals.md#lesson-02-functions-packages-and-modules)

**Already know Go?**
→ Jump to [Phase 1: Storage Layer](phase-1-storage-layer.md)

---

## Troubleshooting

### "go: command not found"

**Problem:** Go is not in your PATH.

**Solution:**
```bash
# Add this to ~/.bashrc or ~/.zshrc
export PATH=$PATH:/usr/local/go/bin

# Then reload your shell
source ~/.bashrc  # or source ~/.zshrc
```

### "cannot find module providing package"

**Problem:** You haven't initialized a Go module.

**Solution:**
```bash
cd ~/kuzu-go-learning
go mod init github.com/yourusername/kuzu-go
```

### Tests fail with "undefined: add"

**Problem:** `calculator_test.go` must be in the same directory as `calculator.go`.

**Solution:**
```bash
# Make sure both files are in the same directory
ls learning/phase0/
# Should show: calculator.go  calculator_test.go
```

### VS Code doesn't show Go syntax highlighting

**Problem:** Go extension not installed.

**Solution:**
1. Press `Ctrl+Shift+X` (Extensions)
2. Search for "Go"
3. Install the official extension
4. Reload window: `Ctrl+Shift+P` → "Reload Window"

---

## Useful Commands Reference

```bash
# Run a Go program
go run filename.go

# Build an executable
go build filename.go

# Run tests
go test

# Run tests with verbose output
go test -v

# Run tests with race detector
go test -race

# Format your code (auto-fixes indentation, spacing)
go fmt

# Check for common mistakes
go vet

# Download dependencies
go mod tidy

# Show Go version
go version

# Show where Go is installed
which go
```

---

## Next Steps

1. **Read the main guide:** [INDEX.md](INDEX.md)
2. **Start Phase 0:** [Go Fundamentals](phase-0-go-fundamentals.md)
3. **Join the community:** Gophers Slack (#newbies channel)

---

## Resources for Day 1

**Interactive Learning:**
- [Tour of Go](https://go.dev/tour/) - 3-4 hours interactive tutorial
- [Go Playground](https://go.dev/play/) - Run Go code in your browser

**Quick References:**
- [Go by Example](https://gobyexample.com/) - Searchable examples
- [Go Cheat Sheet](https://devhints.io/go)

**Videos:**
- [Go in 100 Seconds](https://www.youtube.com/watch?v=446E-r0rXHI) - Quick overview
- [Learn Go Programming](https://www.youtube.com/watch?v=YS4e4q9oBaU) - Full beginner course

---

**Congratulations on completing your first day!** 🎉

Remember: Everyone starts as a beginner. Take it one step at a time, experiment, break things, and most importantly—have fun!

**Ready for more?** → [Phase 0: Go Fundamentals](phase-0-go-fundamentals.md)
