tech-ip-sem2/
  services/
    auth/
      cmd/auth/main.go
      internal/
        http/...
        service/...
      Dockerfile (опционально на этом ПЗ)
    tasks/
      cmd/tasks/main.go
      internal/
        http/...
        service/...
        client/authclient/...
  shared/
    middleware/
      requestid.go
      logging.go
    httpx/
      client.go
  docs/
    pz17_api.md
    pz17_diagram.md (опционально)
  README.md
