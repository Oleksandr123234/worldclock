package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/gin-gonic/gin"
)

type zoneInfo struct {
	Name       string
	OffsetMins int
}

const pageTemplate = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>World Clock</title>
  <style>
    body { font-family: -apple-system, Segoe UI, Roboto, sans-serif; background: #0f1115; color: #e8e8e8; display: flex; justify-content: center; padding-top: 6vh; }
    .card { background: #181b21; border-radius: 12px; padding: 32px 40px; box-shadow: 0 8px 24px rgba(0,0,0,.4); min-width: 420px; }
    h1 { margin: 0 0 20px; font-size: 20px; color: #9ad1ff; }
    table { width: 100%; border-collapse: collapse; }
    td, th { padding: 10px 12px; text-align: left; }
    th { color: #8a8f98; font-weight: 500; font-size: 13px; text-transform: uppercase; border-bottom: 1px solid #2a2f38; }
    tr:not(:last-child) td { border-bottom: 1px solid #22262e; }
    .time { font-variant-numeric: tabular-nums; font-size: 18px; color: #7ee787; }
  </style>
</head>
<body>
  <div class="card">
    <h1>🕒 World Clock</h1>
    <table>
      <tr><th>Timezone</th><th>Local Time</th></tr>
      {{range .Zones}}
      <tr>
        <td>{{.Name}}</td>
        <td class="time" data-offset="{{.OffsetMins}}">--:--:--</td>
      </tr>
      {{end}}
    </table>
  </div>
  <script>
    function pad(n) { return n.toString().padStart(2, "0"); }
    function tick() {
      const nowUTC = Date.now();
      document.querySelectorAll(".time").forEach(function (el) {
        const offsetMins = parseInt(el.dataset.offset, 10);
        const local = new Date(nowUTC + offsetMins * 60000);
        el.textContent = pad(local.getUTCHours()) + ":" + pad(local.getUTCMinutes()) + ":" + pad(local.getUTCSeconds());
      });
    }
    tick();
    setInterval(tick, 1000);
  </script>
</body>
</html>
`

func loadZones() []zoneInfo {
	raw := strings.TrimSpace(os.Getenv("TIMEZONES"))
	if raw == "" {
		log.Fatal("TIMEZONES environment variable is required (comma-separated IANA timezone names, e.g. \"UTC,Europe/Kyiv,America/New_York\")")
	}

	var zones []zoneInfo
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		loc, err := time.LoadLocation(name)
		if err != nil {
			log.Fatalf("invalid timezone %q: %v", name, err)
		}
		_, offsetSec := time.Now().In(loc).Zone()
		zones = append(zones, zoneInfo{Name: name, OffsetMins: offsetSec / 60})
	}

	if len(zones) == 0 {
		log.Fatal("TIMEZONES resolved to an empty list after parsing")
	}
	return zones
}

func setupRouter(zones []zoneInfo) *gin.Engine {
	tmpl := template.Must(template.New("index").Parse(pageTemplate))

	r := gin.Default()
	r.SetHTMLTemplate(tmpl)

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{"Zones": zones})
	})

	return r
}

func main() {
	zones := loadZones()
	r := setupRouter(zones)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("worldclock listening on :%s with zones: %s", port, os.Getenv("TIMEZONES"))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
