# 🔒 Database SSL Configuration

## Overview

This document explains the database SSL configuration for the Bruno API to ensure secure database connections.

## 🚨 Security Issue Fixed

**Previous Issue**: Database connections were hardcoded to use `sslmode=disable`, meaning all database traffic was transmitted in plain text.

**Solution**: Implemented configurable SSL mode with environment-based defaults that enforce SSL in production environments.

## 🔧 Configuration

### Environment Variable

- **`DATABASE_SSL_MODE`**: Controls the SSL mode for database connections

### Valid SSL Modes

| Mode | Description | Use Case |
|------|-------------|----------|
| `disable` | No SSL encryption | Development/testing only |
| `require` | SSL required | Production (minimum security) |
| `verify-ca` | SSL with CA verification | Production (recommended) |
| `verify-full` | SSL with full certificate verification | Production (most secure) |

### Environment-Based Defaults

The application automatically sets appropriate SSL modes based on the environment:

- **Development**: `disable` (for local testing)
- **Staging**: `require` (for testing with SSL)
- **Production**: `require` (enforced SSL)

## 📝 Environment Configuration

### Development (.env)
```bash
# 🔒 Database SSL Configuration
DATABASE_SSL_MODE=disable
```

### Production (.env.production)
```bash
# 🔒 Database SSL Configuration (PRODUCTION)
# 🚨 CRITICAL: Production must use SSL for database connections
DATABASE_SSL_MODE=require
```

## 🔍 Validation

The application validates SSL mode values and will fail to start if an invalid mode is specified.

## 📊 Logging

SSL configuration is logged at startup (without sensitive data):
```
Database connection configuration: {
  "host": "localhost",
  "port": "5432", 
  "database": "bruno_site",
  "user": "postgres",
  "ssl_mode": "require",
  "env": "production"
}
```

## 🧪 Testing

For test environments, you can control SSL mode via:
```bash
TEST_DATABASE_SSL_MODE=disable
```

## 🛡️ Security Recommendations

1. **Production**: Always use `require` or higher
2. **Staging**: Use `require` to test SSL configuration
3. **Development**: `disable` is acceptable for local testing
4. **Never**: Use `disable` in production environments

## 🔄 Migration

If you're upgrading from the previous version:

1. Set `DATABASE_SSL_MODE=require` in production
2. Ensure your database supports SSL
3. Test the connection before deploying
4. Monitor logs for SSL-related issues

## 📚 Additional Resources

- [PostgreSQL SSL Documentation](https://www.postgresql.org/docs/current/ssl-tcp.html)
- [Go PostgreSQL Driver SSL Options](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters)
