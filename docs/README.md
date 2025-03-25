# Strategic Brand Optimization Platform Documentation

Welcome to the Strategic Brand Optimization Platform documentation. This platform enables organizations to track, analyze, and optimize brand performance across social media platforms compared to competitors. The documentation is organized into the following sections:

## Documentation Sections

### [Architecture](./architecture/overview.md)
Overview of the system architecture, communication patterns, and design decisions.

### [API Documentation](./api/api.md)
Comprehensive API guide with detailed instructions for using the GraphQL API.

### [Database Schema](./db/supabase_schema.sql)
Database schema used by the platform, including tables, relationships, and RLS policies.

### [Service Documentation](./services/README.md)
Detailed documentation for each microservice, including capabilities and interfaces.

### [Deployment Guide](./deployment/deployment.md)
Instructions for deploying the platform to various environments.

### [Development Guide](./development/development.md)
Guidelines for setting up a development environment and contributing to the project.

## Quick Start

To get up and running quickly, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/go-competitor.git
   cd go-competitor
   ```

2. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env file with your Supabase credentials
   ```

3. Start the services using Docker Compose:
   ```bash
   docker-compose up
   ```

4. Access the GraphQL API at http://localhost:8080/query

## System Overview

The Strategic Brand Optimization Platform is a microservice architecture built in Go, designed to provide in-depth brand performance analysis compared to competitors. The platform:

1. Collects and normalizes social media data
2. Provides side-by-side competitor analysis
3. Segments audience data for targeted insights
4. Offers AI-powered content optimization recommendations
5. Supports multi-tenant isolation through Supabase RLS

For a more comprehensive overview, see the [Architecture Overview](./architecture/overview.md). 