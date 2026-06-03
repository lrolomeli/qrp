package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const dashboardHTML = `<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>QR Dashboard</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;background:#f5f5f5;color:#333;padding:2rem}
.container{max-width:1200px;margin:0 auto}
h1{margin-bottom:1.5rem;color:#111}
.cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:1rem;margin-bottom:2rem}
.card{background:#fff;border-radius:8px;padding:1.5rem;box-shadow:0 1px 3px rgba(0,0,0,0.1)}
.card h3{font-size:.875rem;color:#666;margin-bottom:.5rem;text-transform:uppercase;letter-spacing:.05em}
.card .value{font-size:2rem;font-weight:700;color:#111}
table{width:100%;border-collapse:collapse;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,0.1);margin-bottom:2rem}
th,td{padding:.75rem 1rem;text-align:left;border-bottom:1px solid #eee}
th{background:#fafafa;font-weight:600;font-size:.875rem;color:#666;text-transform:uppercase;letter-spacing:.05em}
tr:hover{background:#fafafa}
.code{font-family:'SF Mono',Monaco,monospace;font-weight:600;color:#2563eb}
.dest{color:#666;font-size:.875rem;max-width:300px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ip{font-family:monospace;font-size:.875rem;color:#666}
.ua{font-size:.875rem;color:#666;max-width:250px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.ts{font-size:.875rem;color:#888}
h2{margin:1.5rem 0 1rem;color:#222}
a{color:#2563eb;text-decoration:none}
a:hover{text-decoration:underline}
</style>
</head>
<body>
<div class="container">
<h1>QR Redirect Dashboard</h1>

<div class="cards">
  <div class="card">
    <h3>QR Codes</h3>
    <div class="value">{{.CodeCount}}</div>
  </div>
  <div class="card">
    <h3>Total Clicks</h3>
    <div class="value">{{.TotalClicks}}</div>
  </div>
</div>

<h2>Clicks por c&oacute;digo</h2>
<table>
<thead>
  <tr>
    <th>C&oacute;digo</th>
    <th>Destino</th>
    <th>Clicks</th>
  </tr>
</thead>
<tbody>
  {{range .ClicksByCode}}
  <tr>
    <td class="code">{{.Code}}</td>
    <td class="dest" title="{{.DestinationURL}}">{{.DestinationURL}}</td>
    <td>{{.Clicks}}</td>
  </tr>
  {{end}}
</tbody>
</table>

<h2>&Uacute;ltimos clicks</h2>
<table>
<thead>
  <tr>
    <th>Fecha</th>
    <th>C&oacute;digo</th>
    <th>Destino</th>
    <th>IP</th>
    <th>User-Agent</th>
  </tr>
</thead>
<tbody>
  {{range .RecentClicks}}
  <tr>
    <td class="ts">{{.CreatedAt}}</td>
    <td class="code">{{.Code}}</td>
    <td class="dest" title="{{.DestinationURL}}">{{.DestinationURL}}</td>
    <td class="ip">{{.IP}}</td>
    <td class="ua" title="{{.UserAgent}}">{{.UserAgent}}</td>
  </tr>
  {{end}}
</tbody>
</table>
</div>
</body>
</html>`

type DashboardData struct {
	CodeCount    int
	TotalClicks  int
	ClicksByCode []ClickCount
	RecentClicks []Click
}

func main() {
	redirectURL := getEnv("REDIRECT_URL", "https://ejemplo.com")
	port := getEnv("PORT", "8080")

	store, err := NewStore(getEnv("DB_PATH", "/data/qrp.db"))
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	defer store.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		go store.RecordClick("root", redirectURL, r.RemoteAddr, r.UserAgent(), r.Referer())
		http.Redirect(w, r, redirectURL, http.StatusFound)
	})

	mux.HandleFunc("GET /api/dashboard", func(w http.ResponseWriter, r *http.Request) {
		totalClicks, err := store.GetTotalClicks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		clicksByCode, err := store.GetClicksByCode()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		recentClicks, err := store.GetRecentClicks(50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := DashboardData{
			CodeCount:    1,
			TotalClicks:  totalClicks,
			ClicksByCode: clicksByCode,
			RecentClicks: recentClicks,
		}

		tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("QR backend listening on :%s — redirecting to %s", port, redirectURL)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
