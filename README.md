# Decode Apollo config to strcut field

Install via:

```bash
go install github.com/caiyunapp/oap
```

Usage like:

```go
import "github.com/caiyunapp/oap"

type DemoConfig struct {
	Int      int           `apollo:"int"`
	Int8     int8          `apollo:"int8"`
	Int16    int16         `apollo:"int16"`
	Int32    int32         `apollo:"int32"`
	Int64    int64         `apollo:"int64"`
	Uint     uint          `apollo:"uint"`
	Uint8    uint8         `apollo:"uint8"`
	Uint16   uint16        `apollo:"uint16"`
	Uint32   uint32        `apollo:"uint32"`
	Uint64   uint64        `apollo:"uint64"`
	Duration time.Duration `apollo:"duration"`
}

func main(){
	// init your apollo client here
	// ...

	conf := &DemoConfig{}
	if err := oap.Decode(conf, client, make(map[string][]agollo.OpOption)); err != nil {
		panic(err)
	}
}
```

Support types:

- [x] String
- [x] Int
- [x] Bool
- [x] Float32
- [x] Float64
- [x] Struct from JSON or YAML
