// test662 : stdlib5.legsup gen4kdb

package legsup

import (
	"errors"
	"fmt"
	"stdlib5/kio"
	"strconv"
	"strings"
)

// gen4 kdb data node
type G4data struct {
	ByteV  []byte
	StrV   string
	IntV   int
	FloatV float64

	Dtype rune // 'b', 's', 'i', 'f'
}

// gen4 kdb data set
func (node *G4data) Set(data interface{}) error {
	switch v := data.(type) {
	case []byte:
		node.ByteV = v
		node.Dtype = 'b'
	case string:
		node.StrV = v
		node.Dtype = 's'
	case int:
		node.IntV = v
		node.Dtype = 'i'
	case float64:
		node.FloatV = v
		node.Dtype = 'f'
	default:
		node.Dtype = 'n'
		return errors.New("InvalidData")
	}
	return nil
}

// gen4 kdb read db, returns nil if error
func G4DBread(raw string) map[string]G4data {
	out := make(map[string]G4data)
	temp := strings.Split(strings.ToUpper(raw), "\n")
	for _, r := range temp {
		if len(r) < 5 {
			continue
		}
		pos0 := strings.Index(r, "(")
		pos1 := strings.Index(r, ")")
		switch r[pos0+1 : pos1] {
		case "BYTES":
			var new G4data
			tb, err := kio.Bread(r[pos1+1:])
			if err == nil {
				new.Set(tb)
				out[r[0:pos0]] = new
			}
		case "STR":
			var new G4data
			tb, err := kio.Bread(r[pos1+1:])
			if err == nil {
				new.Set(string(tb))
				out[r[0:pos0]] = new
			}
		case "INT":
			var new G4data
			ti, err := strconv.Atoi(r[pos1+1:])
			if err == nil {
				new.Set(ti)
				out[r[0:pos0]] = new
			}
		case "FLOAT":
			var new G4data
			tf, err := strconv.ParseFloat(r[pos1+1:], 64)
			if err == nil {
				new.Set(tf)
				out[r[0:pos0]] = new
			}
		default:
			return nil
		}
	}
	return out
}

// gen4 kdb write db
func G4DBwrite(data map[string]G4data) string {
	out := make([]string, 0)
	for i, r := range data {
		switch r.Dtype {
		case 'b':
			out = append(out, fmt.Sprintf("%s(BYTES)%s", i, kio.Bprint(r.ByteV)))
		case 's':
			out = append(out, fmt.Sprintf("%s(STR)%s", i, kio.Bprint([]byte(r.StrV))))
		case 'i':
			out = append(out, fmt.Sprintf("%s(INT)%d", i, r.IntV))
		case 'f':
			out = append(out, fmt.Sprintf("%s(FLOAT)%.4f", i, r.FloatV))
		}
	}
	return strings.ToUpper(strings.Join(out, "\n"))
}
