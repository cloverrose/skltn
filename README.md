# skltn

Reduce routine work by generating skeleton from template and method's signature

## How to use

### Build

```sh
go build cmd/skltn/skltn.go
```

### Generate skeleton
- `skltn` use clipboard for data in/out.
- You need to copy method signature before you run `skltn`
- e.g. `func (d *db) Update(userID data.UserID, value data.UserDTO) error {`

```sh
./skltn -t toy
```

- Then your clipboard contains generated text

e.g.
```text
Hello db
```

### Options
- `-t TemplateName` (required) use registered `TemplateName`
- `-f` (optional, default false) enable formatting source code. If template is not complete go syntax it cause error.
