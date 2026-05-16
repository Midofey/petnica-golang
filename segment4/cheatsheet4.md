### The value,err pattern
```golang
// many standard library functions return
// a result AND an error
n, err := strconv.Atoi("42")
if err != nil {
    fmt.Println("something went wrong:", err)
} else {
    fmt.Println("parsed:", n)
}

// if you don't care about the error (rare!)
n, _ := strconv.Atoi("42")

// nil means "no error occurred"
// non-nil means something went wrong
```

### Returning errors from functions
```golang
import "errors"

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

// calling it
result, err := divide(10, 0)
if err != nil {
    fmt.Println("error:", err)
    return
} fmt.Println("result:", result)
```

## Wrapping errors with fmt.Errorf
```golang
// fmt.Errorf creates a formatted error message
// use it when you want to include context or variable values
func (p *Player) takeDamage(damage int) error {
    if damage < 0 {
        return fmt.Errorf(
            "takeDamage: damage cannot be negative, got %d", damage,
        )
    }
    p.Health -= damage
    if p.Health < 0 {
        p.Health = 0
    }
    return nil
}

// wrapping an existing error with %w adds context
// while preserving the original error
func (p *Player) resolveAttack(damage int) error {
    err := p.takeDamage(damage)
    if err != nil {
        return fmt.Errorf("resolveAttack: %w", err)
    }
    return nil
}

// the error message now reads:
// "resolveAttack: takeDamage: damage cannot be negative, got -5"
```

### Custom error types
```golang
// a custom error is just a struct that
// has an Error() string method
type InvalidDamageError struct {
    Damage int
    Reason string
}

func (e InvalidDamageError) Error() string {
    return fmt.Sprintf(
        "invalid damage %d: %s", e.Damage, e.Reason,
    )
}

// use it in takeDamage instead of fmt.Errorf
func (p *Player) takeDamage(damage int) error {
    if damage < 0 {
        return InvalidDamageError{
            Damage: damage,
            Reason: "damage cannot be negative",
        }
    }
    p.Health -= damage
    if p.Health < 0 {
        p.Health = 0
    }
    return nil
}

// this is another example of an interface —
// anything with Error() string satisfies
// the error interface in Go
```