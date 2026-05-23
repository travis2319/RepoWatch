# Health check
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


