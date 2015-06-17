**jo** provides a high-performance JSON scanner in Go.

Why? Because the lowest-level JSON parsing primitive provided by the Go
standard library is `json.Unmarshal`ing a byte slice into an `interface{}`
value and inspecting it using runtime reflection... and then crying in the
shower until the pain goes away.

##### Example

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

		// Minify the buffer in-place.
		for r, w = 0, 0; r < n; r++ {
			ev := s.Scan(buf[r])

			// Bail on syntax errors.
			if ev == jo.Error {
				return s.LastError()
			}

			// Ignore whitespace characters.
			if ev&jo.Space == 0 {
				buf[w] = buf[r]
				w++
			}
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
