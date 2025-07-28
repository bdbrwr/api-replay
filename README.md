<p align="center">
  <img src="assets/api-replay.png" alt="API Replay logo">
</p>

<p align="center">
  <a href="https://github.com/bdbrwr/api-replay/stargazers">
    <img src="https://img.shields.io/github/stars/bdbrwr/api-replay?style=social" alt="GitHub stars">
  </a>
  <a href="https://github.com/bdbrwr/api-replay/issues">
    <img src="https://img.shields.io/github/issues/bdbrwr/api-replay?color=blue" alt="GitHub issues">
  </a>
  <a href="https://www.linkedin.com/in/bdbrwr/">
    <img src="https://img.shields.io/badge/LinkedIn-Profile-blue?logo=linkedin&style=flat" alt="LinkedIn">
  </a>
  <a href="https://x.com/bdbrwr">
    <img src="https://img.shields.io/badge/Twitter-@bdbrwr-1DA1F2?logo=twitter&style=flat" alt="Twitter">
  </a>
  <a href="https://blog.boot.dev/news/hackathon-2025/">
    <img src="https://img.shields.io/badge/Boot.dev-Hackathon-A77618?logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzAiIGhlaWdodD0iMzAiIHZpZXdCb3g9IjAgMCAzMCAzMCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMu...." alt="Boot.dev Hackathon">
  </a>
</p>

# API-Replay

**API-Replay** is a developer tool that lets you record and serve API responses locally - ideal for demos, offline use, and mocking private or authenticated APIs.


🎯 Record real API responses  
🔁 Replay them over HTTP from your local machine  
📦 Easily ship demos without exposing secrets or requiring live APIs  

> This project was originally built for the [Boot.Dev Hackaton](https://boot.dev/) and I will likely evolve this beyond the event as my capstone project for the programme

> If you like this project, please hit it with a star ⭐

---

## ✨ Use Cases

- Show off a project that depends on private or rate-limited APIs
- Cache responses behind an API key for development/testing
- Demo a frontend app offline with static data
- Share reproducible API responses for bug reports or testing

---

## ✅ Features

- Supports `GET` request recording with custom headers
- Strip base URLs for clean, logical folder structures
- Query parameter handling with safe file encoding
- Human-readable `.json` response files
- Custom output directory support
- Serve responses over HTTP on configurable ports
- Easy to script or automate (great in CI pipelines)

---

## 🔥 Common APIs That Could Benefit

| API Type             | Usefulness for Emulation | Notes                                                                 |
|----------------------|---------------------------|-----------------------------------------------------------------------|
| Game APIs         | ✅✅✅                      | Blizzard, Riot, Steam, Twitch — ideal for stable demos or dev work   |
| SaaS / REST APIs  | ✅✅✅                      | Stripe, GitHub, Shopify — great for mocking predictable responses     |
| 🛠️ Internal APIs  | ✅✅✅                      | Great for isolating frontend work from backend logic and auth layers  |
| Mail APIs         | ✅✅                        | SendGrid, Mailgun — useful for simulating confirmation flows          |
| Public JSON APIs  | ✅                          | IP lookup, weather, country APIs — good for offline testing           |
| LLM APIs (OpenAI) | ❌                          | Too large & dynamic — token-based & expensive per request             |
| Docs APIs         | ✅ (partial)                | Notion, Confluence — works for static page preview/testing            |
| Shipping APIs     | ✅✅                        | FedEx, UPS — useful for static label/tracking testing                 |
| Analytics APIs    | ✅                          | Good for dashboards and static report development                     |

> ⚠️ Sorry, we cannot replay your AI-powered, streaming, multi-modal OpenAI wrapper with fallback retries and fine-tuned models. Yet.


---

## 🛠️ Installation

We recommend to install into bin, so api-replay can be called from anywhere

```bash
go install github.com/bdbrwr/api-replay@latest
```

Alternatively, you can clone the repo and build it. 
```bash
git clone https://github.com/bdbrwr/api-replay.git
cd api-replay
go build ./cmd/api-replay
```

---

## 🔧 Configuration
The CLI expects `api-replay.yaml` to be present in the directory you run the commands from. 

```yaml
dir: api-replay
port: "1337"
```
The CLI will use these defaults unless overridden

---

## 🧪 Recording API Responses
```bash
api-replay record https://pokeapi.co/api/v2/location-area
```
*stores the json in ./api-replay/api/v2/location-area.json*

#### With Base Url stripped (recommended when you have some defined in your code)
``` bash
api-replay record https://pokeapi.co/api/v2/location-area \
  -B https://pokeapi.co/api/v2
```
*stores the json in ./api-replay/location-area.json*

#### Override the configured output folder for specific requests
``` bash
api-replay record https://pokeapi.co/api/v2/location-area \
  -O demo-cache
```
*stores the json in ./demo-cache/api/v2/location-area.json*

### Add Custom Headers
``` bash
api-replay record https://example.com/protected \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/json"
```

>Note: If the API requires authentication, you are expected to manually retrieve any access tokens using tools like curl or Postman, and then inject them via headers as shown above. api-replay does not currently handle OAuth flows internally.

---
## 💻 Serving Recorded Responses

```bash
api-replay serve
```

#### Or from a specific folder and port
```bash
api-replay serve -D demo-cache -P 3000
```

To access your mocked API at
`http://localhost:3000/<original-path>`

---
## 📁 File Structure
Recorded files follow the structure of the original API path:
```
api-replay/
├── v2/
│   ├── items.json
│   └── users@limit%3D10.json
```

Query parameters are safely encoded in the filenames
>For example:
`/items?limit=10` → `items@limit%3D10.json`

---

## 🧰 Built With

- [Cobra](https://github.com/spf13/cobra) – CLI framework
- [Viper](https://github.com/spf13/viper) – Config management
- [Chi](https://github.com/go-chi/chi) – Lightweight HTTP router

--- 
## 🤝 Contributing 
Feel free to fork and contribute via PRs!
Ideas welcome in [Issues](https://github.com/bdbrwr/api-replay/issues)
