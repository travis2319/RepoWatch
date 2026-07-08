## API Testing with curl

### 1) Health check
```bash
curl http://localhost:4000/health
```

### 2) Validate GitHub token
```bash
curl http://localhost:4000/api/v1/github/validate
```

### 3) Get all repos for an owner
```bash
curl "http://localhost:4000/api/v1/repos?owner=travis2319"
```

### 4) Check a user across all repos of an owner
```bash
curl -X POST http://localhost:4000/api/v1/check \
  -H "Content-Type: application/json" \
  -d '{"owner":"travis2319","user":"ChetanNaikk"}'
```

### 5) Check a user in a single repository
```bash
curl -X POST http://localhost:4000/api/v1/check-single \
  -H "Content-Type: application/json" \
  -d '{"owner":"travis2319","repo":"RepoWatch","user":"ChetanNaikk"}'
```

### 6) Legacy single-repo check endpoint
```bash
curl http://localhost:4000/check-single
```

<!-- # Health check
curl http://localhost:4000/health

# Get all repos of an owner
curl "http://localhost:4000/api/v1/repos?owner=travis2319"

# Check a user across ALL repos of an owner
curl -X POST http://localhost:4000/api/v1/check \
  -H "Content-Type: application/json" \
  -d '{"owner": "travis2319", "user": "ChetanNaikk"}'

# Check a user in ONE specific repo
curl -X POST http://localhost:4000/api/v1/check-single \
  -H "Content-Type: application/json" \
  -d '{"owner": "travis2319", "repo": "RepoWatch", "user": "ChetanNaikk"}'

 -->

