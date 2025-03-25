# Deployment Guide

This guide provides detailed instructions for deploying the Strategic Brand Optimization Platform to various environments. The platform uses a microservice architecture with Docker containers, making it adaptable to different deployment scenarios.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Deployment Options](#deployment-options)
3. [Supabase Setup](#supabase-setup)
4. [Local Deployment](#local-deployment)
5. [Docker Swarm Deployment](#docker-swarm-deployment)
6. [Kubernetes Deployment](#kubernetes-deployment)
7. [AWS Deployment](#aws-deployment)
8. [Environment Variables](#environment-variables)
9. [Health Monitoring](#health-monitoring)
10. [Scaling Considerations](#scaling-considerations)
11. [Backup and Recovery](#backup-and-recovery)
12. [Troubleshooting](#troubleshooting)

## Prerequisites

Before deploying the Strategic Brand Optimization Platform, ensure you have:

- Docker and Docker Compose (for local/Swarm deployments)
- Kubernetes cluster (for Kubernetes deployments)
- A Supabase account with a project created
- Required API keys for social media platforms (if using the Scraper service)
- A domain name for production deployments
- SSL certificates for secure communication

## Deployment Options

The platform supports several deployment options:

1. **Local Development**: Using Docker Compose for development and testing
2. **Docker Swarm**: For simple production deployments with basic orchestration
3. **Kubernetes**: For production deployments with advanced orchestration
4. **AWS ECS**: For Amazon Web Services deployments

Choose the option that best fits your operational requirements and expertise.

## Supabase Setup

The platform uses Supabase for data storage and authentication. Follow these steps to set up Supabase:

1. Create a Supabase account at [supabase.com](https://supabase.com)
2. Create a new project
3. Go to SQL Editor and run the database schema:
   - Copy the schema from `docs/db/supabase_schema.sql`
   - Run the script in the SQL Editor
4. Set up authentication:
   - Go to Authentication > Settings
   - Configure email provider settings
   - Enable "Email confirmations" and "Secure email change"
5. Get your API keys:
   - Go to Project Settings > API
   - Note down the `anon` public key and `service_role` secret key
   - Add these to your environment variables during deployment

## Local Deployment

For local development or small deployments:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/go-competitor.git
   cd go-competitor
   ```

2. Create a `.env` file with required environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your Supabase credentials and other settings
   ```

3. Start the services with Docker Compose:
   ```bash
   docker-compose up -d
   ```

4. Verify the deployment:
   ```bash
   # Check if all services are running
   docker-compose ps
   
   # Check the health of a specific service
   curl http://localhost:9001/health
   
   # Test the GraphQL API
   curl -X POST -H "Content-Type: application/json" --data '{"query": "{ __schema { types { name } } }"}' http://localhost:8080/query
   ```

## Docker Swarm Deployment

For production deployments using Docker Swarm:

1. Initialize a Docker Swarm (if not already done):
   ```bash
   docker swarm init
   ```

2. Create a Docker config for environment variables:
   ```bash
   cat .env | docker config create go-competitor-env -
   ```

3. Deploy the stack:
   ```bash
   docker stack deploy -c docker-compose.prod.yaml go-competitor
   ```

4. Verify the deployment:
   ```bash
   docker stack services go-competitor
   ```

5. Scale specific services:
   ```bash
   docker service scale go-competitor_graphql=3
   ```

## Kubernetes Deployment

For production deployments using Kubernetes:

1. Apply the Kubernetes manifests:
   ```bash
   # Apply namespace
   kubectl apply -f kubernetes/namespace.yaml
   
   # Apply secrets (first create a secrets file from template)
   cp kubernetes/secrets.yaml.template kubernetes/secrets.yaml
   # Edit kubernetes/secrets.yaml with your credentials
   kubectl apply -f kubernetes/secrets.yaml
   
   # Apply all other resources
   kubectl apply -f kubernetes/
   ```

2. Verify the deployment:
   ```bash
   kubectl get pods -n go-competitor
   ```

3. Enable the Ingress controller (if required):
   ```bash
   kubectl apply -f kubernetes/ingress.yaml
   ```

4. Scale specific deployments:
   ```bash
   kubectl scale deployment/graphql -n go-competitor --replicas=3
   ```

## AWS Deployment

For deployment on Amazon Web Services:

1. Set up an AWS ECR repository for each service:
   ```bash
   aws ecr create-repository --repository-name go-competitor/graphql
   aws ecr create-repository --repository-name go-competitor/auth
   # Repeat for each service
   ```

2. Build and push images to ECR:
   ```bash
   # Log in to ECR
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin YOUR_AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com
   
   # Build and push each service
   docker build -t YOUR_AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/go-competitor/graphql:latest -f graphql/Dockerfile .
   docker push YOUR_AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/go-competitor/graphql:latest
   # Repeat for each service
   ```

3. Deploy using CloudFormation or Terraform:
   - Use `aws/cloudformation.yaml` for CloudFormation
   - Use `aws/terraform/` for Terraform

4. Configure AWS ECS:
   - Create a cluster
   - Create task definitions for each service
   - Create services with the desired number of tasks
   - Set up load balancing with AWS Application Load Balancer

## Environment Variables

The platform requires these environment variables:

### Common Variables (required by all services)
| Variable | Description | Required |
|----------|-------------|----------|
| `SUPABASE_URL` | Supabase instance URL | Yes |
| `SUPABASE_ANON_KEY` | Supabase anon key | Yes |
| `SUPABASE_SERVICE_ROLE` | Supabase service role key | Yes |

### Service-Specific Variables
| Service | Variable | Description | Required |
|---------|----------|-------------|----------|
| Auth | `JWT_SECRET` | Secret for signing JWT tokens | Yes |
| GraphQL | `GRAPHQL_PORT` | Port to listen on | No (default: 8080) |
| Scraper | `FB_APP_ID` | Facebook App ID | Only for Facebook scraping |
| Scraper | `FB_APP_SECRET` | Facebook App Secret | Only for Facebook scraping |
| Notification | `SMTP_HOST` | SMTP server for emails | Only for email notifications |

Refer to each service's documentation for a complete list of variables.

## Health Monitoring

All services expose a health check endpoint at `/health`, which returns:

```json
{"status":"UP"}
```

Use these endpoints to:
1. Configure container health checks
2. Set up external monitoring
3. Enable load balancer health probes

For more comprehensive monitoring:
- Set up Prometheus for metrics collection
- Set up Grafana for visualization
- Configure alerts based on service health

## Scaling Considerations

When scaling the platform, consider:

1. **Stateless Services**: Most services are stateless and can be scaled horizontally
2. **Database Connection Pooling**: Configure connection pools appropriately
3. **Resource Allocation**: Allocate resources based on service requirements
4. **API Rate Limits**: Be aware of external API rate limits for the Scraper service
5. **Load Balancing**: Ensure proper load balancing for scaled services

### Recommended Scaling Strategy

| Service | Scaling Recommendation |
|---------|------------------------|
| GraphQL | Scale horizontally based on traffic |
| Auth | Minimal scaling (2-3 instances for high availability) |
| Analytics | Scale based on processing requirements |
| Scraper | Scale carefully, consider API rate limits |

## Backup and Recovery

### Database Backup

Supabase provides automated backups, but you can also:

1. Set up daily PostgreSQL dumps:
   ```bash
   pg_dump -h your-supabase-db.supabase.co -U postgres -d postgres > backup.sql
   ```

2. Store backups securely:
   - Use encrypted storage
   - Set up rotation policies
   - Test restore procedures regularly

### Service Recovery

For service recovery:

1. Each service can be redeployed independently
2. Service dependencies will reconnect automatically
3. Use rolling updates to avoid downtime

## Troubleshooting

Common issues and solutions:

### Service Won't Start

**Symptoms**: Container exits immediately or fails to become healthy

**Solutions**:
- Check logs: `docker logs <container_id>`
- Verify environment variables are set correctly
- Ensure Supabase is accessible
- Check for port conflicts

### Services Can't Communicate

**Symptoms**: "Connection refused" errors in logs

**Solutions**:
- Verify service discovery is working
- Check network connectivity between services
- Ensure service names are correctly referenced in environment variables

### Database Connection Issues

**Symptoms**: Database connection errors in logs

**Solutions**:
- Verify Supabase credentials
- Check network connectivity to Supabase
- Ensure IP allowlisting is configured if necessary

### API Rate Limiting

**Symptoms**: Scraper service reports rate limit errors

**Solutions**:
- Reduce scraping frequency
- Distribute scraping jobs over time
- Consider multiple API keys for rotation 