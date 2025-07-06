# Code Review

## What I Found

### 1. No Input Validation

The code doesn't check if the `name` parameter is valid:

- What if `name` is empty? You'll create a user with no name
- What if someone passes malicious input(e.g. SQL injection)?
- What if the name is way too long?

### 2. Wrong HTTP Method

This looks like it's meant to be called with a GET request (using query parameters), but creating users should use POST with a proper request body.

### 3. Poor Error Handling

There's literally no error handling anywhere:

- If something goes wrong, you'll never know
- Users get a success response even if the operation failed
- Debugging becomes a nightmare

### 4. Data Race Condition (Critical Issue)

The biggest problem here is that the `users` map can be accessed by multiple goroutines at the same time without any protection. This will cause:

- **Data races** - Multiple goroutines writing to the map simultaneously
- **Runtime panics** - Go's map implementation isn't thread-safe
- **Lost data** - Writes might get overwritten or corrupted

### 5. Lack of Response Content

The handler immediately returns `200 OK` but the actual work happens in a goroutine. This means:

- You're telling the client everything worked before you even tried to do anything
- If the user creation fails, the client will never know
- You can't provide meaningful feedback

### 6. Unnecessary goroutine

- Launching a goroutine for every request `(go createUser(...))` can lead to unbounded goroutines under high load, potentially exhausting memory or CPU.
- If createUser grows in complexity, this becomes riskier.

## How to fix it?

```go

// ----------------------------------- User -----------------------------------

// here should be use `struct` to replace simple map
type User struct {
	ID   string
	Name string
	CreateAt int64
}

type UserStore struct{
	mu sync.RWMutex
	UserList map[string]User
}

func (us *UserStore) CreateUser(name string) error {
    us.mu.Lock()
    defer us.mu.Unlock()

    if _, ok := us.UserList[name]; ok {
        return errors.New("user already exists")
    }

    // Safe to modify the map now
    us.UserList[name] = User{
    	ID: uuid.New().String(),
     	Name: name,
     	CreatedAt: time.Now().Unix(),
    }
    return nil
}

var users = UserStore{
	UserList: make(map[string]User),
}

// ----------------------------------- validate name -----------------------------------
func validateName(name string) error {
    if strings.TrimSpace(name) == `` {
        return errors.New("name cannot be empty")
    }
    if len(name) > 32 {
        return errors.New("name too long")
    }
    // regexp to validate username
    if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
        return errors.New("invalid username")
    }
    return nil
}

// ----------------------------------- handler -----------------------------------
// change name to handleCreateUser() here
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct{
    		Name string `json:"name"`
        // can add more fields here ...
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := users.CreateUser(req.Name); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated) // replace with ok status code, cause `created` should be 201 here
    json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

```
