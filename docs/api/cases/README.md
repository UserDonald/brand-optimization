# Integration Use Cases

This directory contains implementation guides for developers who want to integrate with the Strategic Brand Optimization Platform using different technologies and frameworks.

## Available Integration Guides

- [Next.js Integration](./nextjs.md) - How to integrate the platform with a Next.js frontend application
- [Go Client Integration](./go-client.md) - How to create a Go client that communicates with the platform
- [Node.js Integration](./nodejs.md) - How to build a Node.js application that leverages the platform's API

## Use Case Overview

Each guide provides a complete implementation example covering:

1. **Authentication** - How to implement JWT authentication for the specific technology
2. **API Communication** - How to make GraphQL queries and mutations
3. **Error Handling** - Best practices for handling API errors
4. **Real-world Examples** - Common use cases implemented in the specific technology

## General Integration Pattern

Regardless of the technology stack you're using, integrating with the Strategic Brand Optimization Platform follows this general pattern:

1. **Authentication**
   - Implement login functionality to obtain JWT tokens
   - Store tokens securely
   - Set up automatic token refresh

2. **API Client Setup**
   - Configure a GraphQL client with proper authentication
   - Set up error handling and retries

3. **Data Fetching**
   - Query competitor data
   - Fetch metrics and analytics
   - Display insights and recommendations

4. **Data Mutations**
   - Add/update competitors
   - Schedule content posts
   - Configure notifications

## Getting Started

Choose the guide that matches your technology stack to get started. If you don't see a guide for your specific technology, the patterns and approaches shown in the existing guides can be adapted to most modern frameworks and languages. 