# Resume

Edit `resume.md` for content and `template.html` for layout and styles.
The first Markdown heading supplies the page title. Raw HTML is disabled.

Requires Go 1.22 or newer.

```sh
./rebuild.sh                    # Generate index.html for GitHub Pages
go run ./cmd/build -serve      # Preview at http://127.0.0.1:8080
```

Refresh the browser after editing. Use browser print / Save as PDF for a
printable copy. Commit `index.html` alongside source changes.

```sh
go test ./...
go vet ./...
```
