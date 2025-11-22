# Pre-work Week 1-2: Go Fundamentals

**Duration:** 2 weeks | **Time Commitment:** 15-20 hours/week | **Difficulty:** Beginner

## Welcome to Go!

This is your first step towards building a graph database. Don't worry if you're completely new to programming or Go—this lesson assumes zero prior knowledge.

By the end of these two weeks, you'll:
- ✅ Understand Go's basic syntax and type system
- ✅ Write functions and use structs
- ✅ Create a working CLI application
- ✅ Understand interfaces (Go's superpower)
- ✅ Handle errors properly

## Week 1: The Basics

### Day 1: Hello, Go!

#### Setup

First, verify Go is installed:

```bash
go version  # Should show go1.25 or later
```

#### Your First Program

Create a directory and your first file:

```bash
mkdir -p ~/kuzu-go-learning/week1
cd ~/kuzu-go-learning/week1
```

Create `hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

Run it:

```bash
go run hello.go
```

**Understanding what's happening:**
- `package main` - Every Go program needs a main package
- `import "fmt"` - Bring in the formatting package
- `func main()` - The entry point of your program
- `fmt.Println()` - Print with a newline

#### Interactive Exercise 1: Variables

```go
package main

import "fmt"

func main() {
    // Variable declaration
    var name string = "Alice"
    var age int = 30

    // Short declaration (more common)
    city := "New York"
    isActive := true

    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
    fmt.Println("City:", city)
    fmt.Println("Active:", isActive)
}
```

**Try this:**
1. Change the values
2. Add a new variable for `height` (use `float64`)
3. Print it out

**Key concept:** `:=` is "declare and assign" and can infer the type.

#### Interactive Exercise 2: Basic Types

```go
package main

import "fmt"

func main() {
    // Numbers
    var integer int = 42
    var float float64 = 3.14159

    // Strings
    var message string = "Hello, World!"

    // Booleans
    var isTrue bool = true

    // Type conversions (explicit!)
    var x int = 10
    var y float64 = float64(x)  // Must explicitly convert

    fmt.Printf("Integer: %d\n", integer)
    fmt.Printf("Float: %.2f\n", float)
    fmt.Printf("String: %s\n", message)
    fmt.Printf("Boolean: %t\n", isTrue)
    fmt.Printf("Converted: %.1f\n", y)
}
```

**Important:** Go does NOT allow implicit type conversions. `var y float64 = x` would fail.

**Try this:**
1. Try to assign an `int` to a `float64` without conversion—watch it fail
2. Use `fmt.Printf` with different format verbs: `%d` (decimal), `%f` (float), `%s` (string), `%t` (boolean)

### Day 2: Functions

#### Basic Function

```go
package main

import "fmt"

// Simple function
func greet(name string) {
    fmt.Println("Hello,", name)
}

// Function with return value
func add(a int, b int) int {
    return a + b
}

// Multiple return values (very common in Go!)
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}

func main() {
    greet("Alice")

    sum := add(5, 3)
    fmt.Println("Sum:", sum)

    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", result)
    }
}
```

**Key concepts:**
- Functions can return multiple values
- Error handling uses the pattern `value, err := function()`
- Check `if err != nil` after every function that returns an error
- `nil` is Go's "null" or "none"

#### Interactive Exercise 3: Write Your Own Functions

Create `calculator.go`:

```go
package main

import "fmt"

func multiply(a, b int) int {
    // TODO: Implement multiplication
    return 0
}

func isEven(n int) bool {
    // TODO: Return true if n is even
    return false
}

func greetPerson(name string, age int) string {
    // TODO: Return a string like "Hello, Alice! You are 30 years old."
    return ""
}

func main() {
    fmt.Println(multiply(4, 5))        // Should print 20
    fmt.Println(isEven(7))             // Should print false
    fmt.Println(greetPerson("Bob", 25)) // Should print "Hello, Bob! You are 25 years old."
}
```

**Your task:** Implement these three functions. Run with `go run calculator.go`.

### Day 3: Control Flow

#### If Statements

```go
package main

import "fmt"

func main() {
    age := 18

    // Simple if
    if age >= 18 {
        fmt.Println("You are an adult")
    }

    // If-else
    if age < 13 {
        fmt.Println("Child")
    } else if age < 18 {
        fmt.Println("Teenager")
    } else {
        fmt.Println("Adult")
    }

    // If with initialization statement (very Go!)
    if score := 85; score >= 90 {
        fmt.Println("Grade: A")
    } else if score >= 80 {
        fmt.Println("Grade: B")
    }
    // Note: 'score' is only in scope within the if block
}
```

#### For Loops (The Only Loop in Go!)

```go
package main

import "fmt"

func main() {
    // Classic for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    // While-style loop
    count := 0
    for count < 3 {
        fmt.Println("Count:", count)
        count++
    }

    // Infinite loop (use break to exit)
    counter := 0
    for {
        if counter >= 3 {
            break
        }
        fmt.Println("Iteration:", counter)
        counter++
    }

    // Range over slice
    numbers := []int{10, 20, 30}
    for index, value := range numbers {
        fmt.Printf("Index: %d, Value: %d\n", index, value)
    }

    // If you don't need the index
    for _, value := range numbers {
        fmt.Println(value)
    }
}
```

**Key concept:** The blank identifier `_` discards values you don't need.

#### Interactive Exercise 4: FizzBuzz

Classic programming challenge:

```go
package main

import "fmt"

func main() {
    // Print numbers 1-20
    // If divisible by 3: print "Fizz"
    // If divisible by 5: print "Buzz"
    // If divisible by both: print "FizzBuzz"
    // Otherwise: print the number

    for i := 1; i <= 20; i++ {
        // TODO: Implement FizzBuzz
    }
}
```

**Your task:** Implement FizzBuzz. Expected output:
```
1
2
Fizz
4
Buzz
Fizz
...
```

### Day 4: Arrays and Slices

#### Arrays (Fixed Size)

```go
package main

import "fmt"

func main() {
    // Array declaration
    var arr [5]int  // Array of 5 integers, initialized to zeros
    arr[0] = 10
    arr[1] = 20

    // Array literal
    numbers := [5]int{1, 2, 3, 4, 5}

    fmt.Println(arr)
    fmt.Println(numbers)
    fmt.Println("Length:", len(numbers))
}
```

**Important:** Arrays have fixed size. `[5]int` and `[10]int` are different types!

#### Slices (Dynamic, More Common)

```go
package main

import "fmt"

func main() {
    // Slice declaration
    var s []int  // nil slice
    fmt.Println(s, len(s), cap(s))  // [] 0 0

    // Make a slice
    s = make([]int, 5)  // length 5, capacity 5
    s[0] = 10

    // Slice literal
    fruits := []string{"apple", "banana", "cherry"}

    // Append (grows the slice)
    fruits = append(fruits, "date")
    fmt.Println(fruits)

    // Slicing
    subset := fruits[1:3]  // Elements at index 1 and 2
    fmt.Println(subset)    // [banana cherry]

    // Length vs Capacity
    nums := make([]int, 3, 5)  // length 3, capacity 5
    fmt.Printf("Len: %d, Cap: %d\n", len(nums), cap(nums))

    nums = append(nums, 10, 20)
    fmt.Printf("Len: %d, Cap: %d\n", len(nums), cap(nums))  // 5, 5

    nums = append(nums, 30)  // Capacity will grow (usually doubles)
    fmt.Printf("Len: %d, Cap: %d\n", len(nums), cap(nums))  // 6, 10 (or similar)
}
```

**Key concepts:**
- **Length:** Number of elements currently in the slice
- **Capacity:** Total space allocated
- `append` may allocate new memory if capacity is exceeded
- Slices are references to underlying arrays

#### Interactive Exercise 5: Slice Manipulation

```go
package main

import "fmt"

func main() {
    numbers := []int{1, 2, 3, 4, 5}

    // TODO: Find the sum of all numbers
    sum := 0
    // Your code here

    fmt.Println("Sum:", sum)  // Should be 15

    // TODO: Find the maximum number
    max := numbers[0]
    // Your code here

    fmt.Println("Max:", max)  // Should be 5

    // TODO: Create a new slice with only even numbers
    var evens []int
    // Your code here

    fmt.Println("Evens:", evens)  // Should be [2, 4]
}
```

### Day 5: Maps

```go
package main

import "fmt"

func main() {
    // Make a map
    ages := make(map[string]int)
    ages["Alice"] = 30
    ages["Bob"] = 25

    // Map literal
    scores := map[string]int{
        "Alice": 95,
        "Bob":   87,
        "Carol": 92,
    }

    // Access
    fmt.Println("Alice's score:", scores["Alice"])

    // Check if key exists
    score, exists := scores["David"]
    if exists {
        fmt.Println("David's score:", score)
    } else {
        fmt.Println("David not found")
    }

    // Delete
    delete(scores, "Bob")

    // Iterate
    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }
}
```

**Key concepts:**
- Maps are key-value pairs (like dictionaries in Python, objects in JavaScript)
- Always check if a key exists with `value, ok := map[key]`
- Iteration order is random (not guaranteed)

#### Interactive Exercise 6: Word Counter

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    text := "the quick brown fox jumps over the lazy dog the fox"
    words := strings.Fields(text)  // Split by whitespace

    // TODO: Count how many times each word appears
    // Hint: Use a map[string]int

    // Your code here

    // Expected output:
    // the: 3
    // quick: 1
    // brown: 1
    // fox: 2
    // jumps: 1
    // over: 1
    // lazy: 1
    // dog: 1
}
```

### Day 6-7: Structs

#### Defining Structs

```go
package main

import "fmt"

// Define a struct type
type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    // Create a struct
    alice := Person{
        Name: "Alice",
        Age:  30,
        City: "New York",
    }

    // Short form (must be in order)
    bob := Person{"Bob", 25, "San Francisco"}

    // Access fields
    fmt.Println(alice.Name)
    fmt.Println(bob.Age)

    // Modify
    alice.Age = 31
    fmt.Println(alice)
}
```

#### Methods on Structs

```go
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method with pointer receiver (can modify)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}

    fmt.Println("Area:", rect.Area())  // 50

    rect.Scale(2)
    fmt.Println("Scaled:", rect)  // {20 10}
    fmt.Println("New Area:", rect.Area())  // 200
}
```

**Key concepts:**
- **Value receiver** `(r Rectangle)`: Receives a copy, can't modify original
- **Pointer receiver** `(r *Rectangle)`: Receives a reference, can modify original
- Use pointer receivers when you need to modify the struct or when the struct is large

#### Interactive Exercise 7: Bank Account

```go
package main

import "fmt"

type BankAccount struct {
    Owner   string
    Balance float64
}

// TODO: Implement Deposit method
func (b *BankAccount) Deposit(amount float64) {
    // Add amount to balance
}

// TODO: Implement Withdraw method
func (b *BankAccount) Withdraw(amount float64) bool {
    // Subtract amount from balance
    // Return false if insufficient funds
    return false
}

// TODO: Implement GetBalance method
func (b *BankAccount) GetBalance() float64 {
    return 0
}

func main() {
    account := BankAccount{Owner: "Alice", Balance: 100}

    account.Deposit(50)
    fmt.Println("Balance:", account.GetBalance())  // Should be 150

    success := account.Withdraw(30)
    fmt.Println("Withdraw success:", success)  // Should be true
    fmt.Println("Balance:", account.GetBalance())  // Should be 120

    success = account.Withdraw(200)
    fmt.Println("Withdraw success:", success)  // Should be false
    fmt.Println("Balance:", account.GetBalance())  // Should still be 120
}
```

## Week 2: Intermediate Concepts

### Day 8-9: Interfaces

Interfaces are one of Go's most powerful features.

```go
package main

import "fmt"

// Interface definition
type Shape interface {
    Area() float64
}

// Rectangle implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

// Function that accepts any Shape
func printArea(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circ := Circle{Radius: 7}

    printArea(rect)  // Works!
    printArea(circ)  // Also works!
}
```

**Key concepts:**
- Interfaces define behavior (methods)
- Types implement interfaces **implicitly** (no "implements" keyword)
- If a type has all the methods, it implements the interface
- This is called "duck typing": if it walks like a duck and quacks like a duck...

#### Common Interfaces

```go
package main

import (
    "fmt"
    "strings"
)

// Stringer interface (from fmt package)
type Person struct {
    Name string
    Age  int
}

// Implementing fmt.Stringer
func (p Person) String() string {
    return fmt.Sprintf("%s (age %d)", p.Name, p.Age)
}

func main() {
    alice := Person{Name: "Alice", Age: 30}
    fmt.Println(alice)  // Automatically uses String() method
}
```

#### Interactive Exercise 8: Animal Sounds

```go
package main

import "fmt"

// Define an Animal interface
type Animal interface {
    Speak() string
}

// TODO: Create a Dog struct
type Dog struct {
    Name string
}

// TODO: Implement Speak for Dog (return "Woof!")

// TODO: Create a Cat struct
type Cat struct {
    Name string
}

// TODO: Implement Speak for Cat (return "Meow!")

func makeAnimalSpeak(a Animal) {
    fmt.Println(a.Speak())
}

func main() {
    dog := Dog{Name: "Buddy"}
    cat := Cat{Name: "Whiskers"}

    makeAnimalSpeak(dog)  // Should print: Woof!
    makeAnimalSpeak(cat)  // Should print: Meow!
}
```

### Day 10: Error Handling

Go uses explicit error handling, not exceptions.

```go
package main

import (
    "errors"
    "fmt"
)

// Function that returns an error
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{
            Field:   "age",
            Message: "cannot be negative",
        }
    }
    if age > 150 {
        return &ValidationError{
            Field:   "age",
            Message: "unrealistic value",
        }
    }
    return nil
}

func main() {
    // Simple error handling
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)

    // Custom error
    if err := validateAge(-5); err != nil {
        fmt.Println(err)
    }
}
```

**Key concepts:**
- Errors are values that implement the `error` interface
- Always check `if err != nil`
- Return errors, don't panic (except for truly unrecoverable situations)

### Day 11-12: Packages and Imports

#### Creating Your Own Package

Create directory structure:
```
week2/
├── main.go
└── calculator/
    └── calculator.go
```

`calculator/calculator.go`:
```go
package calculator

// Exported function (starts with capital letter)
func Add(a, b int) int {
    return a + b
}

// Unexported function (starts with lowercase)
func multiply(a, b int) int {
    return a * b
}

// Exported function using unexported function
func Square(n int) int {
    return multiply(n, n)
}
```

`main.go`:
```go
package main

import (
    "fmt"
    "week2/calculator"  // Import your package
)

func main() {
    sum := calculator.Add(5, 3)
    square := calculator.Square(4)

    fmt.Println("Sum:", sum)      // 8
    fmt.Println("Square:", square) // 16

    // This would fail (unexported):
    // calculator.multiply(2, 3)
}
```

Initialize a Go module:
```bash
cd week2
go mod init week2
go run main.go
```

**Key concepts:**
- Capital letter = exported (public)
- Lowercase = unexported (private to package)
- Packages organize code

### Day 13-14: Practice Project - TODO CLI

Build a complete command-line TODO application!

**Features:**
- Add tasks
- List tasks
- Mark tasks as complete
- Delete tasks
- Save to file

Create `main.go`:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Task struct {
    ID        int
    Title     string
    Completed bool
}

type TodoList struct {
    Tasks  []Task
    NextID int
}

func (t *TodoList) Add(title string) {
    task := Task{
        ID:        t.NextID,
        Title:     title,
        Completed: false,
    }
    t.Tasks = append(t.Tasks, task)
    t.NextID++
    fmt.Printf("Added task #%d: %s\n", task.ID, task.Title)
}

func (t *TodoList) List() {
    if len(t.Tasks) == 0 {
        fmt.Println("No tasks!")
        return
    }

    for _, task := range t.Tasks {
        status := " "
        if task.Completed {
            status = "✓"
        }
        fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Title)
    }
}

func (t *TodoList) Complete(id int) {
    for i := range t.Tasks {
        if t.Tasks[i].ID == id {
            t.Tasks[i].Completed = true
            fmt.Printf("Completed task #%d\n", id)
            return
        }
    }
    fmt.Printf("Task #%d not found\n", id)
}

func (t *TodoList) Delete(id int) {
    for i := range t.Tasks {
        if t.Tasks[i].ID == id {
            t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
            fmt.Printf("Deleted task #%d\n", id)
            return
        }
    }
    fmt.Printf("Task #%d not found\n", id)
}

func main() {
    todo := &TodoList{NextID: 1}
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println("TODO List - Commands: add, list, complete, delete, quit")

    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }

        input := scanner.Text()
        parts := strings.Fields(input)

        if len(parts) == 0 {
            continue
        }

        command := parts[0]

        switch command {
        case "add":
            if len(parts) < 2 {
                fmt.Println("Usage: add <task title>")
                continue
            }
            title := strings.Join(parts[1:], " ")
            todo.Add(title)

        case "list":
            todo.List()

        case "complete":
            if len(parts) < 2 {
                fmt.Println("Usage: complete <task id>")
                continue
            }
            id, err := strconv.Atoi(parts[1])
            if err != nil {
                fmt.Println("Invalid task ID")
                continue
            }
            todo.Complete(id)

        case "delete":
            if len(parts) < 2 {
                fmt.Println("Usage: delete <task id>")
                continue
            }
            id, err := strconv.Atoi(parts[1])
            if err != nil {
                fmt.Println("Invalid task ID")
                continue
            }
            todo.Delete(id)

        case "quit":
            fmt.Println("Goodbye!")
            return

        default:
            fmt.Println("Unknown command:", command)
        }
    }
}
```

**Try it:**
```bash
go run main.go
> add Buy groceries
> add Write code
> list
> complete 1
> list
> delete 2
> quit
```

#### Your Enhancement Tasks

1. **Save to file:** Add `Save()` and `Load()` methods to save tasks to a JSON file
2. **Priorities:** Add a priority field (high, medium, low)
3. **Due dates:** Add due date tracking
4. **Search:** Add ability to search tasks by keyword

## Week 1-2 Checkpoint

### Self-Assessment Quiz

Answer these without looking at the lessons:

1. What's the difference between `var x int` and `x := 0`?
2. How do you check if a key exists in a map?
3. What's the difference between an array and a slice?
4. When should you use a pointer receiver vs value receiver on a method?
5. How does a type implement an interface in Go?
6. What does `_` mean in `for _, value := range items`?

### Practical Test

Complete this challenge without help:

**Task:** Build a simple phonebook application that:
- Stores name -> phone number mappings
- Has `Add()`, `Get()`, `Delete()`, and `List()` operations
- Handles errors (e.g., adding duplicate names)
- Uses a struct to encapsulate the data
- Implements a `String()` method for pretty printing

If you can do this in 2-3 hours, you're ready for Week 3-4!

## What's Next?

You've learned:
- ✅ Variables, types, and basic syntax
- ✅ Functions and control flow
- ✅ Slices and maps (Go's most important data structures)
- ✅ Structs and methods
- ✅ Interfaces (Go's superpower)
- ✅ Error handling
- ✅ Packages

**Next:** [Week 3-4: Intermediate Concepts](week-3-4-intermediate-concepts.md)

You'll learn:
- Deep dive into pointers and memory
- Testing and benchmarking
- Working with JSON
- Building a persistent key-value store

## Resources

**If you're stuck:**
- [Tour of Go](https://go.dev/tour/) - Interactive tutorial
- [Go by Example](https://gobyexample.com/) - Code examples
- [Gophers Slack](https://invite.slack.golangbridge.org/) - #newbies channel

**If you want more practice:**
- [Exercism Go Track](https://exercism.org/tracks/go) - 100+ exercises
- [Go Programming by Example](https://www.youtube.com/playlist?list=PLzUGFf4GhXBL4GHXVcMMvzgtO8-WEJIoY)

Take your time, practice the exercises, and don't hesitate to ask for help!
