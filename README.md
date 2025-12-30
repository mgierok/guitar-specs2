# Guitar-Specs

Guitar-Specs is a web app for searching, viewing, and comparing guitars. The frontend is built with Next.js and Tailwind CSS, and the backend API is written in Go. The site targets strong SEO, Core Web Vitals performance, and accessibility compliance (including the European Accessibility Act).

## Architecture
- `frontend/` — Next.js (React) app
- `backend/` — Go API
- `shared/` — cross-platform types/utilities for future React Native reuse

## External Dependencies
Required (local build and run):
- Docker Engine or Docker Desktop
- Docker Compose v2
- Git
- PostgreSQL (self-hosted for now)

Optional (for local dev outside containers):
- Node.js (latest LTS) + npm
- Go (latest stable)

Required (AWS deployment):
- AWS account and IAM credentials
- AWS CLI v2

Recommended (AWS infrastructure as code):
- Terraform

## Local Build and Run (Containers)
The app is built and run in containers. Use Docker Compose to start the full stack.

1) Create environment files:
```bash
cp frontend/.env.example frontend/.env
cp backend/.env.example backend/.env
```

2) Build and start services:
```bash
docker compose up --build
```

3) Open the app:
- Frontend: `http://localhost:3000`
- API: `http://localhost:8080`

4) Apply migrations (once Postgres is running):
```bash
docker compose exec db psql -U postgres -d guitar_specs -f /migrations/0001_init.up.sql
```

## Local Testing (Containers)
Run tests inside the containers:

```bash
docker compose run --rm frontend npm test
docker compose run --rm backend go test ./...
```

## AWS Deployment (Containers)
The production site is served at `https://www.guitar-specs.com/` behind Cloudflare CDN. The recommended AWS deployment is ECS on EC2 (lower cost baseline) with an Application Load Balancer (ALB) and ECR for images.

1) Configure AWS CLI:
```bash
aws configure
```

2) Create ECR repositories:
```bash
aws ecr create-repository --repository-name guitar-specs-frontend
aws ecr create-repository --repository-name guitar-specs-backend
```

3) Build and push images:
```bash
# Frontend
aws ecr get-login-password --region <region> | \
  docker login --username AWS --password-stdin <account>.dkr.ecr.<region>.amazonaws.com

docker build -t guitar-specs-frontend -f frontend/Dockerfile frontend

docker tag guitar-specs-frontend:latest \
  <account>.dkr.ecr.<region>.amazonaws.com/guitar-specs-frontend:latest

docker push <account>.dkr.ecr.<region>.amazonaws.com/guitar-specs-frontend:latest

# Backend

docker build -t guitar-specs-backend -f backend/Dockerfile backend

docker tag guitar-specs-backend:latest \
  <account>.dkr.ecr.<region>.amazonaws.com/guitar-specs-backend:latest

docker push <account>.dkr.ecr.<region>.amazonaws.com/guitar-specs-backend:latest
```

4) Provision AWS infrastructure (required via Terraform):
- VPC, subnets, security groups
- ECS cluster on EC2 (Auto Scaling group) + ECS services
- ALB with HTTPS (ACM certificate)
- CloudWatch logs
- Parameter Store or Secrets Manager for env vars

5) Point Cloudflare DNS to the ALB and enable HTTPS.

## Configuration
Document required environment variables in `docs/config.md` as the API and web app evolve. Keep secrets out of the repo and store them in AWS SSM or Secrets Manager for production.

## Test Data
Canonical test dataset: `data/guitars.json`. Keep it aligned with `DATABASE.md`.

## API Endpoints
Base path: `/api/v1`
- `GET /health` — service health check
- `GET /guitars` — list guitars (pagination/filtering to be wired)
- `GET /guitars/{slug}` — guitar detail

## CI
GitHub Actions runs frontend lint and backend tests on PRs and on pushes to `dev`.
