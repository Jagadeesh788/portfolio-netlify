package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Article struct {
	Title   string
	Date    string
	Slug    string
	Content template.HTML
	Summary string
}

type SocialLink struct {
	Name string
	URL  string
}

type PageData struct {
	Title       string
	Articles    []Article
	Skills      []string
	SocialLinks []SocialLink
}

func main() {
	os.MkdirAll("docs", 0755)
	os.MkdirAll("docs/articles", 0755)

	generateHomePage()
	generateArticles()
	copyAssets()

	fmt.Println("Portfolio site generated in 'docs' directory!")
}

func generateHomePage() {
	skills := []string{
		"Go/Golang", "Microservices", "Docker", "Kubernetes",
		"REST APIs", "gRPC", "Kafka", "Redis",
		"AWS", "Git", "OpenTelemetry", "CI/CD",
		"PostgreSQL", "Python", "C/C++",
	}

	articles := getArticles()

	socialLinks := []SocialLink{
		{Name: "GitHub", URL: "https://github.com/Jagadeesh788"},
		{Name: "LinkedIn", URL: "https://www.linkedin.com/in/jagadeesh-patil-lnkdn"},
	}

	data := PageData{
		Title:       "Jagadeesh Patil",
		Articles:    articles,
		Skills:      skills,
		SocialLinks: socialLinks,
	}

	tmpl := template.Must(template.New("index").Parse(indexTemplate))
	file, _ := os.Create("docs/index.html")
	defer file.Close()
	tmpl.Execute(file, data)
}

func generateArticles() {
	articles := getArticles()
	tmpl := template.Must(template.New("article").Parse(articleTemplate))

	for _, article := range articles {
		file, _ := os.Create(fmt.Sprintf("docs/articles/%s.html", article.Slug))
		tmpl.Execute(file, article)
		file.Close()
	}
}

func getArticles() []Article {
	var articles []Article

	filepath.WalkDir("content", func(path string, d fs.DirEntry, err error) error {
		if err != nil || !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, _ := os.ReadFile(path)
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(html.WithUnsafe()),
		)
		var buf bytes.Buffer
		md.Convert(content, &buf)
		htmlContent := buf.String()

		filename := strings.TrimSuffix(filepath.Base(path), ".md")
		title := strings.ReplaceAll(filename, "-", " ")
		title = strings.Title(title)

		articles = append(articles, Article{
			Title:   title,
			Date:    time.Now().Format("Jan 2, 2006"),
			Slug:    filename,
			Content: template.HTML(htmlContent),
			Summary: extractSummary(string(content)),
		})
		return nil
	})

	return articles
}

func extractSummary(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			if len(line) > 150 {
				return line[:150] + "..."
			}
			return line
		}
	}
	return "Read more..."
}

func copyAssets() {
	os.MkdirAll("docs/images", 0755)
	if data, err := os.ReadFile("profile.jpeg"); err == nil {
		os.WriteFile("docs/images/profile.jpeg", data, 0644)
	}

	css := `
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }
header { background: #1a1a2e; color: white; padding: 2rem 0; }
.hero { display: flex; align-items: center; gap: 2rem; }
.profile-img { width: 150px; height: 150px; border-radius: 50%; object-fit: cover; border: 4px solid white; }
.hero-text { flex: 1; }
.hero h1 { font-size: 3rem; margin-bottom: 0.5rem; }
.hero p { font-size: 1.2rem; opacity: 0.9; }
.skills { background: #f8f9fa; padding: 3rem 0; }
.skills-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; margin-top: 2rem; }
.skill-card { background: white; padding: 1rem; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; }
.articles { padding: 3rem 0; }
.articles-grid { display: grid; gap: 2rem; margin-top: 2rem; }
.article-card { border: 1px solid #ddd; border-radius: 8px; padding: 1.5rem; }
.article-card h3 { color: #1a1a2e; margin-bottom: 0.5rem; }
.article-card .date { color: #666; font-size: 0.9rem; margin-bottom: 1rem; }
.btn { display: inline-block; background: #e94560; color: white; padding: 0.5rem 1rem; text-decoration: none; border-radius: 4px; margin-top: 1rem; }
.btn:hover { background: #c73652; }
.social-links { display: flex; justify-content: center; gap: 2rem; margin-bottom: 1rem; }
.social-links a { color: white; text-decoration: none; padding: 0.5rem 1rem; border: 1px solid rgba(255,255,255,0.3); border-radius: 4px; transition: background 0.3s; }
.social-links a:hover { background: rgba(255,255,255,0.1); }
footer { background: #1a1a2e; color: white; text-align: center; padding: 2rem 0; margin-top: 3rem; }
.content h1, .content h2, .content h3 { margin: 1.5rem 0 0.75rem; color: #1a1a2e; }
.content p { margin-bottom: 1rem; }
.content ul, .content ol { margin: 1rem 0 1rem 2rem; }
.content li { margin-bottom: 0.4rem; }
.content pre { background: #f4f4f4; padding: 1rem; border-radius: 6px; overflow-x: auto; margin: 1rem 0; }
.content code { background: #f4f4f4; padding: 0.2rem 0.4rem; border-radius: 3px; font-family: monospace; font-size: 0.9em; }
.content pre code { background: none; padding: 0; }
@media (max-width: 768px) { .hero { flex-direction: column; text-align: center; } .hero h1 { font-size: 2rem; } .skills-grid { grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); } }
`
	os.WriteFile("docs/style.css", []byte(css), 0644)
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <header>
        <div class="container">
            <div class="hero">
                <img src="images/profile.jpeg" alt="Jagadeesh Patil" class="profile-img">
                <div class="hero-text">
                    <h1>Jagadeesh Patil</h1>
                    <p>Backend Engineer & Software Developer</p>
                    <p style="font-size: 1rem; margin-top: 0.5rem;">Bengaluru | jagadeeshpatil2244@gmail.com</p>
                </div>
            </div>
        </div>
    </header>

    <section class="skills">
        <div class="container">
            <h2>Technical Skills</h2>
            <div class="skills-grid">
                {{range .Skills}}
                <div class="skill-card">
                    <strong>{{.}}</strong>
                </div>
                {{end}}
            </div>
        </div>
    </section>

    <section class="articles">
        <div class="container">
            <h2>Latest Articles</h2>
            <div class="articles-grid">
                {{range .Articles}}
                <article class="article-card">
                    <h3>{{.Title}}</h3>
                    <div class="date">{{.Date}}</div>
                    <p>{{.Summary}}</p>
                    <a href="articles/{{.Slug}}.html" class="btn">Read More</a>
                </article>
                {{end}}
            </div>
        </div>
    </section>

    <footer>
        <div class="container">
            <div class="social-links">
                {{range .SocialLinks}}
                <a href="{{.URL}}" target="_blank">{{.Name}}</a>
                {{end}}
            </div>
            <p>&copy; 2026 Jagadeesh Patil. Built with Go.</p>
        </div>
    </footer>
</body>
</html>`

const articleTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="../style.css">
</head>
<body>
    <header>
        <div class="container">
            <h1><a href="../index.html" style="color: white; text-decoration: none;">← Back to Portfolio</a></h1>
        </div>
    </header>

    <main class="container" style="padding: 3rem 0;">
        <article>
            <h1>{{.Title}}</h1>
            <div class="date" style="margin-bottom: 2rem;">{{.Date}}</div>
            <div class="content">{{.Content}}</div>
        </article>
    </main>

    <footer>
        <div class="container">
            <p>&copy; 2026 Jagadeesh Patil. Built with Go.</p>
        </div>
    </footer>
</body>
</html>`
