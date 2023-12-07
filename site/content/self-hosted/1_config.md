---
title: "Configuration for Self Hosted Setup"
date: 2023-12-07T01:10:25+04:00
weight: 10
menuTitle: "Self-Hosted config"
---

```yml
bot:
  config:
    # Self-hosted configuration
    api:
      address: "localhost:8080"  # Address for the API server to listen on

    persistence:
      type: database  # Type of persistence to use (memory or database)
      db_config:
        user: "username"  # Database username
        password: "password"  # Database password
        host: "localhost"  # Database host
        port: 5432  # Database port
        database: "bot_db"  # Database name
        no_ssl: false  # Disable SSL for database connection

    assets:
      provider: fs  # Assets provider (fs for filesystem)
      params:
        root: "/tmp/assets"  # Root directory for filesystem assets

    paymentProviders:
      - name: stripe  # Payment provider name
        token: "your_stripe_token"  # Stripe API token

  handlers:
    # Handlers configuration

  api:
    # API handlers configuration
```

This example demonstrates the self-hosted configuration, including the API server settings, persistence type,
database configuration, assets provider, and payment providers.

## API Configuration (api)

```yml
api:
  address: "localhost:8080"
```

 * `address`: Specifies the address the API server should listen on. In this example,
 the server will listen on localhost at port 8080.

## Persistence Configuration (persistence)

```yml
persistence:
  type: database
  db_config:
    user: "username"
    password: "password"
    host: "localhost"
    port: 5432
    database: "bot_db"
    no_ssl: false
```

 * `type`: Specifies the type of persistence to use. It can be either memory or database.
 * `db_config`: Configuration specific to database persistence.
   * `user`: Database username for connection.
   * `password`: Database password for connection.
   * `host`: Database host address.
   * `port`: Database port.
   * `database`: Database name.
   * `no_ssl`: A boolean indicating whether to disable SSL for the database connection.

## Assets Configuration (assets)

```yml
assets:
  provider: fs
  params:
    root: "/tmp/assets"
```

 * `provider`: Specifies the assets provider. In this case, it's set to fs, indicating the filesystem provider.
 * `params`: Additional parameters specific to the chosen provider.
   * `root`: The root directory for the filesystem assets.

## Payment Providers Configuration (paymentProviders)

```yml
paymentProviders:
  - name: stripe
    token: "your_stripe_token"
```

 * `paymentProviders`: A list of payment providers with their respective configurations.
   * `name`: The name of the payment provider.
   * `token`: The API token associated with the payment provider (e.g., Stripe). Replace "your\_stripe\_token" with the actual token.
