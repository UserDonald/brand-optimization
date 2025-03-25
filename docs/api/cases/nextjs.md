# Next.js Integration Guide

This guide demonstrates how to integrate the Strategic Brand Optimization Platform with a Next.js application. We'll cover authentication, data fetching, and common use cases.

## Table of Contents

1. [Setup](#setup)
2. [Authentication](#authentication)
3. [Data Fetching](#data-fetching)
4. [Common Use Cases](#common-use-cases)

## Setup

### Dependencies

Install the required packages:

```bash
npm install @apollo/client graphql jwt-decode
# or with yarn
yarn add @apollo/client graphql jwt-decode
```

### Project Structure

We recommend the following structure for your Next.js project:

```
my-nextjs-app/
├── components/           # React components
├── lib/                  # Utility functions
│   ├── apollo-client.js  # Apollo client configuration
│   └── auth.js           # Authentication utilities
└── pages/                # Next.js pages
```

## Authentication

### Apollo Client Setup with Authentication

```javascript
// lib/apollo-client.js
import { ApolloClient, InMemoryCache, HttpLink, ApolloLink } from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { useMemo } from 'react';

let apolloClient;

function createApolloClient() {
  // Auth link adds the token to requests
  const authLink = new ApolloLink((operation, forward) => {
    const token = typeof window !== 'undefined' ? localStorage.getItem('accessToken') : '';
    
    operation.setContext({
      headers: {
        authorization: token ? `Bearer ${token}` : ''
      }
    });
    
    return forward(operation);
  });

  // Error handling link for token refresh
  const errorLink = onError(({ graphQLErrors, operation, forward }) => {
    if (graphQLErrors) {
      for (const err of graphQLErrors) {
        if (err.extensions?.code === 'AUTHENTICATION_ERROR') {
          return new Promise(resolve => {
            refreshToken()
              .then(newToken => {
                // Retry with new token
                const oldHeaders = operation.getContext().headers;
                operation.setContext({
                  headers: {
                    ...oldHeaders,
                    authorization: `Bearer ${newToken}`
                  }
                });
                resolve(forward(operation));
              })
              .catch(() => {
                // Redirect to login on failure
                if (typeof window !== 'undefined') {
                  window.location.href = '/login';
                }
                resolve();
              });
          });
        }
      }
    }
  });

  // HTTP link to the API
  const httpLink = new HttpLink({
    uri: 'http://localhost:8080/query',
  });

  return new ApolloClient({
    link: ApolloLink.from([authLink, errorLink, httpLink]),
    cache: new InMemoryCache()
  });
}

// Token refresh function
async function refreshToken() {
  const refreshToken = localStorage.getItem('refreshToken');
  
  const response = await fetch('http://localhost:8080/query', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      query: `
        mutation RefreshToken($refreshToken: String!) {
          refreshToken(refreshToken: $refreshToken) {
            accessToken
            refreshToken
          }
        }
      `,
      variables: { refreshToken }
    })
  });
  
  const data = await response.json();
  
  if (data.errors) {
    throw new Error('Token refresh failed');
  }
  
  const { accessToken, refreshToken: newRefreshToken } = data.data.refreshToken;
  localStorage.setItem('accessToken', accessToken);
  localStorage.setItem('refreshToken', newRefreshToken);
  
  return accessToken;
}

// Initialize Apollo Client with SSR support
export function initializeApollo(initialState = null) {
  const _apolloClient = apolloClient ?? createApolloClient();

  // If your page has Next.js data fetching methods that use Apollo Client,
  // the initial state gets hydrated here
  if (initialState) {
    _apolloClient.cache.restore(initialState);
  }

  // For SSR, always create a new Apollo Client
  if (typeof window === 'undefined') return _apolloClient;

  // Create the Apollo Client once in the client
  if (!apolloClient) apolloClient = _apolloClient;
  return _apolloClient;
}

export function useApollo(initialState) {
  const store = useMemo(() => initializeApollo(initialState), [initialState]);
  return store;
}
```

### Login Component

```jsx
// pages/login.js
import { useState } from 'react';
import { useMutation, gql } from '@apollo/client';
import { useRouter } from 'next/router';

const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      accessToken
      refreshToken
      user {
        id
        email
      }
    }
  }
`;

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [login, { loading, error }] = useMutation(LOGIN_MUTATION);
  const router = useRouter();

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    try {
      const { data } = await login({
        variables: { email, password }
      });
      
      localStorage.setItem('accessToken', data.login.accessToken);
      localStorage.setItem('refreshToken', data.login.refreshToken);
      
      router.push('/dashboard');
    } catch (err) {
      console.error('Login failed', err);
    }
  };

  return (
    <div>
      <h1>Login</h1>
      {error && <p>Error: {error.message}</p>}
      
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input 
            type="email" 
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        
        <div>
          <label>Password:</label>
          <input 
            type="password" 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        
        <button type="submit" disabled={loading}>
          {loading ? 'Logging in...' : 'Login'}
        </button>
      </form>
    </div>
  );
}
```

### Wrapping the App with Apollo Provider

```jsx
// pages/_app.js
import { ApolloProvider } from '@apollo/client';
import { useApollo } from '../lib/apollo-client';

function MyApp({ Component, pageProps }) {
  const apolloClient = useApollo(pageProps.initialApolloState);
  
  return (
    <ApolloProvider client={apolloClient}>
      <Component {...pageProps} />
    </ApolloProvider>
  );
}

export default MyApp;
```

## Data Fetching

### Client-Side Data Fetching

```jsx
// components/CompetitorList.js
import { useQuery, gql } from '@apollo/client';

const GET_COMPETITORS = gql`
  query GetCompetitors {
    competitors {
      id
      name
      platform
    }
  }
`;

export default function CompetitorList() {
  const { loading, error, data } = useQuery(GET_COMPETITORS);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div>
      <h2>Competitors</h2>
      <ul>
        {data.competitors.map(competitor => (
          <li key={competitor.id}>
            {competitor.name} ({competitor.platform})
          </li>
        ))}
      </ul>
    </div>
  );
}
```

### Server-Side Rendering with Data

```jsx
// pages/dashboard.js
import { gql } from '@apollo/client';
import { initializeApollo } from '../lib/apollo-client';
import CompetitorList from '../components/CompetitorList';

const GET_COMPETITORS = gql`
  query GetCompetitors {
    competitors {
      id
      name
      platform
    }
  }
`;

export default function Dashboard() {
  return (
    <div>
      <h1>Dashboard</h1>
      <CompetitorList />
    </div>
  );
}

export async function getServerSideProps(context) {
  // Get client headers to forward auth token if present
  const apolloClient = initializeApollo();
  
  try {
    // Prefetch data
    await apolloClient.query({
      query: GET_COMPETITORS
    });
    
    return {
      props: {
        initialApolloState: apolloClient.cache.extract()
      }
    };
  } catch (error) {
    // If auth error, we'll let client-side handle the redirect
    return {
      props: {
        initialApolloState: apolloClient.cache.extract()
      }
    };
  }
}
```

## Common Use Cases

### Adding a Competitor

```jsx
// components/AddCompetitor.js
import { useState } from 'react';
import { useMutation, gql } from '@apollo/client';

const ADD_COMPETITOR = gql`
  mutation AddCompetitor($input: AddCompetitorInput!) {
    addCompetitor(input: $input) {
      id
      name
      platform
    }
  }
`;

export default function AddCompetitor() {
  const [name, setName] = useState('');
  const [platform, setPlatform] = useState('');
  const [addCompetitor, { loading }] = useMutation(ADD_COMPETITOR, {
    refetchQueries: ['GetCompetitors'],
    onError: (error) => {
      console.error('Failed to add competitor:', error);
    }
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    await addCompetitor({
      variables: {
        input: { name, platform }
      }
    });
    
    // Reset form
    setName('');
    setPlatform('');
  };

  return (
    <form onSubmit={handleSubmit}>
      <h3>Add Competitor</h3>
      <div>
        <label>Name:</label>
        <input 
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>
      <div>
        <label>Platform:</label>
        <select 
          value={platform}
          onChange={(e) => setPlatform(e.target.value)}
          required
        >
          <option value="">Select...</option>
          <option value="instagram">Instagram</option>
          <option value="twitter">Twitter</option>
          <option value="facebook">Facebook</option>
          <option value="linkedin">LinkedIn</option>
        </select>
      </div>
      <button type="submit" disabled={loading}>
        {loading ? 'Adding...' : 'Add Competitor'}
      </button>
    </form>
  );
}
```

### Displaying Comparison Metrics

```jsx
// components/ComparisonMetrics.js
import { useQuery, gql } from '@apollo/client';

const GET_COMPARISON = gql`
  query GetComparison($competitorId: ID!, $dateRange: DateRangeInput!) {
    compareMetrics(competitorID: $competitorId, dateRange: $dateRange) {
      competitor {
        aggregates {
          totalLikes
          totalShares
          totalComments
          avgEngagementRate
        }
      }
      personal {
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
        engagementRateRatio
      }
    }
  }
`;

export default function ComparisonMetrics({ competitorId }) {
  const { loading, error, data } = useQuery(GET_COMPARISON, {
    variables: {
      competitorId,
      dateRange: {
        startDate: "2023-01-01",
        endDate: "2023-12-31"
      }
    },
    skip: !competitorId
  });

  if (!competitorId) return <p>Select a competitor to compare</p>;
  if (loading) return <p>Loading comparison...</p>;
  if (error) return <p>Error: {error.message}</p>;

  const { competitor, personal, ratios } = data.compareMetrics;

  return (
    <div>
      <h3>Comparison Metrics</h3>
      <div className="metrics-grid">
        <div>
          <h4>Your Metrics</h4>
          <p>Likes: {personal.aggregates.totalLikes}</p>
          <p>Shares: {personal.aggregates.totalShares}</p>
          <p>Comments: {personal.aggregates.totalComments}</p>
          <p>Engagement: {personal.aggregates.avgEngagementRate.toFixed(2)}%</p>
        </div>
        
        <div>
          <h4>Competitor Metrics</h4>
          <p>Likes: {competitor.aggregates.totalLikes}</p>
          <p>Shares: {competitor.aggregates.totalShares}</p>
          <p>Comments: {competitor.aggregates.totalComments}</p>
          <p>Engagement: {competitor.aggregates.avgEngagementRate.toFixed(2)}%</p>
        </div>
        
        <div>
          <h4>Comparison Ratios</h4>
          <p>Likes: {(ratios.likesRatio * 100).toFixed(0)}%</p>
          <p>Shares: {(ratios.sharesRatio * 100).toFixed(0)}%</p>
          <p>Engagement: {(ratios.engagementRateRatio * 100).toFixed(0)}%</p>
        </div>
      </div>
    </div>
  );
}
```

For more examples and best practices, refer to the complete [Next.js example repository](https://github.com/example/strategic-brand-nextjs-example). 