### Defining and calling functions
```golang
// no parameters, no return value
func greet() {
    fmt.Println("hello!")
}

// with parameters
func greetName(name string) {
    fmt.Printf("hello, %s!\n", name)
}

// with a return value
func add(a int, b int) int {
    return a + b
}

// calling functions
greet()
greetName("alex")
result := add(3, 4)
```

### Multiple return values
```golang
// return two values
func minMax(a, b int) (int, int) {
    if a < b {
        return a, b
    }
    return b, a
}

// capturing both values
min, max := minMax(10, 3)

// ignoring one value with _
_, max := minMax(10, 3)

// common pattern: value + error
n, err := strconv.Atoi("42")
if err != nil {
    fmt.Println("something went wrong")
}
```

### Slices
```golang
// declare and initialise
nums := []int{1, 2, 3}
names := []string{"alice", "bob"}

// declare empty, then append
var scores []int
scores = append(scores, 95)
scores = append(scores, 87)

// indexing (zero-based)
first := nums[0]
last  := nums[len(nums)-1]

// length
n := len(nums) // 3
```

### Maps
```golang
// declare and initialise
ages := map[string]int{
    "alice": 17,
    "bob":   16,
}

// declare empty, then add entries
scores := make(map[string]int)
scores["alice"] = 95

// reading a value
a := ages["alice"]       // 17

// check if a key exists
val, ok := ages["carol"]
if !ok {
    fmt.Println("not found")
}

// delete an entry
delete(ages, "bob")
```

### For range
```golang
// range over a slice — index and value
for i, v := range nums {
    fmt.Printf("index %d: %d\n", i, v)
}

// ignore the index
for _, v := range nums {
    fmt.Println(v)
}

// range over a map — key and value (order not guaranteed)
for key, val := range ages {
    fmt.Printf("%s is %d\n", key, val)
}

// ignore the value — keys only
for key := range ages {
    fmt.Println(key)
}
```

### Reading numbers and strings one after another
```golang
var reader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

//read ints as strings, then convert
func readInt(prompt string) int {
	for {
		input := readLine(prompt)
		n, err := strconv.Atoi(input)
		if err == nil {
			return n
		}
		fmt.Println("Please enter a valid number.")
	}
}
```