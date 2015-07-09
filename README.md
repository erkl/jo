Package **jo** provides a simple, high-performance JSON scanner for Go. Why?
Because the lowest level parsing functionality provided by the standard library
is the `encoding/json` package's `json.Unmarshal` function, which is clearly
too high level for many use cases.


#### Example

Below is a function which minifies raw JSON while it's being copied.

```go
func minify(dst io.Writer, src io.Reader) error {
	var buf = make([]byte, 4096)
	var s = jo.NewScanner()
	var w, r int

	for {
		// Read the next chunk of data.
		n, err := src.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Minify the buffer in place.
		for r, w = 0, 0; r < n; r++ {
			ev := s.Scan(buf[r])

			// Bail on syntax errors.
			if ev == jo.Error {
				return s.LastError()
			}

			// Ignore whitespace characters.
			if ev&jo.Space != 0 {
				continue
			}

			buf[w] = buf[r]
			w++
		}

		// Write the now compressed buffer.
		_, err = dst.Write(buf[:w])
		if err != nil {
			return err
		}
	}

	// Check for syntax errors caused by incomplete values.
	if ev := s.End(); ev == jo.Error {
		return s.LastError()
	}

	return nil
}
```


#### License

```
Copyright (c) 2015, Erik Lundin.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE
OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
```
