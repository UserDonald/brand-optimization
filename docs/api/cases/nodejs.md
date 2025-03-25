# Node.js Integration Guide

This guide demonstrates how to integrate the Strategic Brand Optimization Platform with a Node.js application. We'll cover authentication, data fetching, and common use cases.

## Table of Contents

1. [Setup](#setup)
2. [Authentication](#authentication)
3. [GraphQL Queries](#graphql-queries)
4. [Error Handling](#error-handling)
5. [Common Use Cases](#common-use-cases)

## Setup

### Dependencies

Install the required packages:

```bash
npm install graphql @apollo/client cross-fetch jsonwebtoken
# or with yarn
yarn add graphql @apollo/client cross-fetch jsonwebtoken
```

### Project Structure

We recommend organizing your Node.js integration with the following structure:

```
nodejs-client/
├── src/
│   ├── auth/          # Authentication utilities
│   ├── api/           # API client and GraphQL operations
│   ├── config/        # Configuration settings
│   └── index.js       # Main entry point
├── package.json
└── README.md
```

## Authentication

### Setting Up Authentication

```javascript
// src/auth/index.js
const fetch = require('cross-fetch');

class AuthClient {
  constructor(apiUrl) {
    this.apiUrl = apiUrl;
    this.accessToken = null;
    this.refreshToken = null;
    this.expiresAt = null;
  }

  async login(email, password) {
    const query = `
      mutation Login($email: String!, $password: String!) {
        login(email: $email, password: $password) {
          accessToken
          refreshToken
          tokenType
          expiresIn
          user {
            id
            email
            firstName
            lastName
            tenantId
            role
          }
        }
      }
    `;

    try {
      const response = await fetch(this.apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          query,
          variables: {
            email,
            password
          }
        })
      });

      const result = await response.json();
      
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      const { accessToken, refreshToken, expiresIn, user } = result.data.login;
      
      this.accessToken = accessToken;
      this.refreshToken = refreshToken;
      this.expiresAt = Date.now() + expiresIn * 1000;
      
      return {
        user,
        accessToken,
        refreshToken
      };
    } catch (error) {
      throw new Error(`Authentication failed: ${error.message}`);
    }
  }

  async refreshAccessToken() {
    const query = `
      mutation RefreshToken($refreshToken: String!) {
        refreshToken(refreshToken: $refreshToken) {
          accessToken
          refreshToken
          tokenType
          expiresIn
        }
      }
    `;

    try {
      const response = await fetch(this.apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          query,
          variables: {
            refreshToken: this.refreshToken
          }
        })
      });

      const result = await response.json();
      
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      const { accessToken, refreshToken, expiresIn } = result.data.refreshToken;
      
      this.accessToken = accessToken;
      this.refreshToken = refreshToken;
      this.expiresAt = Date.now() + expiresIn * 1000;
      
      return {
        accessToken,
        refreshToken
      };
    } catch (error) {
      throw new Error(`Token refresh failed: ${error.message}`);
    }
  }

  async getAccessToken() {
    // If token is still valid, return it
    if (this.accessToken && Date.now() < this.expiresAt) {
      return this.accessToken;
    }
    
    // Otherwise, refresh the token
    if (this.refreshToken) {
      const { accessToken } = await this.refreshAccessToken();
      return accessToken;
    }
    
    throw new Error('Not authenticated. Please log in first.');
  }

  async logout() {
    const query = `
      mutation {
        logout
      }
    `;

    try {
      const token = await this.getAccessToken();
      
      const response = await fetch(this.apiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ query })
      });

      const result = await response.json();
      
      // Clear tokens regardless of response
      this.accessToken = null;
      this.refreshToken = null;
      this.expiresAt = null;
      
      if (result.errors) {
        throw new Error(result.errors[0].message);
      }

      return result.data.logout;
    } catch (error) {
      throw new Error(`Logout failed: ${error.message}`);
    }
  }
}

module.exports = AuthClient;
```

### Using the Auth Client

```javascript
// src/index.js
const AuthClient = require('./auth');
const { API_URL } = require('./config');

async function main() {
  try {
    // Initialize authentication client
    const authClient = new AuthClient(API_URL);
    
    // Login
    const { user } = await authClient.login('user@example.com', 'password123');
    console.log(`Logged in as: ${user.firstName} ${user.lastName} (Tenant: ${user.tenantId})`);
    
    // Use authClient in your API operations
    // ...
    
    // Logout when done
    await authClient.logout();
    console.log('Logged out successfully');
  } catch (error) {
    console.error('Error:', error.message);
  }
}

main();
```

## GraphQL Queries

### Setting Up Apollo Client

```javascript
// src/api/client.js
const { ApolloClient, InMemoryCache, HttpLink } = require('@apollo/client');
const fetch = require('cross-fetch');

function createApiClient(authClient) {
  // Create HTTP link with authentication
  const httpLink = new HttpLink({
    uri: authClient.apiUrl,
    fetch: async (uri, options) => {
      try {
        const token = await authClient.getAccessToken();
        options.headers.authorization = `Bearer ${token}`;
      } catch (error) {
        console.error('Failed to get access token:', error.message);
      }
      return fetch(uri, options);
    }
  });

  // Create Apollo Client
  return new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache(),
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'network-only',
        errorPolicy: 'all',
      },
      query: {
        fetchPolicy: 'network-only',
        errorPolicy: 'all',
      },
      mutate: {
        errorPolicy: 'all',
      },
    },
  });
}

module.exports = { createApiClient };
```

### Making GraphQL Queries

```javascript
// src/api/operations.js
const { gql } = require('@apollo/client');

// Queries
const GET_COMPETITORS = gql`
  query GetCompetitors {
    competitors {
      id
      name
      platform
    }
  }
`;

const GET_COMPETITOR_METRICS = gql`
  query GetCompetitorMetrics($competitorId: ID!, $dateRange: DateRangeInput!) {
    competitorMetrics(competitorID: $competitorId, dateRange: $dateRange) {
      competitorID
      postID
      likes
      shares
      comments
      engagementRate
      postedAt
    }
  }
`;

const COMPARE_METRICS = gql`
  query CompareMetrics($competitorId: ID!, $dateRange: DateRangeInput!) {
    compareMetrics(competitorID: $competitorId, dateRange: $dateRange) {
      competitor {
        metrics {
          likes
          shares
          comments
          engagementRate
        }
        aggregates {
          totalLikes
          totalShares
          totalComments
          avgEngagementRate
        }
      }
      personal {
        metrics {
          likes
          shares
          comments
          engagementRate
        }
        aggregates {
          totalLikes
          totalShares
          totalComments
          avgEngagementRate
        }
      }
      ratios {
        likesRatio
        sharesRatio
        commentsRatio
        engagementRateRatio
      }
    }
  }
`;

// Mutations
const ADD_COMPETITOR = gql`
  mutation AddCompetitor($input: AddCompetitorInput!) {
    addCompetitor(input: $input) {
      id
      name
      platform
    }
  }
`;

const SCHEDULE_POST = gql`
  mutation SchedulePost($input: SchedulePostInput!) {
    schedulePost(input: $input) {
      id
      content
      scheduledTime
      platform
      format
      status
    }
  }
`;

module.exports = {
  GET_COMPETITORS,
  GET_COMPETITOR_METRICS,
  COMPARE_METRICS,
  ADD_COMPETITOR,
  SCHEDULE_POST
};
```

### Using the API Client

```javascript
// src/api/index.js
const { createApiClient } = require('./client');
const operations = require('./operations');

class ApiService {
  constructor(authClient) {
    this.client = createApiClient(authClient);
    this.operations = operations;
  }

  async getCompetitors() {
    const { data } = await this.client.query({
      query: this.operations.GET_COMPETITORS
    });
    return data.competitors;
  }

  async getCompetitorMetrics(competitorId, startDate, endDate) {
    const { data } = await this.client.query({
      query: this.operations.GET_COMPETITOR_METRICS,
      variables: {
        competitorId,
        dateRange: {
          startDate,
          endDate
        }
      }
    });
    return data.competitorMetrics;
  }

  async compareMetrics(competitorId, startDate, endDate) {
    const { data } = await this.client.query({
      query: this.operations.COMPARE_METRICS,
      variables: {
        competitorId,
        dateRange: {
          startDate,
          endDate
        }
      }
    });
    return data.compareMetrics;
  }

  async addCompetitor(tenantID, name, platform) {
    const { data } = await this.client.mutate({
      mutation: this.operations.ADD_COMPETITOR,
      variables: {
        input: {
          tenantID,
          name,
          platform
        }
      }
    });
    return data.addCompetitor;
  }

  async schedulePost(content, scheduledTime, platform, format) {
    const { data } = await this.client.mutate({
      mutation: this.operations.SCHEDULE_POST,
      variables: {
        input: {
          content,
          scheduledTime,
          platform,
          format
        }
      }
    });
    return data.schedulePost;
  }
}

module.exports = ApiService;
```

## Error Handling

### Custom Error Handling Utilities

```javascript
// src/utils/error-handling.js
class ApiError extends Error {
  constructor(message, code, details = null) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
  }
}

function handleGraphQLErrors(errors) {
  if (!errors || errors.length === 0) {
    return null;
  }
  
  const firstError = errors[0];
  let code = 'UNKNOWN_ERROR';
  
  if (firstError.extensions) {
    code = firstError.extensions.code || code;
  }
  
  return new ApiError(
    firstError.message,
    code,
    errors
  );
}

function wrapApiCall(apiPromise) {
  return apiPromise.catch(error => {
    if (error.networkError) {
      throw new ApiError(
        'Network error occurred',
        'NETWORK_ERROR',
        error.networkError
      );
    }
    
    if (error.graphQLErrors) {
      throw handleGraphQLErrors(error.graphQLErrors);
    }
    
    throw error;
  });
}

module.exports = {
  ApiError,
  handleGraphQLErrors,
  wrapApiCall
};
```

### Using Error Handling

```javascript
// src/index.js
const AuthClient = require('./auth');
const ApiService = require('./api');
const { wrapApiCall } = require('./utils/error-handling');
const { API_URL } = require('./config');

async function main() {
  try {
    // Initialize authentication client
    const authClient = new AuthClient(API_URL);
    
    // Login
    const { user } = await authClient.login('user@example.com', 'password123');
    console.log(`Logged in as: ${user.firstName} ${user.lastName} (Tenant: ${user.tenantId})`);
    
    // Initialize API service
    const api = new ApiService(authClient);
    
    try {
      // Get competitors with error handling
      const competitors = await wrapApiCall(api.getCompetitors());
      console.log('Competitors:', competitors);
      
      // Add a competitor with error handling
      const newCompetitor = await wrapApiCall(
        api.addCompetitor(user.tenantId, 'New Competitor', 'instagram')
      );
      console.log('Added competitor:', newCompetitor);
    } catch (apiError) {
      console.error(`API Error (${apiError.code}):`, apiError.message);
      // Handle specific error codes
      if (apiError.code === 'AUTHENTICATION_ERROR') {
        console.log('Please login again');
      }
    }
    
    // Logout when done
    await authClient.logout();
    console.log('Logged out successfully');
  } catch (error) {
    console.error('Error:', error.message);
  }
}

main();
```

## Common Use Cases

### 1. Comparing Metrics with a Competitor

```javascript
// examples/compare-metrics.js
const AuthClient = require('../src/auth');
const ApiService = require('../src/api');
const { API_URL } = require('../src/config');

async function compareWithCompetitor(competitorId) {
  // Initialize clients
  const authClient = new AuthClient(API_URL);
  const api = new ApiService(authClient);
  
  try {
    // Login
    await authClient.login('user@example.com', 'password123');
    
    // Set date range for the last 30 days
    const endDate = new Date().toISOString().split('T')[0];
    const startDate = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)
      .toISOString().split('T')[0];
    
    // Get comparison metrics
    const comparison = await api.compareMetrics(competitorId, startDate, endDate);
    
    // Print results
    console.log('\n--- Comparison Results ---');
    console.log('Your metrics:');
    console.log(`- Total Likes: ${comparison.personal.aggregates.totalLikes}`);
    console.log(`- Total Shares: ${comparison.personal.aggregates.totalShares}`);
    console.log(`- Total Comments: ${comparison.personal.aggregates.totalComments}`);
    console.log(`- Avg Engagement Rate: ${comparison.personal.aggregates.avgEngagementRate.toFixed(2)}%`);
    
    console.log('\nCompetitor metrics:');
    console.log(`- Total Likes: ${comparison.competitor.aggregates.totalLikes}`);
    console.log(`- Total Shares: ${comparison.competitor.aggregates.totalShares}`);
    console.log(`- Total Comments: ${comparison.competitor.aggregates.totalComments}`);
    console.log(`- Avg Engagement Rate: ${comparison.competitor.aggregates.avgEngagementRate.toFixed(2)}%`);
    
    console.log('\nPerformance ratios (you vs competitor):');
    console.log(`- Likes: ${(comparison.ratios.likesRatio * 100).toFixed(0)}%`);
    console.log(`- Shares: ${(comparison.ratios.sharesRatio * 100).toFixed(0)}%`);
    console.log(`- Comments: ${(comparison.ratios.commentsRatio * 100).toFixed(0)}%`);
    console.log(`- Engagement Rate: ${(comparison.ratios.engagementRateRatio * 100).toFixed(0)}%`);
    
  } catch (error) {
    console.error('Error comparing metrics:', error.message);
  } finally {
    // Logout
    try {
      await authClient.logout();
    } catch (error) {
      console.warn('Logout failed:', error.message);
    }
  }
}

// Run the comparison
compareWithCompetitor('competitor-123');
```

### 2. Scheduling Multiple Social Media Posts

```javascript
// examples/schedule-posts.js
const AuthClient = require('../src/auth');
const ApiService = require('../src/api');
const { API_URL } = require('../src/config');

async function scheduleContentPosts(posts) {
  // Initialize clients
  const authClient = new AuthClient(API_URL);
  const api = new ApiService(authClient);
  
  try {
    // Login
    await authClient.login('user@example.com', 'password123');
    
    console.log('Scheduling posts...');
    
    // Schedule each post
    const scheduledPosts = [];
    for (const post of posts) {
      try {
        const result = await api.schedulePost(
          post.content,
          post.scheduledTime,
          post.platform,
          post.format
        );
        
        scheduledPosts.push(result);
        console.log(`Scheduled post for ${post.platform} at ${post.scheduledTime}`);
      } catch (error) {
        console.error(`Failed to schedule post: ${error.message}`);
      }
    }
    
    console.log('\nSuccessfully scheduled posts:');
    scheduledPosts.forEach((post, index) => {
      console.log(`${index + 1}. ${post.platform} - ${post.scheduledTime}`);
      console.log(`   Content: ${post.content.substring(0, 50)}${post.content.length > 50 ? '...' : ''}`);
      console.log(`   Status: ${post.status}`);
      console.log(`   ID: ${post.id}`);
      console.log('');
    });
    
  } catch (error) {
    console.error('Error scheduling posts:', error.message);
  } finally {
    // Logout
    try {
      await authClient.logout();
    } catch (error) {
      console.warn('Logout failed:', error.message);
    }
  }
}

// Example posts to schedule
const postsToSchedule = [
  {
    content: 'Excited to announce our new feature launch! #innovation',
    scheduledTime: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(), // tomorrow
    platform: 'twitter',
    format: 'text'
  },
  {
    content: 'Check out our latest product demo video',
    scheduledTime: new Date(Date.now() + 2 * 24 * 60 * 60 * 1000).toISOString(), // day after tomorrow
    platform: 'linkedin',
    format: 'video'
  },
  {
    content: 'Behind the scenes look at our team building our latest product',
    scheduledTime: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000).toISOString(), // in 3 days
    platform: 'instagram',
    format: 'image'
  }
];

// Run the function
scheduleContentPosts(postsToSchedule);
```

### 3. Setting Up a Monitoring Dashboard

```javascript
// examples/monitoring-dashboard.js
const express = require('express');
const AuthClient = require('../src/auth');
const ApiService = require('../src/api');
const { API_URL } = require('../src/config');

// Create Express app
const app = express();
const port = 3000;

// Initialize clients
const authClient = new AuthClient(API_URL);
let api = null;

// Middleware to check authentication
async function ensureAuthenticated(req, res, next) {
  try {
    // Try to get token, which will refresh if needed
    await authClient.getAccessToken();
    next();
  } catch (error) {
    res.redirect('/login');
  }
}

// Setup routes
app.use(express.urlencoded({ extended: true }));
app.use(express.json());

// Login page
app.get('/login', (req, res) => {
  res.send(`
    <h1>Login</h1>
    <form method="post" action="/login">
      <div>
        <label>Email:</label>
        <input type="email" name="email" required />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" name="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
  `);
});

// Login handler
app.post('/login', async (req, res) => {
  try {
    const { email, password } = req.body;
    const { user } = await authClient.login(email, password);
    
    // Initialize API after login
    api = new ApiService(authClient);
    
    console.log(`Logged in as: ${user.firstName} ${user.lastName}`);
    res.redirect('/dashboard');
  } catch (error) {
    res.send(`Login failed: ${error.message}`);
  }
});

// Dashboard page
app.get('/dashboard', ensureAuthenticated, async (req, res) => {
  try {
    // Get competitors
    const competitors = await api.getCompetitors();
    
    // Build HTML
    let html = `
      <h1>Competitor Monitoring Dashboard</h1>
      <p><a href="/logout">Logout</a></p>
      
      <h2>Competitors (${competitors.length})</h2>
      <ul>
    `;
    
    // Add competitors
    for (const competitor of competitors) {
      html += `
        <li>
          <strong>${competitor.name}</strong> (${competitor.platform})
          <a href="/competitor/${competitor.id}">View Details</a>
        </li>
      `;
    }
    
    html += `
      </ul>
      
      <h2>Add Competitor</h2>
      <form method="post" action="/competitor/add">
        <div>
          <label>Name:</label>
          <input type="text" name="name" required />
        </div>
        <div>
          <label>Platform:</label>
          <select name="platform" required>
            <option value="instagram">Instagram</option>
            <option value="twitter">Twitter</option>
            <option value="facebook">Facebook</option>
            <option value="linkedin">LinkedIn</option>
          </select>
        </div>
        <button type="submit">Add Competitor</button>
      </form>
    `;
    
    res.send(html);
  } catch (error) {
    res.send(`Error: ${error.message}`);
  }
});

// Competitor details page
app.get('/competitor/:id', ensureAuthenticated, async (req, res) => {
  try {
    const competitorId = req.params.id;
    
    // Set date range for the last 30 days
    const endDate = new Date().toISOString().split('T')[0];
    const startDate = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000)
      .toISOString().split('T')[0];
    
    // Get comparison metrics
    const comparison = await api.compareMetrics(competitorId, startDate, endDate);
    
    // Build HTML
    let html = `
      <h1>Competitor Comparison</h1>
      <p><a href="/dashboard">Back to Dashboard</a></p>
      
      <h2>Performance (Last 30 Days)</h2>
      
      <h3>Your Metrics</h3>
      <ul>
        <li>Total Likes: ${comparison.personal.aggregates.totalLikes}</li>
        <li>Total Shares: ${comparison.personal.aggregates.totalShares}</li>
        <li>Total Comments: ${comparison.personal.aggregates.totalComments}</li>
        <li>Avg Engagement Rate: ${comparison.personal.aggregates.avgEngagementRate.toFixed(2)}%</li>
      </ul>
      
      <h3>Competitor Metrics</h3>
      <ul>
        <li>Total Likes: ${comparison.competitor.aggregates.totalLikes}</li>
        <li>Total Shares: ${comparison.competitor.aggregates.totalShares}</li>
        <li>Total Comments: ${comparison.competitor.aggregates.totalComments}</li>
        <li>Avg Engagement Rate: ${comparison.competitor.aggregates.avgEngagementRate.toFixed(2)}%</li>
      </ul>
      
      <h3>Performance Ratios (You vs Competitor)</h3>
      <ul>
        <li>Likes: ${(comparison.ratios.likesRatio * 100).toFixed(0)}%</li>
        <li>Shares: ${(comparison.ratios.sharesRatio * 100).toFixed(0)}%</li>
        <li>Comments: ${(comparison.ratios.commentsRatio * 100).toFixed(0)}%</li>
        <li>Engagement Rate: ${(comparison.ratios.engagementRateRatio * 100).toFixed(0)}%</li>
      </ul>
    `;
    
    res.send(html);
  } catch (error) {
    res.send(`Error: ${error.message}`);
  }
});

// Add competitor handler
app.post('/competitor/add', ensureAuthenticated, async (req, res) => {
  try {
    const { name, platform } = req.body;
    
    // Use a dummy tenant ID - in a real app, you'd get this from the session
    const tenantId = 'tenant-123';
    
    const competitor = await api.addCompetitor(tenantId, name, platform);
    res.redirect('/dashboard');
  } catch (error) {
    res.send(`Error adding competitor: ${error.message}`);
  }
});

// Logout handler
app.get('/logout', async (req, res) => {
  try {
    await authClient.logout();
    res.redirect('/login');
  } catch (error) {
    res.send(`Logout failed: ${error.message}`);
  }
});

// Start server
app.listen(port, () => {
  console.log(`Dashboard running at http://localhost:${port}`);
});
```

For more examples and detailed API documentation, refer to the [API Documentation](../api.md). 