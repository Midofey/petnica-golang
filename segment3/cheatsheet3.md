### Pointers
```golang
// a pointer holds a memory address
var x int = 42
p := &x           // p points to x

// read the value at the address
fmt.Println(*p)   // 42

// change the value through the pointer
*p = 100
fmt.Println(x)   // 100 - x changed!

// pointer type is written as *T
var q *int = &x
```

### Passing pointers to functions
```golang
// without pointer - original unchanged
func triple(n int) {
    n = n * 3       // only changes the copy
}

// with pointer - original is modified
func triplePtr(n *int) {
    *n = *n * 3    // modifies the original
}

x := 5
triple(x)         // x is still 5
triplePtr(&x)     // x is now 15
```

### Defining structs
```golang
type Player struct {
    Name   string
    Score  int
    Active bool
}

// create with field names (recommended)
p := Player{
    Name:   "alice",
    Score:  100,
    Active: true,
}

// create with positional values
p2 := Player{"bob", 80, true}

// access fields
fmt.Println(p.Name)
p.Score = 120
```

### Methods on structs
```golang
// value receiver — works on a copy
func (p Player) describe() string {
    return fmt.Sprintf(
        "%s: %d pts", p.Name, p.Score,
    )
}

// pointer receiver — modifies the original
func (p *Player) addScore(points int) {
    p.Score += points
}

// calling methods
fmt.Println(p.describe())
p.addScore(50)
fmt.Println(p.Score)  // 170
```

### Value receiver vs pointer receiver - when to use which
```golang
// use a VALUE receiver when the method only reads fields — it does not modify the struct
func (p Player) isActive() bool {
    return p.Active
}

// use a POINTER receiver when the method needs to change one or more fields
func (p *Player) deactivate() {
    p.Active = false
}

// Go automatically takes the address when needed — both of these work:
p.deactivate()    // Go rewrites this as (&p).deactivate()
```

### Structs in slices and maps
```golang
// slice of structs
players := []Player{
    {Name: "alice", Score: 100, Active: true},
    {Name: "bob",   Score: 80,  Active: true},
}

for _, p := range players {
    fmt.Println(p.describe())
}

// map of structs
roster := map[string]Player{
    "alice": {Name: "alice", Score: 100, Active: true},
}
```

### Interfaces
```golang
// this interface exists in Go natively
type Stringer interface {
    String() string
}
```

```golang
// if your struct has a method with this exact
// signature, Go treats it as its string representation
func (p Player) String() string {
    return fmt.Sprintf(
        "%s | HP: %d | Score: %d",
        p.Name, p.Health, p.Score,
    )
}

// fmt.Println will call String() automatically
p := Player{Name: "alice", Health: 100, Score: 200}
fmt.Println(p)
// alice | HP: 100 | Score: 200

// so will %v and %s in Printf
fmt.Printf("Player: %v\n", p)
// Player: alice | HP: 100 | Score: 200
```

```golang
// without String(), fmt.Println prints
// the raw struct fields like this:
fmt.Println(p)
// {alice 100 200 1}
```
