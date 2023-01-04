"use strict";
module.exports = {
  development: {
    client: "pg",
    connection: {
      host: process.env["DATASOURCES_PEOPLE_ADDR"] || "localhost",
      port: process.env["DATASOURCES_PEOPLE_PORT"] || 5432,
      database: process.env["DATASOURCES_PEOPLE_HOST"] || "people_db",
      user: process.env["DATASOURCES_PEOPLE_USER"] || "people",
      password:
        process.env["DATASOURCES_PEOPLE_PASSWORD"] || "localhost",
    },
    pool: {
      min: process.env["DATASOURCES_PEOPLE_OPTIONS_POOL_MAX"] || 1,
      max: process.env["DATASOURCES_PEOPLE_OPTIONS_POOL_MIN"] || 10,
    },
    migrations: {
      tableName: "knex_migrations",
    }, 
  },
  test: {
    client: "sqlite3",
    useNullAsDefault: true,
    connection: {
      filename: "./test.sqlite3",
    },
  }
};
