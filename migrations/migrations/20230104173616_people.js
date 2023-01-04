exports.up = function (knex) {
  return knex.schema.createTable("people", function (table) {
    table.increments("id").primary();
    table.string("name", 255);
    table.string("surname", 255);

    table.comment("Table from people");
    table.index(
      ["id"],
      "idx_people_1"
    );
    table.engine("InnoDB");
  });
};

exports.down = function (knex) {
  return knex.schema.dropTable("people");
};
